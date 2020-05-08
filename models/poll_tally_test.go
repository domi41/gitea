// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
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
