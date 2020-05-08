// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/setting"
	"fmt"
	"strings"

	//"code.gitea.io/gitea/modules/references"
	"code.gitea.io/gitea/modules/timeutil"
	//"fmt"
	"xorm.io/xorm"
)

// A Poll on <Subject> with the issues of a repository as candidates.
type Poll struct {
	ID       int64       `xorm:"pk autoincr"`
	RepoID   int64       `xorm:"INDEX"`
	Repo     *Repository `xorm:"-"`
	AuthorID int64       `xorm:"INDEX"`
	Author   *User       `xorm:"-"`
	//Index    int64       `xorm:"UNIQUE(repo_index)"` // Index in one repository.
	// When the poll is applied to all the issues, the subject should be an issue's trait.
	// eg: Quality, Importance, Urgency, Wholeness, Relevance…
	Subject string `xorm:"name"`
	// Description may be used to describe at length the constitutive details of that poll.
	// Eg: Rationale, Deliberation Consequences, Schedule, Modus Operandi…
	// It can be written in the usual gitea-flavored markdown.
	Description         string `xorm:"TEXT"`
	RenderedDescription string `xorm:"-"`
	Ref                 string // Do we need this?  Are we even using it?  WHat is it?

	Gradation           string `xorm:"-"`
	AreCandidatesIssues bool   // unused

	DeadlineUnix timeutil.TimeStamp `xorm:"INDEX"`
	CreatedUnix  timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix  timeutil.TimeStamp `xorm:"INDEX updated"`
	ClosedUnix   timeutil.TimeStamp `xorm:"INDEX"`

	// No idea how xorm works -- help!
	//Judgments         []*Judgment    `xorm:"-"`
	//Judgments         JudgmentList   `xorm:"-"`
}

// PollList is a list of polls offering additional functionality (perhaps)
type PollList []*Poll

func (poll *Poll) GetGradationList() []string {
	list := make([]string, 0, 6)

	// Placeholder until user customization somehow (poll.Gradation?)
	// - 🤮😒😐🙂😀🤩
	// - 😫😒😐😌😀😍  (more support, apparently)
	// - …
	list = append(list, "😫")
	list = append(list, "😒")
	list = append(list, "😐")
	list = append(list, "😌")
	list = append(list, "😀")
	list = append(list, "😍")

	return list
}

func (poll *Poll) GetCandidatesIDs() (_ []int64, err error) {
	ids := make([]int64, 0, 10)
	if err := x.Table("judgment").
		Select("DISTINCT candidate_id").
		Where("poll_id = ?", poll.ID).
		//And("updated_unix >= ?", start).
		//GroupBy("context_hash").
		//OrderBy("max( id ) desc").
		Find(&ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (poll *Poll) GetJudgmentOnCandidate(judge *User, candidateID int64) (judgmernt *Judgment) {
	judgment, err := getJudgmentOfJudgeOnPollCandidate(x, judge.ID, poll.ID, candidateID)
	if nil != err {
		return nil
	}

	return judgment
}

func (poll *Poll) GetResult() (results *PollResult) {
	// The deliberator should probably be a parameter of this function,
	// and upstream we could fetch it from context or settings.
	deliberator := &PollNaiveDeliberator{
		UseHighMean: false,
	}

	results, err := deliberator.Deliberate(poll)
	if nil != err {
		return nil // What should we do here?
	}

	return results
}

func (poll *Poll) CountGrades(candidateID int64, grade uint8) (_ uint64, err error) {
	rows := make([]int64, 0, 2)
	//amount := 0
	if err := x.Table("judgment").
		Select("COUNT(*) as amount").
		Where("poll_id = ?", poll.ID).
		And("candidate_id = ?", candidateID).
		And("grade = ?", grade).
		//GroupBy("context_hash").
		//OrderBy("max( id ) desc").
		Find(&rows); err != nil {
		return 0, err
	}
	if 1 != len(rows) {
		return 0, fmt.Errorf("wrong amount of COUNT()")
	}

	amount := uint64(rows[0])
	return amount, nil
}

// $ figlet -w 120 "Create"
//   ____                _
//  / ___|_ __ ___  __ _| |_ ___
// | |   | '__/ _ \/ _` | __/ _ \
// | |___| | |  __/ (_| | ||  __/
//  \____|_|  \___|\__,_|\__\___|
//

type CreatePollOptions struct {
	//Type  PollType  // for inline polls with their own candidates?
	Author      *User
	Repo        *Repository
	Subject     string
	Description string
	//Grades      string
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
		AuthorID:            opts.Author.ID,
		Author:              opts.Author,
		RepoID:              opts.Repo.ID,
		Repo:                opts.Repo,
		Subject:             opts.Subject,
		Description:         opts.Description,
		AreCandidatesIssues: true,
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

//  ____                _
// |  _ \ ___  __ _  __| |
// | |_) / _ \/ _` |/ _` |
// |  _ <  __/ (_| | (_| |
// |_| \_\___|\__,_|\__,_|
//

// GetPolls returns the (paginated) list of polls of a given repository and status.
func GetPolls(repoID int64, page int) (PollList, error) {
	polls := make([]*Poll, 0, setting.UI.IssuePagingNum)
	sess := x.Where("repo_id = ?", repoID)
	//sess := x.Where("repo_id = ? AND is_closed = ?", repoID, isClosed)
	if page > 0 {
		sess = sess.Limit(setting.UI.IssuePagingNum, (page-1)*setting.UI.IssuePagingNum)
	}

	return polls, sess.Find(&polls)
}

func getPollByRepoID(e Engine, repoID, id int64) (*Poll, error) {
	m := new(Poll)
	has, err := e.ID(id).Where("repo_id=?", repoID).Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrPollNotFound{ID: id, RepoID: repoID}
	}
	return m, nil
}

// GetPollByRepoID returns the poll in a repository.
func GetPollByRepoID(repoID, id int64) (*Poll, error) {
	return getPollByRepoID(x, repoID, id)
}

//  _   _           _       _
// | | | |_ __   __| | __ _| |_ ___
// | | | | '_ \ / _` |/ _` | __/ _ \
// | |_| | |_) | (_| | (_| | ||  __/
//  \___/| .__/ \__,_|\__,_|\__\___|
//       |_|

func updatePoll(e Engine, m *Poll) error {
	m.Subject = strings.TrimSpace(m.Subject)
	_, err := e.ID(m.ID).AllCols().
		// Do some extra work here, like updating stats?
		//SetExpr("num_closed_issues", builder.Select("count(*)").From("issue").Where(
		//	builder.Eq{
		//		"poll_id": m.ID,
		//		"is_closed":    true,
		//	},
		//)).
		Update(m)
	return err
}

// UpdatePoll updates information of given poll.
func UpdatePoll(m *Poll) error {
	sess := x.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := updatePoll(sess, m); err != nil {
		return err
	}

	//if err := updatePollCompleteness(sess, m.ID); err != nil {
	//	return err
	//}

	return sess.Commit()
}

//  ____       _      _
// |  _ \  ___| | ___| |_ ___
// | | | |/ _ \ |/ _ \ __/ _ \
// | |_| |  __/ |  __/ ||  __/
// |____/ \___|_|\___|\__\___|
//

// DeletePollByRepoID deletes a poll from a repository.
func DeletePollByRepoID(repoID, id int64) error {
	m, err := GetPollByRepoID(repoID, id)
	if err != nil {
		if IsErrPollNotFound(err) {
			return nil
		}
		return err
	}

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.ID(m.ID).Delete(new(Poll)); err != nil {
		return err
	}

	return sess.Commit()
}
