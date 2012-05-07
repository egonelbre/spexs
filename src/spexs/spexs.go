package spexs

import "time"

const (
	patternsBufferSize = 128
)

type Char byte

type Patterns chan *Pattern

type Pooler interface {
	Take() (*Pattern, bool)
	Put(*Pattern)
	Len() int
}

type ExtenderFunc func(p *Pattern, ref *Reference) Patterns
type FilterFunc func(p *Pattern, ref *Reference) bool
type FitnessFunc func(p *Pattern) float64
type PostProcessFunc func(p *Pattern, s *Setup) error

type Setup struct {
	Ref *Reference
	Out Pooler
	In  Pooler

	Extender ExtenderFunc

	Extendable  FilterFunc
	Outputtable FilterFunc

	PostProcess PostProcessFunc
}

func NewPatterns() Patterns {
	return make(Patterns, patternsBufferSize)
}

func Run(s *Setup) {
	for {
		p, valid := s.In.Take()
		if !valid {
			return
		}

		extensions := s.Extender(p, s.Ref)
		for extended := range extensions {
			if s.Extendable(extended, s.Ref) {
				s.In.Put(extended)
			}
			if s.Outputtable(extended, s.Ref) {
				s.Out.Put(extended)
			}
		}

		if s.PostProcess(p, s) != nil {
			break
		}
	}
}

func Parallel(f func(), routines int) {
	stop := make(chan int, routines)
	for i := 0; i < routines; i += 1 {
		go func(rtn int) {
			defer func() { stop <- 1 }()
			f()
		}(i)
	}
	for i := 0; i < routines; i += 1 {
		<-stop
	}
}

func RunParallel(s *Setup, routines int) {
	quit := make(chan int, routines)
	counter := make(chan int, routines)

	for i := 0; i < routines; i += 1 {
		go func(rtn int){
				main:for {
					p, valid := s.In.Take()
					for !valid {
						p, valid = s.In.Take()
						if valid {
							break
						}
						counter <- -1
						select {
						case <-time.After(100*time.Millisecond):
						case <-quit: 
							break main
						}
						counter <- 1
					}

					extensions := s.Extender(p, s.Ref)
					for extended := range extensions {
						if s.Extendable(extended, s.Ref) {
							s.In.Put(extended)
						}
						if s.Outputtable(extended, s.Ref) {
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
		value := <- counter
		count += value
	}

	for i := 0; i < routines; i += 1 {
		quit <- 1
	}
}