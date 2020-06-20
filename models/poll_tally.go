// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
)

//   ____               _        _____     _ _
//  / ___|_ __ __ _  __| | ___  |_   _|_ _| | |_   _
// | |  _| '__/ _` |/ _` |/ _ \   | |/ _` | | | | | |
// | |_| | | | (_| | (_| |  __/   | | (_| | | | |_| |
//  \____|_|  \__,_|\__,_|\___|   |_|\__,_|_|_|\__, |
//                                             |___/

type PollCandidateGradeTally struct {
	Grade       uint8
	Amount      uint64
	CreatedUnix timeutil.TimeStamp
}

//   ____                _ _     _       _         _____     _ _
//  / ___|__ _ _ __   __| (_) __| | __ _| |_ ___  |_   _|_ _| | |_   _
// | |   / _` | '_ \ / _` | |/ _` |/ _` | __/ _ \   | |/ _` | | | | | |
// | |__| (_| | | | | (_| | | (_| | (_| | ||  __/   | | (_| | | | |_| |
//  \____\__,_|_| |_|\__,_|_|\__,_|\__,_|\__\___|   |_|\__,_|_|_|\__, |
//                                                               |___/

type PollCandidateTally struct {
	Poll            *Poll
	CandidateID     int64                      // Issue Index (or internal candidate index, later on)
	Grades          []*PollCandidateGradeTally // Sorted by grade (0 == REJECT)
	JudgmentsAmount uint64
	CreatedUnix     timeutil.TimeStamp
}

func (pct *PollCandidateTally) Copy() (_ *PollCandidateTally) {

	grades := make([]*PollCandidateGradeTally, 0, 8)
	for _, grade := range pct.Grades {
		grades = append(grades, &PollCandidateGradeTally{
			Grade:       grade.Grade,
			Amount:      grade.Amount,
			CreatedUnix: grade.CreatedUnix,
		})
	}

	return &PollCandidateTally{
		Poll:            pct.Poll,
		CandidateID:     pct.CandidateID,
		Grades:          grades,
		JudgmentsAmount: pct.JudgmentsAmount,
		CreatedUnix:     pct.CreatedUnix,
	}
}

func (pct *PollCandidateTally) GetMedian() (_ uint8) {

	if 0 == pct.JudgmentsAmount {
		return uint8(0)
		//return 0 // to test
	}

	adjustedTotal := pct.JudgmentsAmount - 1
	//if opts.UseHighMedian {
	//	adjustedTotal := pct.JudgmentsAmount + 1
	//}
	medianIndex := adjustedTotal / 2 // Euclidean div
	cursorIndex := uint64(0)
	for _, grade := range pct.Grades {
		if 0 == grade.Amount {
			continue
		}

		startIndex := cursorIndex
		cursorIndex += grade.Amount
		endIndex := cursorIndex
		if (startIndex <= medianIndex) && (medianIndex < endIndex) {
			return grade.Grade
		}
	}
	println("warning: GetMedian defaulting to 0")
	return uint8(0)
}

func (pct *PollCandidateTally) GetBiggestGroup(aroundGrade uint8) (groupSize int, groupSign int, groupGrade uint8) {
	belowGroupSize := 0
	belowGroupSign := -1
	belowGroupGrade := uint8(0)

	aboveGroupSize := 0
	aboveGroupSign := 1
	aboveGroupGrade := uint8(0)

	for k, _ := range pct.Poll.GetGradationList() {
		grade := uint8(k)
		//for _, grade := range pct.Poll.GetGrades() {
		if grade < aroundGrade {
			belowGroupSize += int(pct.Grades[grade].Amount)
			belowGroupGrade = grade
		}
		if grade > aroundGrade {
			aboveGroupSize += int(pct.Grades[grade].Amount)
			if 0 == aboveGroupGrade {
				aboveGroupGrade = grade
			}
		}
	}

	if aboveGroupSize > belowGroupSize {
		return aboveGroupSize, aboveGroupSign, aboveGroupGrade
	}
	return belowGroupSize, belowGroupSign, belowGroupGrade
}

func (pct *PollCandidateTally) RegradeJudgments(fromGrade uint8, toGrade uint8) {
	if toGrade == fromGrade {
		return
	}
	pct.Grades[toGrade].Amount += pct.Grades[fromGrade].Amount
	pct.Grades[fromGrade].Amount = 0
}

//  _____     _ _
// |_   _|_ _| | |_   _
//   | |/ _` | | | | | |
//   | | (_| | | | |_| |
//   |_|\__,_|_|_|\__, |
//                |___/

type PollTally struct {
	Poll               *Poll
	MaxJudgmentsAmount uint64 // per candidate, will help including default grade 0=TO_REJECT
	Candidates         []*PollCandidateTally
	CreatedUnix        timeutil.TimeStamp
}

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
