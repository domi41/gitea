// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/auth"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/context"
	//"code.gitea.io/gitea/modules/setting"
	//"code.gitea.io/gitea/modules/timeutil"
	//"time"
)

const (
	tplPollsIndex base.TplName = "repo/polls/polls_index"
	tplPollsNew   base.TplName = "repo/polls/polls_new"
)

// IndexPolls renders an index of all the polls
func IndexPolls(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.polls.index.title")
	//ctx.Data["PageIsMilestones"] = true

	page := ctx.QueryInt("page") // 0 if not defined ?
	if page <= 1 {
		page = 1
	}

	polls, err := models.GetPolls(ctx.Repo.Repository.ID, page)
	if err != nil {
		ctx.ServerError("GetPolls", err)
		return
	}

	ctx.Data["Polls"] = polls

	//pager := context.NewPagination(total, setting.UI.IssuePagingNum, page, 5)
	//pager.AddParam(ctx, "state", "State")
	//ctx.Data["Page"] = pager

	ctx.HTML(200, tplPollsIndex)
}

// NewPoll renders the "new poll" page with its form
func NewPoll(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.polls.new")

	//ctx.Data["DateLang"] = setting.DateLang(ctx.Locale.Language())

	ctx.HTML(200, tplPollsNew)
}

// NewPollPost processes the "new poll" form and redirects
func NewPollPost(ctx *context.Context, form auth.CreatePollForm) {
	ctx.Data["Title"] = ctx.Tr("repo.polls.new")
	//ctx.Data["DateLang"] = setting.DateLang(ctx.Locale.Language())

	if ctx.HasError() {
		ctx.HTML(200, tplPollsNew)
		return
	}

	if _, err := models.CreatePoll(&models.CreatePollOptions{
		Author:      ctx.User,
		Repo:        ctx.Repo.Repository,
		Subject:     form.Subject,
		Description: form.Description,
	}); err != nil {
		ctx.ServerError("CreatePoll", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("repo.polls.create_success", form.Subject))
	ctx.Redirect(ctx.Repo.RepoLink + "/polls")
}

// EditPoll renders editing poll page
func EditPoll(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.polls.edit")
	//ctx.Data["PageIsPolls"] = true
	ctx.Data["PageIsEditPoll"] = true
	//ctx.Data["DateLang"] = setting.DateLang(ctx.Locale.Language())

	m, err := models.GetPollByRepoID(ctx.Repo.Repository.ID, ctx.ParamsInt64(":id"))
	if err != nil {
		if models.IsErrPollNotFound(err) {
			ctx.NotFound("", nil)
		} else {
			ctx.ServerError("GetPollByRepoID", err)
		}
		return
	}
	ctx.Data["subject"] = m.Subject
	ctx.Data["description"] = m.Description
	ctx.HTML(200, tplPollsNew)
}

// EditPollPost response for edting poll
func EditPollPost(ctx *context.Context, form auth.CreatePollForm) {
	ctx.Data["Title"] = ctx.Tr("repo.polls.edit")
	//ctx.Data["PageIsPolls"] = true
	ctx.Data["PageIsEditPoll"] = true
	//ctx.Data["DateLang"] = setting.DateLang(ctx.Locale.Language())

	if ctx.HasError() {
		ctx.HTML(200, tplPollsNew)
		return
	}

	//if len(form.Deadline) == 0 {
	//	form.Deadline = "9999-12-31"
	//}
	//deadline, err := time.ParseInLocation("2006-01-02", form.Deadline, time.Local)
	//if err != nil {
	//	ctx.Data["Err_Deadline"] = true
	//	ctx.RenderWithErr(ctx.Tr("repo.polls.invalid_due_date_format"), tplPollNew, &form)
	//	return
	//}

	//deadline = time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 23, 59, 59, 0, deadline.Location())
	m, err := models.GetPollByRepoID(ctx.Repo.Repository.ID, ctx.ParamsInt64(":id"))
	if err != nil {
		if models.IsErrPollNotFound(err) {
			ctx.NotFound("", nil)
		} else {
			ctx.ServerError("GetPollByRepoID", err)
		}
		return
	}

	m.Subject = form.Subject
	m.Description = form.Description

	if err = models.UpdatePoll(m); err != nil {
		ctx.ServerError("UpdatePoll", err)
		return
	}

	ctx.Flash.Success(ctx.Tr("repo.polls.edit_success", m.Subject))
	ctx.Redirect(ctx.Repo.RepoLink + "/polls")
}

// DeletePoll delete a poll and redirects
func DeletePoll(ctx *context.Context) {
	if err := models.DeletePollByRepoID(ctx.Repo.Repository.ID, ctx.ParamsInt64(":id")); err != nil {
		ctx.Flash.Error("DeletePollByRepoID: " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("repo.polls.deletion_success"))
	}

	ctx.JSON(200, map[string]interface{}{
		"redirect": ctx.Repo.RepoLink + "/polls",
	})
}
