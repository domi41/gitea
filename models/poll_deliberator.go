// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"code.gitea.io/gitea/modules/timeutil"
	"fmt"
	"sort"
)

type PollDeliberator interface {
	Deliberate(poll *Poll) (result *PollResult, err error)
}

//  _   _       _
// | \ | | __ _(_)_   _____
// |  \| |/ _` | \ \ / / _ \
// | |\  | (_| | |\ V /  __/
// |_| \_|\__,_|_| \_/ \___|
//

type PollNaiveDeliberator struct {
	UseHighMean bool // should default to false ; strategy for even number of judgments
}

func (deli *PollNaiveDeliberator) Deliberate(poll *Poll) (_ *PollResult, err error) {

	naiveTallier := &PollNaiveTallier{}
	pollTally, err := naiveTallier.Tally(poll)
	if nil != err {
		return nil, err
	}

	candidates := make(PollCandidateResults, 0, 64)

	for _, candidateTally := range pollTally.Candidates {

		medianGrade := candidateTally.GetMedian()
		//candidateScore := candidateTally.GetScore()
		candidateScore := deli.GetScore(candidateTally)

		candidates = append(candidates, &PollCandidateResult{
			Poll:        poll,
			CandidateID: candidateTally.CandidateID,
			Position:    0, // We set it below after the Sort
			MedianGrade: medianGrade,
			Score:       candidateScore,
			Tally:       candidateTally,
			CreatedUnix: timeutil.TimeStampNow(),
		})

	}

	sort.Sort(sort.Reverse(candidates))

	previousScore := ""
	for key, candidate := range candidates {
		position := uint64(key + 1)
		if (previousScore == candidate.Score) && (key > 0) {
			position = candidates[key-1].Position
		}
		candidate.Position = position
		previousScore = candidate.Score
	}

	result := &PollResult{
		Poll:        poll,
		Tally:       pollTally,
		Candidates:  candidates,
		CreatedUnix: timeutil.TimeStampNow(),
	}

	return result, nil
}

//  ____                 _
// / ___|  ___ ___  _ __(_)_ __   __ _
// \___ \ / __/ _ \| '__| | '_ \ / _` |
//  ___) | (_| (_) | |  | | | | | (_| |
// |____/ \___\___/|_|  |_|_| |_|\__, |
//                               |___/

func (deli *PollNaiveDeliberator) GetScore(pct *PollCandidateTally) (_ string) {
	score := ""

	//ct := *pct
	// FIXME: Not THAT naive :)
	score += fmt.Sprintf("%03d", pct.GetMedian())

	//println("Score: "+score)

	return score
}

/*
// Assume that each candidate has the same amount of judgments = MAX_JUDGES.
// (best fill with 0=REJECT to allow posterior candidate addition, cf. <paper>)
for each Candidate
	ct = CandidateTally(Candidate) // sums of judgments, per grade, basically
	score = "" // score is a string but could be raw bits
	// When we append integers to score below,
	// consider that we concatenate the string representation including leading zeroes
	// up to the amount of digits we need to store 2 * MAX_JUDGES,
	// or the raw bits (unsigned and led as well) in case of a byte array.

	for i in range(MAX_GRADES)
		grade = ct.median()
		score.append(grade) // three digits will suffice for int8
		// Collect biggest of the two groups outside of the median.
		// Group Grade is the group"s grade adjacent to the median group
		// Group Sign is:
		// - +1 if the group promotes higher grades (adhesion)
		// - -1 if the group promotes lower grades (contestation)
		// - ±0 if there is no spoon
		group_size, group_sign, group_grade = ct.get_biggest_group()
		// MAX_JUDGES is to deal with negative values lexicographically
		score.append(MAX_JUDGES + groups_sign * group_size)
		// Move the median grades into the group grades
		ct.regrade_judgments(grade, groups_grade)

	// Use it later in a bubble sort or whatever
	Candidate.score = score

*/
