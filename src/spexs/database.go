package spexs

import (
	"bytes"
	"set"
)

type Sequence struct {
	start, end int
	Count  map[int]int
}

type TokenGroup struct {
	Token    Token
	Elems    []Token
	Name     string
	FullName string
}

type TokenInfo struct {
	Token Token
	Name  string
	Count int
}

type Database struct {
	Alphabet map[Token]*TokenInfo
	Groups   map[Token]*TokenGroup

	PosToSequence []int
	FullSequence []Token
	Sequences   []*Sequence
	Total       []int // total number sequence for each section
	TotalTokens []int // total number of tokens for each section

	Separator string // separator for joining pattern

	nameToToken map[string]Token
	strToSeq    map[string]int
	lastToken   Token
}

const nullToken = Token(0)

func NewDatabase(estimatedSize int) *Database {
	return &Database{
		Alphabet: make(map[Token]*TokenInfo),
		Groups:   make(map[Token]*TokenGroup),

		PosToSequence: make([]int, 0, estimatedSize),
		FullSequence: make([]Token, 0, estimatedSize),
		Sequences:   make([]*Sequence, 0, estimatedSize),
		Total:       make([]int, 0),
		TotalTokens: make([]int, 0),

		Separator: "",

		strToSeq:    make(map[string]int),
		nameToToken: make(map[string]Token),
		lastToken:   Token(1),
	}
}

func (db *Database) AddAllPositions(s set.Set) {
	for i, v := range db.FullSequence {
		if v != nullToken {
			s.Add(i)
		}
	}
}

func (db *Database) GetToken(pos int) (token Token, ok bool, nextPos int) {
	t := db.FullSequence[pos]
	if t == nullToken {
		return 0, false, 0
	}
	return t, true, pos + 1
}

func (db *Database) mkNewToken() Token {
	newToken := db.lastToken
	db.lastToken += 1
	return newToken
}

func (db *Database) Matches(s set.Set) []int {
	counted := make(map[int]bool, s.Len())
	count := make([]int, len(db.Total))

	for _, p := range s.Iter() {
		si := db.PosToSequence[p]
		if counted[si] {
			continue
		}
		counted[si] = true

		seq := db.Sequences[si]
		for sec, c := range seq.Count {
			count[sec] += c
		}
	}

	return count
}

func (db *Database) Occs(s set.Set) []int {
	count := make([]int, len(db.Total))
	for _, p := range s.Iter() {
		si := db.PosToSequence[p]
		seq := db.Sequences[si]
		for sec, c := range seq.Count {
			count[sec] += c
		}
	}
	return count
}

func (db *Database) Seqs(s set.Set) []int {
	counted := make(map[int]bool, 30)
	count := make([]int, len(db.Total))
	for _, p := range s.Iter() {
		si := db.PosToSequence[p]
		if counted[si] {
			continue
		}
		counted[si] = true
		seq := db.Sequences[si]
		for sec, _ := range seq.Count {
			count[sec] += 1
		}
	}
	return count
}

func (db *Database) AddGroup(group *TokenGroup) Token {
	token := db.mkNewToken()
	group.Token = token
	db.Groups[token] = group
	return token
}

func (db *Database) AddToken(tokenName string) Token {
	token := db.mkNewToken()
	db.nameToToken[tokenName] = token
	db.Alphabet[token] = &TokenInfo{token, tokenName, 0}
	return token
}

func (db *Database) ToTokens(raw []string) []Token {
	tokens := make([]Token, len(raw))
	for i, name := range raw {
		token, ok := db.nameToToken[name]
		if !ok {
			token = db.AddToken(name)
		}
		tokens[i] = token
	}
	return tokens
}

func (db *Database) MakeSection() int {
	db.Total = append(db.Total, 0)
	db.TotalTokens = append(db.TotalTokens, 0)
	return len(db.Total) - 1
}

func sum(count []int) int {
	total := 0
	for _, val := range count {
		total += val
	}
	return total
}

func (db *Database) addToTokenCount(sec int, tokens []Token, count int) {
	db.TotalTokens[sec] += count * len(tokens)
	for _, t := range tokens {
		db.Alphabet[t].Count += count
	}
}

func (db *Database) AddSequence(sec int, raw []string, count int) {
	db.Total[sec] += count
	tokens := db.ToTokens(raw)
	db.addToTokenCount(sec, tokens, count)

	hash := hashTokens(tokens)
	if si, ok := db.strToSeq[hash]; ok {
		seq := db.Sequences[si]
		seq.Count[sec] += count
		return
	}
	
	seq := &Sequence{0, 0, make(map[int]int)}
	seq.Count[sec] = count
	
	// add sequence tokens to a single array
	seq.start = len(db.FullSequence)
	db.FullSequence = append(db.FullSequence, tokens...)
	seq.end = len(db.FullSequence)
	// add separator
	db.FullSequence = append(db.FullSequence, nullToken)
	
	// add sequence info to sequence list
	db.Sequences = append(db.Sequences, seq)
	// cache the sequence index
	si := len(db.Sequences) - 1
	db.strToSeq[hash] = si

	// add sequence indexes for positions
	db.PosToSequence = append(db.PosToSequence, make([]int, len(tokens) + 1)...)
	for i := seq.start; i < seq.end; i += 1 {
		db.PosToSequence[i] = si
	}
	db.PosToSequence[seq.end] = 0
}

func hashTokens(toks []Token) string {
	var buf bytes.Buffer
	for _, t := range toks {
		buf.WriteRune(rune(t))
	}
	return string(buf.Bytes())
}