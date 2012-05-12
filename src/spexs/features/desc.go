package features

import (
	. "spexs"
)

type Func func(*Pattern, *Reference) float64

type Desc struct {
	Name string
	Desc string
	Func Func
}

func Get(name string) (*Desc, bool) {
	for _, e := range All {
		if e.Name == name {
			return &e, true
		}
	}
	return nil, false
}

type StrFunc func(*Pattern, *Reference) string

type StrDesc struct {
	Name string
	Desc string
	Func StrFunc
}

func GetStr(name string) (*StrDesc, bool) {
	for _, e := range Str {
		if e.Name == name {
			return &e, true
		}
	}
	return nil, false
}