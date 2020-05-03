// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	//"code.gitea.io/gitea/modules/references"
	"code.gitea.io/gitea/modules/timeutil"
	//"fmt"
	"xorm.io/xorm"
)

// A Poll on <Subject> with the issues of a repository as candidates.
type Poll struct {
	ID       int64       `xorm:"pk autoincr"`
	RepoID   int64       `xorm:"INDEX UNIQUE(repo_index)"`
	Repo     *Repository `xorm:"-"`
	Index    int64       `xorm:"UNIQUE(repo_index)"` // Index in one repository.
	AuthorID int64       `xorm:"INDEX"`
	Author   *User       `xorm:"-"`
	// When the poll is applied to all the issues, the subject should be an issue's trait.
	// eg: Quality, Importance, Urgency, Wholeness, Relevance…
	Subject string `xorm:"name"`
	// Content may be used to describe at length the constitutive details of that poll.
	// Eg: Rationale, Deliberation Consequences, Modus Operandi…
	Content         string `xorm:"TEXT"`
	RenderedContent string `xorm:"-"`
	Ref             string

	DeadlineUnix timeutil.TimeStamp `xorm:"INDEX"`
	CreatedUnix  timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix  timeutil.TimeStamp `xorm:"INDEX updated"`
	ClosedUnix   timeutil.TimeStamp `xorm:"INDEX"`

	//Judgments         []*Judgment    `xorm:"-"`
	//Judgments         JudgmentList   `xorm:"-"`
}

type CreatePollOptions struct {
	//Type  PollType  // for inline polls with their own candidates?
	Author *User
	Repo   *Repository

	Subject string
	Content string
}

func CreatePoll(opts *CreatePollOptions) (poll *Poll, err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	poll, err = createPoll(sess, opts)
	if err != nil {
		return nil, err
	}

	if err = sess.Commit(); err != nil {
		return nil, err
	}

	return poll, nil
}

func createPoll(e *xorm.Session, opts *CreatePollOptions) (_ *Poll, err error) {
	poll := &Poll{
		AuthorID: opts.Author.ID,
		Author:   opts.Author,
		Content:  opts.Content,
		RepoID:   opts.Repo.ID,
		Repo:     opts.Repo,
	}
	if _, err = e.Insert(poll); err != nil {
		return nil, err
	}

	if err = opts.Repo.getOwner(e); err != nil {
		return nil, err
	}

	//if err = updatePollInfos(e, opts, poll); err != nil {
	//	return nil, err
	//}

	return poll, nil
}
