package opal

/*
 * Set data-structure
 */
type Set struct {
	data map[string]bool
}

/*
 *Construct a new-set from a slice
 */
func NewSet(elems []string) *Set {
	data := map[string]bool{}

	for _, elem := range elems {
		data[elem] = true
	}

	return &Set{data}
}

/*
 * Does the element have the set?
 */
func (set *Set) Has(elem string) bool {
	_, ok := set.data[elem]

	return ok
}

/*
 * Add element to the set
 */
func (set *Set) Add(elem string) {
	set.data[elem] = true
}

/*
 * Get set cardinality
 */
func (set *Set) Size() int {
	return len(set.data)
}
