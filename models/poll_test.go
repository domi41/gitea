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

func TestPoll_Create(t *testing.T) {
	assert.NoError(t, PrepareTestDatabase())

	repo := AssertExistsAndLoadBean(t, &Repository{ID: 1}).(*Repository)
	user := AssertExistsAndLoadBean(t, &User{ID: 2}).(*User)

	var opts = CreatePollOptions{
		Repo:    repo,
		Author:  user,
		Subject: "Quality",
	}
	poll, err := CreatePoll(&opts)

	assert.NoError(t, err)
	assert.NotNil(t, poll)
}
