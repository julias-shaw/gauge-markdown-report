/*----------------------------------------------------------------
 *  Copyright (c) ThoughtWorks, Inc.
 *  Licensed under the Apache License, Version 2.0
 *  See LICENSE in the project root for license information.
 *----------------------------------------------------------------*/
package regenerate

import (
	"fmt"
	"os"

	"github.com/getgauge/gauge-proto/go/gauge_messages"
	"github.com/getgauge/html-report/env"
	"github.com/getgauge/html-report/logger"
	"github.com/getgauge/html-report/mdgen"
	"google.golang.org/protobuf/proto"
)

// loadResult reads a proto-serialized SuiteExecutionResult from inputFile
// and converts it to the internal SuiteResult shape.
//
// Pure: no filesystem writes, no logger.Fatal, no os.Exit. The caller
// decides what to do with errors. Extracted from Report so the parsing
// half can be unit-tested without spinning up the full regenerate flow.
func loadResult(inputFile, projectRoot string) (*mdgen.SuiteResult, error) {
	b, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", inputFile, err)
	}
	psr := &gauge_messages.ProtoSuiteResult{}
	if err := proto.Unmarshal(b, psr); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w", inputFile, err)
	}
	return mdgen.ToSuiteResult(projectRoot, psr), nil
}

// Report regenerates a Markdown report from a previously persisted
// last_run_result (proto-serialized SuiteResult).
func Report(inputFile, reportsDir, pRoot string) {
	res, err := loadResult(inputFile, pRoot)
	if err != nil {
		logger.Fatal(err.Error())
	}
	env.CreateDirectory(reportsDir)
	if err := mdgen.GenerateReports(res, reportsDir); err != nil {
		logger.Fatalf("Failed to regenerate report: %s", err.Error())
	}
}
