// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	//"code.gitea.io/gitea/modules/references"
	"code.gitea.io/gitea/modules/timeutil"
	"fmt"

	//"fmt"
	"xorm.io/xorm"
)

// A Judgment on a Poll
type Judgment struct {
	ID      int64 `xorm:"pk autoincr"`
	PollID  int64 `xorm:"INDEX UNIQUE(poll_judge_candidate)"`
	Poll    *Poll `xorm:"-"`
	JudgeID int64 `xorm:"INDEX UNIQUE(poll_judge_candidate)"`
	Judge   *User `xorm:"-"`
	// Either an Issue ID or an index in the list of Candidates (for inline polls)
	CandidateID int64 `xorm:"UNIQUE(poll_judge_candidate)"`
	// There may be other graduations
	// 0 = to reject
	// 1 = poor
	// 2 = passable
	// 3 = good
	// 4 = very good
	// 5 = excellent
	// Make sure 0 always means *something* in your graduation
	// Various graduations are provided <???>.
	Grade int8

	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
}

type CreateJudgmentOptions struct {
	Poll        *Poll
	Judge       *User
	Grade       int8
	CandidateID int64
}

type UpdateJudgmentOptions struct {
	Poll        *Poll
	Judge       *User
	Grade       int8
	CandidateID int64
}

type DeleteJudgmentOptions struct {
	Poll        *Poll
	Judge       *User
	CandidateID int64
}

func CreateJudgment(opts *CreateJudgmentOptions) (judgment *Judgment, err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	judgment, err = createJudgment(sess, opts)
	if err != nil {
		return nil, err
	}

	if err = sess.Commit(); err != nil {
		return nil, err
	}

	return judgment, nil
}

func createJudgment(e *xorm.Session, opts *CreateJudgmentOptions) (_ *Judgment, err error) {
	judgment := &Judgment{
		PollID:      opts.Poll.ID,
		Poll:        opts.Poll,
		JudgeID:     opts.Judge.ID,
		Judge:       opts.Judge,
		CandidateID: opts.CandidateID,
		Grade:       opts.Grade,
	}
	//e.Find()
	if _, err = e.Insert(judgment); err != nil {
		return nil, err
	}

	//if err = updatePollInfos(e, opts, poll); err != nil {
	//	return nil, err
	//}

	return judgment, nil
}

func getJudgmentByID(e Engine, id int64) (*Judgment, error) {
	repo := new(Judgment)
	has, err := e.ID(id).Get(repo)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrJudgmentNotFound{}
	}
	return repo, nil
}

func getJudgmentOfJudgeOnPollCandidate(e Engine, judgeID int64, pollID int64, candidateID int64) (judgment *Judgment, err error) {
	// We could probably use only one SQL query instead of two here.
	// No idea how this ORM works, and sprinting past it with snippet copy-pasting.
	judgmentsIds := make([]int64, 0, 1)
	if err = e.Table("judgment").
		Cols("id").
		Where("`judgment`.judge_id = ?", judgeID).
		And("`judgment`.poll_id = ?", pollID).
		And("`judgment`.candidate_id = ?", candidateID).
		Limit(1).
		Find(&judgmentsIds); err != nil {
		return nil, fmt.Errorf("Find Judgment: %v", err)
	}

	if 0 == len(judgmentsIds) {
		return nil, ErrJudgmentNotFound{}
	}

	judgment, errj := getJudgmentByID(e, judgmentsIds[0])
	if errj != nil {
		return nil, errj
	}

	return judgment, nil
}

func UpdateJudgment(opts *UpdateJudgmentOptions) (judgment *Judgment, err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	judgment, errJ := getJudgmentOfJudgeOnPollCandidate(sess, opts.Judge.ID, opts.Poll.ID, opts.CandidateID)
	if nil != errJ {
		return nil, errJ
	}

	judgment.Grade = opts.Grade

	_, err = sess.ID(judgment.ID).
		Cols("grade", "updated_unix").
		Update(judgment)
	if err != nil {
		return nil, err
	}

	if err = sess.Commit(); err != nil {
		return nil, err
	}

	return judgment, nil
}

func DeleteJudgment(opts *DeleteJudgmentOptions) (err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	judgment, errJ := getJudgmentOfJudgeOnPollCandidate(sess, opts.Judge.ID, opts.Poll.ID, opts.CandidateID)
	if nil != errJ {
		return errJ
	}

	if _, errD := sess.Delete(judgment); nil != errD {
		return errD
	}

	if err = sess.Commit(); nil != err {
		return err
	}

	return nil
}
