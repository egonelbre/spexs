package main

import (
	"io"
	. "spexs"
	pool "spexs/pool"
)

const MAX_POOL_SIZE = 1024 * 1024 * 1024

type TrieFitnessCreator func(interface{}) FitnessFunc
type TrieExtenderCreator func(interface{}) ExtenderFunc

type PrinterFunc func(io.Writer, *Pattern, *Reference)

type AppSetup struct {
	Setup
	Fitness FitnessFunc
	Printer PrinterFunc
}

func lengthFitness(p *Pattern) float64 {
	return 1 / float64(p.Len())
}

type PatternFilterCreator func(limit int) FilterFunc

func CreateInput(conf Conf, setup AppSetup) Pooler {
	//in := NewPriorityPool(lengthFitness, MAX_POOL_SIZE, true)
	in := pool.NewLifo()
	in.Put(NewFullPattern(setup.Ref))
	return in
}

func CreateOutput(conf Conf, setup AppSetup, f FitnessFunc) Pooler {
	if conf.Output.Queue == "lifo" {
		return pool.NewLifo()
	}
	size := conf.Output.Count
	if size < 0 {
		size = MAX_POOL_SIZE
	}
	return pool.NewPriority(f, size, conf.Output.Sort == "asc")
}

func CreateSetup(conf Conf) AppSetup {
	var s AppSetup
	s.Ref = CreateReference(conf)

	s.In = CreateInput(conf, s)
	s.Fitness = CreateFitness(conf, s)
	s.Out = CreateOutput(conf, s, s.Fitness)

	s.Extender = CreateExtender(conf, s)
	s.Extendable = CreateFilter(conf.Extension.Filter, s)
	s.Outputtable = CreateFilter(conf.Output.Filter, s)

	s.PostProcess = func(p *Pattern, s *Setup) error {
		p.Occs(s.Ref,0)
		p.Count(s.Ref,0)
		p.Pos.Clear()
		return nil
	}

	s.Printer = CreatePrinter(conf, s)

	return s
}

func initSetup() {
	initFitnesses()
	initFilters()
}
