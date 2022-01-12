package opal

type Set struct {
	data map[string]bool
}

func NewSet(elems []string) *Set {
	data := map[string]bool{}

	for _, elem := range elems {
		data[elem] = true
	}

	return &Set{data}
}

func (set *Set) Has(elem string) bool {
	_, ok := set.data[elem]

	return ok
}
