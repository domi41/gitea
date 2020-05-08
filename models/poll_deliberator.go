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
