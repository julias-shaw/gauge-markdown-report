/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/
package regenerate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/getgauge/gauge-proto/go/gauge_messages"
	"google.golang.org/protobuf/proto"
)

// loadResult is the unit-testable half of Report: read proto bytes, unmarshal,
// transform. It does no filesystem writes, so we can verify the read /
// unmarshal / transform steps without involving GenerateReports.

func TestLoadResult_HappyPath(t *testing.T) {
	tmp := t.TempDir()
	psr := &gauge_messages.ProtoSuiteResult{
		ProjectName: "demo",
		Environment: "ci",
		Tags:        "smoke",
	}
	b, err := proto.Marshal(psr)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	in := filepath.Join(tmp, "last_run_result")
	if err := os.WriteFile(in, b, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	got, err := loadResult(in, "/proj")
	if err != nil {
		t.Fatalf("loadResult: %v", err)
	}
	if got == nil {
		t.Fatal("loadResult returned nil suite")
	}
	if got.ProjectName != "demo" {
		t.Errorf("ProjectName = %q, want demo", got.ProjectName)
	}
	if got.Environment != "ci" {
		t.Errorf("Environment = %q, want ci", got.Environment)
	}
	if got.Tags != "smoke" {
		t.Errorf("Tags = %q, want smoke", got.Tags)
	}
}

func TestLoadResult_MissingFile(t *testing.T) {
	_, err := loadResult(filepath.Join(t.TempDir(), "does-not-exist"), "/proj")
	if err == nil {
		t.Fatal("expected an error for missing input file")
	}
	if !strings.Contains(err.Error(), "read") {
		t.Errorf("error should mention read failure, got: %v", err)
	}
}

func TestLoadResult_UnparseableProto(t *testing.T) {
	in := filepath.Join(t.TempDir(), "garbage")
	if err := os.WriteFile(in, []byte("not a valid proto"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}
	_, err := loadResult(in, "/proj")
	if err == nil {
		t.Fatal("expected an error for unparseable proto")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("error should mention unmarshal failure, got: %v", err)
	}
}
