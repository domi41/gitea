// Copyright 2014 The Gogs Authors. All rights reserved.
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
	CandidateID     int64 // Issue Index (or internal candidate index, later on)
	Grades          []*PollCandidateGradeTally
	JudgmentsAmount uint64
	CreatedUnix     timeutil.TimeStamp
}

type PollTally struct {
	Poll               *Poll
	MaxJudgmentsAmount uint64 // per candidate, will help including default grade 0=TO_REJECT
	Candidates         []*PollCandidateTally
	CreatedUnix        timeutil.TimeStamp
}

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
