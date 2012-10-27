package extenders

import . "spexs"

type queryMap map[Token]*Query

//TODO: method queryMap.toQuerys
//TODO: method queryMap.addLoc(token, pos)
func toQuerys(queryMap queryMap) Querys {
	querys := make(Querys, len(queryMap))
	i := 0
	for _, q := range queryMap {
		querys[i] = q
		i += 1
	}
	return querys
}

func extend(base *Query, db *Database, querys queryMap) {
	for _, val := range base.Loc.Iter() {
		i, pos := DecodePos(val)

		token, ok, next := db.GetToken(i, pos)
		if !ok {
			continue
		}

		q, exists := querys[token]
		if !exists {
			q = NewQuery(base, RegToken{token, false, false})
			querys[token] = q
		}

		q.Loc.Add(EncodePos(i, next))
	}
}

func Simplex(base *Query) Querys {
	querys := make(queryMap)
	extend(base, base.Db, querys)
	return toQuerys(querys)
}

func combine(base *Query, db *Database, querys queryMap, isStar bool) {
	for _, group := range db.Groups {
		q := NewQuery(base, RegToken{group.Token, true, isStar})
		querys[group.Token] = q
		for _, token := range group.Elems {
			single, ok := querys[token]
			if ok {
				q.Loc.AddSet(single.Loc)
			}
		}
	}
}

func Groupex(base *Query) Querys {
	querys := make(queryMap)
	extend(base, base.Db, querys)
	combine(base, base.Db, querys, false)
	return toQuerys(querys)
}

func max(a uint, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func starExtend(base *Query, db *Database, querys queryMap) {
	lastPos := make(map[uint]uint, base.Loc.Len())

	for _, val := range base.Loc.Iter() {
		i, pos := DecodePos(val)
		lastPos[i] = max(lastPos[i], pos)
	}

	for i, last := range lastPos {
		var q *Query
		token, ok, next := db.GetToken(i, last)
		for ok {
			q, ok = querys[token]
			if !ok {
				q = NewQuery(base, RegToken{token, false, true})
				querys[token] = q
			}
			q.Loc.Add(EncodePos(i, next))
			token, ok, next = db.GetToken(i, next)
		}
	}
}

func Starex(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	stars := make(queryMap)
	starExtend(base, base.Db, stars)
	return append(toQuerys(patterns), toQuerys(stars)...)
}

func Regex(base *Query) Querys {
	patterns := make(queryMap)
	extend(base, base.Db, patterns)
	combine(base, base.Db, patterns, false)
	stars := make(queryMap)
	starExtend(base, base.Db, stars)
	combine(base, base.Db, stars, true)
	return append(toQuerys(patterns), toQuerys(stars)...)
}
