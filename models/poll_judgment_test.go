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

func TestPollJudgment_Create(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())
	repo := AssertExistsAndLoadBean(t, &Repository{ID: 2}).(*Repository)
	user := AssertExistsAndLoadBean(t, &User{ID: 2}).(*User)
	//issue := AssertExistsAndLoadBean(t, &Issue{ID: 2}).(*Issue)
	//assert.Equal(t, repo, issue.Repo)

	poll, errp := CreatePoll(&CreatePollOptions{
		Repo:    repo,
		Author:  user,
		Subject: "Quality",
	})

	assert.NoError(t, errp)
	assert.NotNil(t, poll)

	judgment, err := CreateJudgment(&CreateJudgmentOptions{
		Judge:       user,
		Poll:        poll,
		Grade:       3,
		CandidateID: 1,
	})

	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	// Emit another Judgment, on another Candidate

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       user,
		Poll:        poll,
		Grade:       1,
		CandidateID: 2,
	})

	assert.NoError(t, err)
	assert.NotNil(t, judgment)

	// Cannot create another judgment on the same poll and candidate

	judgment, err = CreateJudgment(&CreateJudgmentOptions{
		Judge:       user,
		Poll:        poll,
		Grade:       0,
		CandidateID: 2,
	})

	assert.Error(t, err)
	assert.Nil(t, judgment)

	// … you have to update it

	judgment, err = UpdateJudgment(&UpdateJudgmentOptions{
		Judge:       user,
		Poll:        poll,
		Grade:       0,
		CandidateID: 2,
	})

	assert.NoError(t, err)
	assert.NotNil(t, judgment)
}
