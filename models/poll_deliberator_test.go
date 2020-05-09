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

func TestNaiveDeliberator(t *testing.T) {
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

	pnd := &PollNaiveDeliberator{}

	// No judgments yet

	result, err := pnd.Deliberate(poll)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Poll)
	assert.Equal(t, poll, result.Poll)
	assert.NotNil(t, result.Tally)
	assert.NotNil(t, result.Candidates)
	assert.Len(t, result.Candidates, 0)

	// Add a Judgment from Ali

	judgment, err := CreateJudgment(&CreateJudgmentOptions{
		Judge:       userAli,
		Poll:        poll,
		Grade:       3,
		CandidateID: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	result, err = pnd.Deliberate(poll)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Tally)
	assert.NotNil(t, result.Candidates)
	assert.Len(t, result.Candidates, 1)

	// Add a judgment from Bob

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userBob,
		Poll:        poll,
		Grade:       2,
		CandidateID: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	result, err = pnd.Deliberate(poll)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Tally)
	assert.NotNil(t, result.Candidates)
	assert.Len(t, result.Candidates, 1)
	assert.Equal(t, int64(0), result.Candidates[0].Position)

	// Add another judgment from Bob

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userBob,
		Poll:        poll,
		Grade:       3,
		CandidateID: 2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	result, err = pnd.Deliberate(poll)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Candidates)
	assert.Len(t, result.Candidates, 2)
	assert.Equal(t, int64(0), result.Candidates[0].Position)
	assert.Equal(t, int64(1), result.Candidates[1].Position)
	assert.Equal(t, int64(1), result.Candidates[0].CandidateID)
	assert.Equal(t, int64(2), result.Candidates[1].CandidateID)

}
