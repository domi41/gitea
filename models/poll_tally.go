// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
)

type PollCandidateGradeTally struct {
	Grade       uint8
	Amount      uint64
	CreatedUnix timeutil.TimeStamp
}

type PollCandidateTally struct {
	Poll            *Poll
	CandidateID     int64                      // Issue Index (or internal candidate index, later on)
	Grades          []*PollCandidateGradeTally // sorted by grade
	JudgmentsAmount uint64
	CreatedUnix     timeutil.TimeStamp
}

func (pct *PollCandidateTally) GetMedian() (_ uint8) {
	medianIndex := pct.JudgmentsAmount / 2
	cursorIndex := uint64(0)
	for _, grade := range pct.Grades {
		if 0 < grade.Amount {
			cursorIndex += grade.Amount
			if cursorIndex >= medianIndex {
				return grade.Grade
			}
		}
	}
	println("warning: GetMedian defaulting to 0")
	return uint8(0)
}

//func (pct *PollCandidateTally) GetScore() (_ string) {
//
//}

type PollTally struct {
	Poll               *Poll
	MaxJudgmentsAmount uint64 // per candidate, will help including default grade 0=TO_REJECT
	Candidates         []*PollCandidateTally
	CreatedUnix        timeutil.TimeStamp
}

//// PollCandidateGrades implements sort.Interface based on the Score field.
//type PollCandidateGrades []*PollCandidateGradeTally
//
//func (a PollCandidateGrades) Len() int           { return len(a) }
//func (a PollCandidateGrades) Less(i, j int) bool { return a[i].Score < a[j].Score }
//func (a PollCandidateGrades) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

//  _____     _ _ _
// |_   _|_ _| | (_) ___ _ __
//   | |/ _` | | | |/ _ \ '__|
//   | | (_| | | | |  __/ |
//   |_|\__,_|_|_|_|\___|_|
//

type PollTallier interface {
	Tally(poll *Poll) (tally *PollTally, err error)
}

//  _   _       _
// | \ | | __ _(_)_   _____
// |  \| |/ _` | \ \ / / _ \
// | |\  | (_| | |\ V /  __/
// |_| \_|\__,_|_| \_/ \___|
//

type PollNaiveTallier struct{}

func (tallier *PollNaiveTallier) Tally(poll *Poll) (_ *PollTally, err error) {

	gradation := poll.GetGradationList()

	candidatesIDs, errG := poll.GetCandidatesIDs()
	if nil != errG {
		return nil, errG
	}

	candidates := make([]*PollCandidateTally, 0, 64)
	maximumAmount := uint64(0)

	for _, candidateID := range candidatesIDs {

		grades := make([]*PollCandidateGradeTally, 0, 8)
		judgmentsAmount := uint64(0)

		for gradeInt, _ := range gradation {
			grade := uint8(gradeInt)
			amount, errC := poll.CountGrades(candidateID, grade)
			if nil != errC {
				return nil, errC
			}

			judgmentsAmount += amount
			grades = append(grades, &PollCandidateGradeTally{
				Grade:       grade,
				Amount:      amount,
				CreatedUnix: timeutil.TimeStampNow(),
			})
		}

		//maximumAmount = util.Max(judgmentsAmount, maximumAmount)
		if maximumAmount < judgmentsAmount {
			maximumAmount = judgmentsAmount
		}

		candidates = append(candidates, &PollCandidateTally{
			Poll:            poll,
			CandidateID:     candidateID,
			Grades:          grades,
			JudgmentsAmount: judgmentsAmount,
			CreatedUnix:     timeutil.TimeStampNow(),
		})
	}

	tally := &PollTally{
		Poll:               poll,
		MaxJudgmentsAmount: maximumAmount,
		Candidates:         candidates,
		CreatedUnix:        timeutil.TimeStampNow(),
	}

	return tally, nil
}
