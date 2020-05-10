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
	CandidateID int64  // Issue Index (or internal candidate index, later on)
	Position    uint64 // Two Candidates may share the same Position (perfect equality)
	MedianGrade uint8
	Score       string
	Tally       *PollCandidateTally
	CreatedUnix timeutil.TimeStamp
}

// PollCandidateResults implements sort.Interface based on the Score field.
type PollCandidateResults []*PollCandidateResult

func (a PollCandidateResults) Len() int           { return len(a) }
func (a PollCandidateResults) Less(i, j int) bool { return a[i].Score < a[j].Score }
func (a PollCandidateResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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
	switch result.MedianGrade {
	case 0:
		return "red"
	case 1:
		return "red"
	case 2:
		return "orange"
	case 3:
		return "yellow"
	case 4:
		return "olive"
	case 5:
		return "green"
	default:
		return "green"
	}
}
