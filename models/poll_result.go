// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	//"code.gitea.io/gitea/modules/references"
	"code.gitea.io/gitea/modules/timeutil"
)

type PollCandidateJudgmentTally struct {
	Grade       int8
	Amount      int64
	CreatedUnix timeutil.TimeStamp
}

type PollCandidateTally struct {
	Poll            *Poll
	CandidateID     int64 // Issue Index (or internal candidate index, later on)
	Judgments       []*PollCandidateJudgmentTally
	JudgmentsAmount int64
	CreatedUnix     timeutil.TimeStamp
}

type PollTally struct {
	Poll               *Poll
	MaxJudgmentsAmount int64 // per candidate, will help including default grade 0=TO_REJECT
	Candidates         []*PollCandidateTally
	CreatedUnix        timeutil.TimeStamp
}

type PollCandidateResult struct {
	Poll        *Poll
	CandidateID int64 // Issue Index (or internal candidate index, later on)
	Position    int64 // Two Candidates may share the same Position (perfect equality)
	MedianGrade int8
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
	return result.Candidates[0] // FIXME
}

func (result *PollCandidateResult) GetColorWord() (_ string) {
	return "red" // FIXME
}
