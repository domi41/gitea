// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
	//"sort"
	"testing"
	//"time"

	"github.com/stretchr/testify/assert"
)

func TestTally(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	repo := AssertExistsAndLoadBean(t, &Repository{ID: 1}).(*Repository)
	userAli := AssertExistsAndLoadBean(t, &User{ID: 1}).(*User)
	userBob := AssertExistsAndLoadBean(t, &User{ID: 2}).(*User)

	poll, err := CreatePoll(&CreatePollOptions{
		Repo:    repo,
		Author:  userAli,
		Subject: "Demand",
	})
	assert.NoError(t, err)
	assert.NotNil(t, poll)

	pnt := &PollNaiveTallier{}

	// No judgments yet

	tally, err := pnt.Tally(poll)
	assert.NoError(t, err)
	assert.NotNil(t, tally)

	assert.Len(t, tally.Candidates, 0)
	assert.Equal(t, uint64(0), tally.MaxJudgmentsAmount)

	// Add a Judgment from Ali

	judgment, err := CreateJudgment(&CreateJudgmentOptions{
		Judge:       userAli,
		Poll:        poll,
		Grade:       3,
		CandidateID: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	tally, err = pnt.Tally(poll)
	assert.NoError(t, err)
	assert.NotNil(t, tally)

	assert.Len(t, tally.Candidates, 1)
	assert.Equal(t, uint64(1), tally.MaxJudgmentsAmount)

	// Add a judgment from Bob

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userBob,
		Poll:        poll,
		Grade:       2,
		CandidateID: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	tally, err = pnt.Tally(poll)
	assert.NoError(t, err)
	assert.NotNil(t, tally)

	assert.Len(t, tally.Candidates, 1)
	assert.Equal(t, uint64(2), tally.MaxJudgmentsAmount)

	// Add another judgment from Ali, on another candidate

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userAli,
		Poll:        poll,
		Grade:       3,
		CandidateID: 2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	tally, err = pnt.Tally(poll)
	assert.NoError(t, err)
	assert.NotNil(t, tally)

	assert.Len(t, tally.Candidates, 2)
	assert.Equal(t, uint64(2), tally.MaxJudgmentsAmount)
}

func TestPollCandidateTally_GetMedian(t *testing.T) {
	tally := buildCandidateTally([]int{2, 3, 5, 7, 11, 13})
	assert.Equal(t, uint8(4), tally.GetMedian())

	tally = buildCandidateTally([]int{0, 0, 0, 0, 0, 0})
	assert.Equal(t, uint8(0), tally.GetMedian())

	tally = buildCandidateTally([]int{0, 0, 0, 1, 0, 0})
	assert.Equal(t, uint8(3), tally.GetMedian())

	tally = buildCandidateTally([]int{0, 0, 1, 0, 0, 1})
	assert.Equal(t, uint8(2), tally.GetMedian())

	tally = buildCandidateTally([]int{1, 0, 1, 0, 0, 1})
	assert.Equal(t, uint8(2), tally.GetMedian())

	tally = buildCandidateTally([]int{1, 0, 1, 0, 0, 3})
	assert.Equal(t, uint8(5), tally.GetMedian())

	tally = buildCandidateTally([]int{0, 2, 2})
	assert.Equal(t, uint8(1), tally.GetMedian())
}

func buildCandidateTally(grades []int) (_ *PollCandidateTally) {
	things := make([]*PollCandidateGradeTally, 0, len(grades))
	totalAmount := 0
	for grade, amount := range grades {
		things = append(things, &PollCandidateGradeTally{
			Grade:       uint8(grade),
			Amount:      uint64(amount),
			CreatedUnix: timeutil.TimeStampNow(),
		})
		totalAmount += amount
	}

	return &PollCandidateTally{
		Poll:            nil, // mock me
		CandidateID:     0,
		Grades:          things,
		JudgmentsAmount: uint64(totalAmount),
		CreatedUnix:     timeutil.TimeStampNow(),
	}
}
