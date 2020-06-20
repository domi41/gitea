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
	userCho := AssertExistsAndLoadBean(t, &User{ID: 3}).(*User)

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
	// Candidate 1 : SOMEWHAT_GOOD
	// Candidate 2 :

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
	// Candidate 1 : PASSABLE SOMEWHAT_GOOD
	// Candidate 2 :

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
	assert.Equal(t, uint64(1), result.Candidates[0].Position)

	// Add another judgment from Bob
	// Candidate 1 : PASSABLE SOMEWHAT_GOOD
	// Candidate 2 : GOOD

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userBob,
		Poll:        poll,
		Grade:       4,
		CandidateID: 2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	result, err = pnd.Deliberate(poll)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Candidates)
	assert.Len(t, result.Candidates, 2)
	assert.Equal(t, uint64(1), result.Candidates[0].Position)
	assert.Equal(t, uint64(2), result.Candidates[1].Position)
	assert.Equal(t, int64(1), result.Candidates[0].CandidateID)
	assert.Equal(t, int64(2), result.Candidates[1].CandidateID)
	assert.Equal(t, uint8(2), result.Candidates[0].MedianGrade)
	assert.Equal(t, uint8(0), result.Candidates[1].MedianGrade)

	// Add another 2 judgments from Cho and one from Ali
	// Candidate 1 : PASSABLE SOMEWHAT_GOOD GOOD
	// Candidate 2 : SOMEWHAT_GOOD SOMEWHAT_GOOD GOOD

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userCho,
		Poll:        poll,
		Grade:       4,
		CandidateID: 1,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userCho,
		Poll:        poll,
		Grade:       3,
		CandidateID: 2,
	})
	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       userAli,
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
	assert.Equal(t, uint64(1), result.Candidates[0].Position)
	assert.Equal(t, uint64(2), result.Candidates[1].Position)
	assert.Equal(t, int64(2), result.Candidates[0].CandidateID)
	assert.Equal(t, int64(1), result.Candidates[1].CandidateID)
	//println("C1", result.Candidates[1].MedianGrade)
	assert.Equal(t, uint8(3), result.Candidates[0].MedianGrade)
	assert.Equal(t, uint8(3), result.Candidates[1].MedianGrade)

	assert.Equal(t, uint64(1), result.Candidates[0].Tally.Grades[4].Amount)
	assert.Equal(t, uint64(2), result.Candidates[0].Tally.Grades[3].Amount)
	assert.Equal(t, uint64(1), result.Candidates[1].Tally.Grades[2].Amount)
	assert.Equal(t, uint64(1), result.Candidates[1].Tally.Grades[3].Amount)
	assert.Equal(t, uint64(1), result.Candidates[1].Tally.Grades[4].Amount)
}
