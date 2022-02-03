package opal

type Pair struct {
	Key   string
	Value string
	Count int
}

type Counter struct {
	data map[string]*Pair
}

/*
 * Construct a new counter object
 */
func NewCounter() Counter {
	return Counter{
		map[string]*Pair{},
	}
}

/*
 * Add a key-value pair to the counter, incrementing if the
 * key is already present
 */
func (count *Counter) Add(key string, value string) {
	if pair, ok := count.data[key]; ok {
		pair.Count += 1
	} else {
		count.data[key] = &Pair{
			Key:   key,
			Value: value,
			Count: 1,
		}
	}
}

/*
 * Enumerate values with duplicates
 */
func (count *Counter) Duplicates() []string {
	repeated := []string{}

	for _, data := range count.data {
		if data.Count > 1 {
			repeated = append(repeated, data.Value)
		}
	}

	return repeated
}
