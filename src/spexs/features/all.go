package features

import (
	. "spexs"
	"stats/hyper"
)

var All = [...]Desc{
	{"query-seqs",
		"total number of query sequences",
		func(q *Query) float64 {
			return float64(q.Db.Sections[0].Count)
		}},
	{"back-seqs",
		"total number of background sequences",
		func(q *Query) float64 {
			return float64(q.Db.Sections[1].Count)
		}},
	{"query-match-seqs",
		"number of matching query sequences",
		func(q *Query) float64 {
			return float64(q.MatchSeqs()[0])
		}},
	{"back-match-seqs",
		"number of matching background sequences",
		func(q *Query) float64 {
			return float64(q.MatchSeqs()[1])
		}},
	{"query-match-occs",
		"number of occurences in query",
		func(q *Query) float64 {
			return float64(q.MatchOccs()[0])
		}},
	{"back-match-occs",
		"number of occurences in background",
		func(q *Query) float64 {
			return float64(q.MatchOccs()[1])
		}},
	{"query-match-seqs-prop",
		"percentage of matching sequences in query",
		func(q *Query) float64 {
			seqs := q.MatchSeqs()[0]
			total := q.Db.Sections[0].Count
			return float64(seqs) / float64(total)
		}},
	{"back-match-seqs-prop",
		"percentage of matching sequences in background",
		func(q *Query) float64 {
			seqs := q.MatchSeqs()[1]
			total := q.Db.Sections[1].Count
			return float64(seqs) / float64(total)
		}},

	{"match-hyper-up-pvalue",
		"hypergeometric split q-value",
		func(q *Query) float64 {
			seqs := q.MatchSeqs()
			pvalue := hyper.Split(seqs[0], seqs[1],
				q.Db.Sections[0].Count, q.Db.Sections[1].Count)
			return pvalue
		}},
	{"match-hyper-up-pvalue-approx",
		"approximate hypergeometric split q-value (~5 significant digits)",
		func(q *Query) float64 {
			seqs := q.MatchSeqs()
			pvalue := hyper.SplitApprox(seqs[0], seqs[1],
				q.Db.Sections[0].Count, q.Db.Sections[1].Count)
			return pvalue
		}},
	{"match-hyper-down-pvalue",
		"hypergeometric split q-value down",
		func(q *Query) float64 {
			seqs := q.MatchSeqs()
			pvalue := hyper.SplitDown(seqs[0], seqs[1],
				q.Db.Sections[0].Count, q.Db.Sections[1].Count)
			return pvalue
		}},
	{"match-ratio",
		"ratio of (matches in query + 1) / (matches in background + 1)",
		func(q *Query) float64 {
			seqs := q.MatchSeqs()
			return float64(seqs[0]+1) / float64(seqs[1]+1)
		}},

	{"match-hyper-optimal-pvalue",
		"finds optimal hypergeometric p-value for the input",
		func(q *Query) float64 {
			return q.FindOptimalSplit()
		}},
	{"match-hyper-optimal-seqs",
		"how many sequences were in optimal hypergeometric",
		func(q *Query) float64 {
			return float64(q.FindOptimalSplitSeqs())
		}},
	{"match-hyper-optimal-matches",
		"how many matches were in optimal hypergeometric",
		func(q *Query) float64 {
			return float64(q.FindOptimalSplitMatches())
		}},

	{"pat-length",
		"length of the pattern",
		func(q *Query) float64 {
			t := 0
			for _, e := range q.Pat {
				t += 1
				if e.IsStar {
					t += 1
				}
			}
			return float64(t)
		}},
	{"pat-chars",
		"count of characters in pattern",
		func(q *Query) float64 {
			t := 0
			for _, e := range q.Pat {
				if !e.IsGroup {
					t += 1
				}
			}
			return float64(t)
		}},
	{"pat-groups",
		"count of groups in pattern",
		func(q *Query) float64 {
			t := 0
			for _, e := range q.Pat {
				if e.IsGroup {
					t += 1
				}
			}
			return float64(t)
		}},
	{"pat-stars",
		"count of stars in pattern",
		func(q *Query) float64 {
			t := 0
			for _, e := range q.Pat {
				if e.IsStar {
					t += 1
				}
			}
			return float64(t)
		}},
}
