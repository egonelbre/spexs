package extenders

import (
	. "github.com/egonelbre/spexs2/search"
	"github.com/egonelbre/spexs2/utils"
)

type CreateFunc func(Setup, []byte) Extender

var All = [...]Extender{
	Simple,
	Group,
	Star,
	StarGreedy,
	Regex,
	RegexGreedy,
}

func wrap(f Extender) CreateFunc {
	return func(s Setup, data []byte) Extender {
		return f
	}
}

func Get(name string) (Extender, bool) {
	for _, fn := range All {
		if utils.FuncName(fn) == name {
			return fn, true
		}
	}
	return nil, false
}

var Help = `
  Simple : uses the sequence tokens to discover the patterns  ( ACCT )
  Group : uses additionally defined groups in Alphabet.Groups ( AC[CT]T )
  Star : uses matching anything in the pattern ( AC.*T )
  StarGreedy : matches greedily anything in the pattern ( AC.*T )
  Regex : uses both group and star token in the pattern ( A[CT].*T )
  RegexGreedy : uses both group and star token in the pattern ( A[CT].*T )
`
