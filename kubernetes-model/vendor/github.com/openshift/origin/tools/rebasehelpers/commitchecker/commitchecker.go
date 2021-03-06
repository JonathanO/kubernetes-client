/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/openshift/origin/tools/rebasehelpers/util"
)

func main() {
	var start, end string
	flag.StringVar(&start, "start", "master", "The start of the revision range for analysis")
	flag.StringVar(&end, "end", "HEAD", "The end of the revision range for analysis")
	flag.Parse()

	commits, err := util.CommitsBetween(start, end)
	if err != nil {
		if err == util.ErrNotCommit {
			fmt.Fprintf(os.Stderr, "WARNING: one of the provided commits does not exist, not a true branch\n")
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "ERROR: couldn't find commits from %s..%s: %v\n", start, end, err)
		os.Exit(1)
	}

	// TODO: Filter out bump commits for now until we decide how to deal with
	// them correctly.
	// TODO: ...along with subtree merges.
	nonbumpCommits := []util.Commit{}
	for _, commit := range commits {
		var lastDescriptionLine string
		if descriptionLen := len(commit.Description); descriptionLen > 0 {
			lastDescriptionLine = commit.Description[descriptionLen-1]
		}
		if !strings.HasPrefix(commit.Summary, "bump(") && !strings.HasPrefix(lastDescriptionLine, "git-subtree-split:") {
			nonbumpCommits = append(nonbumpCommits, commit)
		}
	}

	errs := []string{}
	for _, validate := range AllValidators {
		if err := validate(nonbumpCommits); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(os.Getenv("RESTORE_AND_VERIFY_GODEPS")) > 0 {
		// Godeps verifies all commits, including bumps and UPSTREAM
		if err := ValidateGodeps(commits); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "%s\n", strings.Join(errs, "\n\n"))
		os.Exit(2)
	}
}
