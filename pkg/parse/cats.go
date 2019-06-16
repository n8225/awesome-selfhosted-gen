package parse

import (
	"sort"
)

//Cats struct
type Cats struct {
	Cat   string
	Count int
}

// MakeCats puts unique Categories and a count into the Cats struct.
func MakeCats(entries []Entry) []Cats {
	catsl := []Cats{}
	c := new(Cats)
	var tmp []string
	for _, e := range entries {
		tmp = append(tmp, e.Cat)
	}
	catMap := make(map[string]int)
	for _, item := range tmp {
		_, exist := catMap[item]
		if exist {
			catMap[item]++
		} else {
			catMap[item] = 1
		}
	}
	for k, v := range catMap {
		c.Cat = k
		c.Count = v
		catsl = append(catsl, *c)
	}
	sort.Slice(catsl, func(i, j int) bool {
		return catsl[i].Cat < catsl[j].Cat
	})
	return catsl
}
