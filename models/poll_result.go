// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	//"code.gitea.io/gitea/modules/references"
	"code.gitea.io/gitea/modules/timeutil"
)

type PollCandidateResult struct {
	Poll        *Poll
	CandidateID int64 // Issue Index (or internal candidate index, later on)
	Position    int64 // Two Candidates may share the same Position (perfect equality)
	MedianGrade uint8
	Tally       *PollCandidateTally
	CreatedUnix timeutil.TimeStamp
}

type PollResult struct {
	Poll        *Poll
	Tally       *PollTally
	Candidates  []*PollCandidateResult
	CreatedUnix timeutil.TimeStamp
}

func (result *PollResult) GetCandidate(candidateID int64) (_ *PollCandidateResult) {
	for _, candidate := range result.Candidates {
		if candidate.CandidateID == candidateID {
			return candidate
		}
	}
	return nil
}

func (result *PollCandidateResult) GetColorWord() (_ string) {
	return "red" // FIXME
}
