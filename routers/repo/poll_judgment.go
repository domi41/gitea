// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/setting"
	"path"

	//"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/context"
	//"code.gitea.io/gitea/modules/setting"
	//"code.gitea.io/gitea/modules/timeutil"
	//"time"
)

type CreateJudgmentResponse struct {
	Judgment *models.Judgment
}

// Creates, Updates or Deletes a Judgment depending on the parameters
func EmitJudgment(ctx *context.Context) {
	judge := ctx.User

	grade := uint8(ctx.QueryInt("grade")) // 0 if not defined
	if grade < 1 {
		grade = 0
	}

	pollId := ctx.ParamsInt64(":id")
	candidateID := ctx.QueryInt64("candidate")

	poll, errP := models.GetPollByRepoID(ctx.Repo.Repository.ID, pollId)
	if nil != errP {
		ctx.NotFound("EmitJudgment.GetPollByRepoID", errP)
		return
	}

	judgment := poll.GetJudgmentOnCandidate(ctx.User, candidateID)
	if nil != judgment {

		if judgment.Grade == grade {
			// Delete a judgment if it exists and is the same as submitted.
			// Not obvious nor usual behavior, but pretty handy for now.
			errD := models.DeleteJudgment(&models.DeleteJudgmentOptions{
				Poll:        poll,
				Judge:       judge,
				CandidateID: candidateID,
			})
			if nil != errD {
				ctx.ServerError("EmitJudgment.DeleteJudgment", errD)
				return
			} else {
				ctx.Flash.Success(ctx.Tr("repo.polls.judgments.delete.success"))
			}

		} else {

			_, errU := models.UpdateJudgment(&models.UpdateJudgmentOptions{
				Poll:        poll,
				Judge:       judge,
				Grade:       grade,
				CandidateID: candidateID,
			})
			if nil != errU {
				ctx.ServerError("EmitJudgment.UpdateJudgment", errU)
				return
			} else {
				ctx.Flash.Success(ctx.Tr("repo.polls.judgments.update.success"))
			}

		}

	} else {

		_, errC := models.CreateJudgment(&models.CreateJudgmentOptions{
			Poll:        poll,
			Judge:       judge,
			Grade:       grade,
			CandidateID: candidateID,
		})
		if nil != errC {
			ctx.ServerError("EmitJudgment.EmitJudgment", errC)
			return
		} else {
			ctx.Flash.Success(ctx.Tr("repo.polls.judgments.create.success"))
		}

	}

	redirectPath := ctx.QueryTrim("redirect")
	ctx.Redirect(path.Join(setting.AppSubURL, redirectPath))

	// Nice, but yields too much (passwords, lol)
	//ctx.JSON(200, &CreateJudgmentResponse{
	//	Judgment: judgment,
	//})
}
