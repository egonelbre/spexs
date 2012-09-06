package main

import (
	"debugger"
	"fmt"
	. "spexs"
)

var dbg = debugger.New()

func attachDebugger(s *AppSetup) {
	debugger.RunShell(dbg)
	f := s.Extender
	s.Extender = func(q *Query, db *Database) Querys {
		result := f(q, db)
		dbg.Break(func() {
			fmt.Fprintf(dbg.Logout, "extending: %v\n", q.String(db, true))
			for _, extended := range result {
				fmt.Fprintf(dbg.Logout, "  | %v\n", q.String(db, true))
				if *verbose {
					fmt.Fprintf(dbg.Logout, "\tE:%v\tO:%v\n",
						s.Extendable(extended, db),
						s.Outputtable(extended, db))
					fmt.Fprintf(dbg.Logout, "      ")
					s.Printer(dbg.Logout, extended, db)
				}
			}
		})
		return result
	}
}
