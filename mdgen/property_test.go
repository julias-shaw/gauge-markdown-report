/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/

package mdgen

import (
	"bytes"
	"math/rand"
	"testing"
	"testing/quick"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// Property tests are deliberately narrow. Markdown rendering has lots of edge
// cases that don't matter (a stray paragraph break, a different bullet glyph
// — none of which a real reader would notice) and curating examples for them
// is busywork. The two properties below catch the failure modes that actually
// hurt: non-idempotent escaping (would silently double-escape on regenerate)
// and a renderer that crashes or produces output the reference parser
// rejects on shapes the curated tests don't cover.
//
// Per PLAN.md §4.6 we cap input size and pin the seed for reproducibility.

// TestEscapeMDIdempotent_Property complements the curated cases in
// format_test.go. testing/quick generates random byte strings — including
// inputs that mix already-escaped specials with unescaped ones, which is the
// regime where a naive escaper drifts.
func TestEscapeMDIdempotent_Property(t *testing.T) {
	f := func(s string) bool {
		once := escapeMD(s)
		twice := escapeMD(once)
		return once == twice
	}
	cfg := &quick.Config{
		MaxCount: 500,
		Rand:     rand.New(rand.NewSource(1)),
	}
	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}

// TestRenderIndexParsesAsGFM_Property asserts that for any randomly-shaped
// SuiteResult, RenderIndex produces a non-empty document that goldmark can
// parse and walk. goldmark itself is permissive (it never errors), so the
// real signal is that the produced AST contains at least one heading — a
// proxy for "the renderer didn't fall through to an empty / malformed
// document on this input".
func TestRenderIndexParsesAsGFM_Property(t *testing.T) {
	saved := projectRoot
	projectRoot = "/proj"
	t.Cleanup(func() { projectRoot = saved })

	gm := goldmark.New(goldmark.WithExtensions(extension.GFM))
	rng := rand.New(rand.NewSource(2026))

	for i := 0; i < 100; i++ {
		suite := genSuite(rng, 5)
		var buf bytes.Buffer
		if err := RenderIndex(&buf, suite); err != nil {
			t.Fatalf("iter %d: RenderIndex returned err: %v", i, err)
		}
		if buf.Len() == 0 {
			t.Fatalf("iter %d: RenderIndex produced no output", i)
		}
		doc := gm.Parser().Parse(text.NewReader(buf.Bytes()))
		if !hasHeading(doc) {
			t.Fatalf("iter %d: produced output has no heading\n--- output ---\n%s", i, buf.String())
		}
	}
}

func hasHeading(n ast.Node) bool {
	found := false
	_ = ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if _, ok := node.(*ast.Heading); ok {
			found = true
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	return found
}

// genSuite builds a small, deterministic SuiteResult from rng. Specs vary in
// status, name, and tag content (drawn from a charset that includes Markdown
// specials) so the renderer's escaping paths get exercised.
func genSuite(rng *rand.Rand, maxSpecs int) *SuiteResult {
	n := rng.Intn(maxSpecs) + 1
	specs := make([]*spec, n)
	for i := range specs {
		specs[i] = genSpec(rng, i)
	}
	return &SuiteResult{
		ProjectName:         randomToken(rng),
		Environment:         randomToken(rng),
		Tags:                randomToken(rng),
		Timestamp:           "Jan 2, 2026",
		ExecutionTime:       rng.Int63n(60_000),
		PassedSpecsCount:    rng.Intn(n + 1),
		FailedSpecsCount:    rng.Intn(n + 1),
		PassedScenarioCount: rng.Intn(20),
		FailedScenarioCount: rng.Intn(20),
		SpecResults:         specs,
	}
}

func genSpec(rng *rand.Rand, idx int) *spec {
	statuses := []status{pass, fail, skip}
	return &spec{
		SpecHeading:     randomToken(rng),
		FileName:        "/proj/specs/" + randomToken(rng) + ".spec",
		ExecutionTime:   rng.Int63n(10_000),
		ExecutionStatus: statuses[rng.Intn(len(statuses))],
		Tags:            []string{randomToken(rng)},
	}
}

// tokenChars deliberately includes Markdown specials so generated names
// exercise escapeMD. Limiting to a small charset keeps property failures
// debuggable when they happen.
var tokenChars = []rune("abcXY *|`<>_")

func randomToken(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	out := make([]rune, n)
	for i := range out {
		out[i] = tokenChars[rng.Intn(len(tokenChars))]
	}
	return string(out)
}
