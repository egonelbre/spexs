package main

import (
	. "spexs/trie"
	"text/template"
	"io"
	"log"
	"fmt"
)

type printerArgs struct{
	Regexp string
	Fitness float64
	Length int
	Count int
	PValue float64
}

func CreatePrinter(conf Conf, setup AppSetup) PrinterFunc {
	tmpl, err := template.New("").Parse(conf.Output.Format)
	if err != nil {
		log.Println("Unable to create template based on output format.")
		log.Fatal(err)
	}

	f := func(out io.Writer, pat *Pattern, ref *Reference) {
		if pat == nil {
			fmt.Println(conf.Output.Format)
			return
		}

		node := printerArgs{
			Regexp : setup.Ref.ReplaceGroups(pat.String()),
			PValue : pat.PValue(ref),
			Fitness : setup.Fitness(pat),
			Length : pat.Len(),
			Count : pat.Pos.Len(),
		}

		err = tmpl.Execute(out, node)
		if err != nil {
			log.Println("Unable to output pattern.")
			log.Fatal(err)
		}
	}

	//TODO: test printer

	return f
}