package spexs

import "time"

type Token uint
type Querys []*Query

type Pooler interface {
	Take() (*Query, bool)
	Put(*Query)
	Len() int
}

type ExtenderFunc func(p *Query) Querys
type FilterFunc func(p *Query) bool
type FitnessFunc func(p *Query) float64
type FeatureFunc func(p *Query) float64, string
type PostProcessFunc func(p *Query, s *Setup) error

type Setup struct {
	Db  *Database
	Out Pooler
	In  Pooler

	Extender ExtenderFunc

	Extendable  FilterFunc
	Outputtable FilterFunc

	PostProcess PostProcessFunc
}

func prepareSpexs(s *Setup) {
	maxSeq := 0
	for _, seq := range s.Db.Sequences {
		if seq.Len > maxSeq {
			maxSeq = seq.Len
		}
	}

	for i := uint(0); i < 32; i += 1 {
		if 1<<i > maxSeq {
			PosOffset = i
			break
		}
	}

	s.In.Put(NewEmptyQuery(s.Db))
}

func Run(s *Setup) {
	prepareSpexs(s)
	for {
		p, valid := s.In.Take()
		if !valid {
			return
		}

		extensions := s.Extender(p)
		for _, extended := range extensions {
			if s.Extendable(extended) {
				s.In.Put(extended)
			}
			if s.Outputtable(extended) {
				s.Out.Put(extended)
			}
		}

		if s.PostProcess(p, s) != nil {
			break
		}
	}
}

func RunParallel(s *Setup, routines int) {
	prepareSpexs(s)

	quit := make(chan int, routines)
	counter := make(chan int, routines)

	for i := 0; i < routines; i += 1 {
		go func(rtn int) {
		main:
			for {
				p, valid := s.In.Take()
				for !valid {
					p, valid = s.In.Take()
					if valid {
						break
					}
					counter <- -1
					select {
					case <-time.After(100 * time.Millisecond):
					case <-quit:
						break main
					}
					counter <- 1
				}

				extensions := s.Extender(p)
				for _, extended := range extensions {
					if s.Extendable(extended) {
						s.In.Put(extended)
					}
					if s.Outputtable(extended) {
						s.Out.Put(extended)
					}
				}

				if s.PostProcess(p, s) != nil {
					break
				}
			}
		}(i)
	}

	count := routines
	for count > 0 {
		value := <-counter
		count += value
	}

	for i := 0; i < routines; i += 1 {
		quit <- 1
	}
}
