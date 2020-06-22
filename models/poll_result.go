// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
	"strconv"
)

//   ____                _ _     _       _       ____                 _ _
//  / ___|__ _ _ __   __| (_) __| | __ _| |_ ___|  _ \ ___  ___ _   _| | |_
// | |   / _` | '_ \ / _` | |/ _` |/ _` | __/ _ \ |_) / _ \/ __| | | | | __|
// | |__| (_| | | | | (_| | | (_| | (_| | ||  __/  _ <  __/\__ \ |_| | | |_
//  \____\__,_|_| |_|\__,_|_|\__,_|\__,_|\__\___|_| \_\___||___/\__,_|_|\__|
//

type PollCandidateResult struct {
	Poll         *Poll
	CandidateID  int64  // Issue Index (or internal candidate index, later on)
	Position     uint64 // Two Candidates may share the same Position (perfect equality)
	MedianGrade  uint8
	Score        string
	Tally        *PollCandidateTally
	MeritProfile *PollCandidateMeritProfile
	CreatedUnix  timeutil.TimeStamp
}

func (result *PollCandidateResult) GetColorWord() (_ string) {
	return result.Poll.GetGradeColorWord(result.MedianGrade)
}

func (result *PollCandidateResult) GetCandidateName() (_ string) { // FIXME
	isssue, err := GetIssueByID(result.CandidateID + 1000)
	if nil != err {
		return "Candidate #" + strconv.FormatInt(result.CandidateID, 10)
	}
	return isssue.Title
}

//   ____                _ _     _       _       ____                 _ _
//  / ___|__ _ _ __   __| (_) __| | __ _| |_ ___|  _ \ ___  ___ _   _| | |_ ___
// | |   / _` | '_ \ / _` | |/ _` |/ _` | __/ _ \ |_) / _ \/ __| | | | | __/ __|
// | |__| (_| | | | | (_| | | (_| | (_| | ||  __/  _ <  __/\__ \ |_| | | |_\__ \
//  \____\__,_|_| |_|\__,_|_|\__,_|\__,_|\__\___|_| \_\___||___/\__,_|_|\__|___/
//

// PollCandidateResults implements sort.Interface based on the Score field.
type PollCandidateResults []*PollCandidateResult

func (a PollCandidateResults) Len() int           { return len(a) }
func (a PollCandidateResults) Less(i, j int) bool { return a[i].Score < a[j].Score }
func (a PollCandidateResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//  ____                 _ _
// |  _ \ ___  ___ _   _| | |_
// | |_) / _ \/ __| | | | | __|
// |  _ <  __/\__ \ |_| | | |_
// |_| \_\___||___/\__,_|_|\__|
//

type PollResult struct {
	Poll        *Poll
	Tally       *PollTally
	Candidates  PollCandidateResults
	CreatedUnix timeutil.TimeStamp
}

func (result *PollResult) GetCandidate(candidateID int64) (_ *PollCandidateResult) {
	// A `for` loop is pretty inefficient, index this somehow?
	for _, candidate := range result.Candidates {
		if candidate.CandidateID == candidateID {
			return candidate
		}
	}
	return nil
}
