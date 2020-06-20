// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
	"math"
)

type PollCandidateMeritProfile struct {
	Poll            *Poll
	CandidateID     int64    // Issue Index (or internal candidate index, later on)
	Position        uint64   // Two Candidates may share the same Position (perfect equality)
	Grades          []uint64 // Amount of Judgments per Grade
	JudgmentsAmount uint64
	//HalfJudgmentsAmount uint64 // Hack for the template, since we can't do arithmetic
	CreatedUnix timeutil.TimeStamp
}

//  ____
// / ___|_   ____ _
// \___ \ \ / / _` |
//  ___) \ V / (_| |
// |____/ \_/ \__, |
//            |___/
//

type CartesianVector2 struct {
	X float64
	Y float64
}

type TwoCirclePoints struct {
	Start *CartesianVector2
	End   *CartesianVector2
}

func (merit *PollCandidateMeritProfile) GetGradeAngle(gradeID int, gapHalfAngle float64) (_ float64) {
	TAU := 6.283185307179586 // Sigh ; how is this not in Golang yet?
	totalAngle := TAU - gapHalfAngle*2.0
	return totalAngle * float64(merit.Grades[gradeID]) / float64(merit.JudgmentsAmount)
}

func (merit *PollCandidateMeritProfile) GetCirclePoints(gradeID int, radius float64, halfGap float64) (_ *TwoCirclePoints) {
	TAU := 6.283185307179586 // Sigh ; how is this not in Golang yet?
	gapHalfAngle := halfGap
	//gapHalfAngle = 0.39192267544687825 * 0.96  // asin(GOLDEN_RATIO-1)
	//gapHalfAngle = TAU / 4.0  // Hemicycle
	//gapHalfAngle = 0.0  // Camembert (du fromage!)
	totalAngle := TAU - gapHalfAngle*2.0
	lowerGradeJudgmentsAmount := uint64(0)
	for i := 0; i < gradeID; i++ {
		lowerGradeJudgmentsAmount += merit.Grades[i]
	}
	totalJudgments := float64(merit.JudgmentsAmount)
	startingAngle := gapHalfAngle + totalAngle*float64(lowerGradeJudgmentsAmount)/totalJudgments
	angle := totalAngle * float64(merit.Grades[gradeID]) / totalJudgments
	endingAngle := startingAngle + angle

	//println("Angles", gradeID, merit.Grades[gradeID], totalJudgments, startingAngle, endingAngle)
	return &TwoCirclePoints{
		Start: &CartesianVector2{
			X: radius * math.Cos(startingAngle),
			Y: radius * math.Sin(startingAngle),
		},
		End: &CartesianVector2{
			X: radius * math.Cos(endingAngle),
			Y: radius * math.Sin(endingAngle),
		},
	}
}

func (merit *PollCandidateMeritProfile) GetColorWord(gradeID int) (_ string) {
	return merit.Poll.GetGradeColorWord(uint8(gradeID))
}

func (merit *PollCandidateMeritProfile) GetColorCode(gradeID int) (_ string) {
	return merit.Poll.GetGradeColorCode(uint8(gradeID))
}
