// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

type PollDeliberator interface {
	Deliberate(poll *Poll) (result *PollResult, err error)
}

type PollNaiveDeliberator struct {
	UseHighMean bool // should default to false ; strategy for even number of judgments
}

func (deli *PollNaiveDeliberator) Deliberate(poll *Poll) (_ *PollResult, err error) {

	c := make([]*PollCandidateResult, 0, 100)

	c = append(c, &PollCandidateResult{
		Poll:        poll,
		CandidateID: 0,
		Position:    0,
		MedianGrade: 5,
		//Tally:       *PollCandidateTally
		//CreatedUnix: timeutil.TimeStamp
	})

	result := &PollResult{
		Poll:        poll,
		Tally:       nil, // FIXME
		Candidates:  c,   // FIXME
		CreatedUnix: 0,   // FIXME
	}

	return result, nil
}
