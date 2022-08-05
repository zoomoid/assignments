package util

type empty struct{}

type Set map[string]empty

func NewSet(items ...string) Set {
	ss := Set{}
	ss.Insert(items...)
	return ss
}

func (s Set) Insert(items ...string) Set {
	for _, item := range items {
		s[item] = empty{}
	}
	return s
}

func (s Set) Delete(items ...string) Set {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s Set) Has(item string) bool {
	_, contained := s[item]
	return contained
}

func (s Set) Len() int {
	return len(s)
}
