package parse

import (
	"sort"
)

type Cat struct {
	Cat   string `json:"Cat"`
	Count int `json:"Count"`
}


// MakeCatMap puts Categories and IDs into the Cat map.
func MakeCatMap(entries []Entry) map[string][]int {
	catIDs := make(map[string][]int)
	for _, e := range entries {
		catIDs[e.Cat] = append(catIDs[e.Cat], e.ID)
	}
	return SortIDs(catIDs)
}

// MakeCatList sorts Categories and adds the count into a slice.
func MakeCatList(catIDs map[string][]int) []Cat {
	cats := []Cat{}
	for k, v := range catIDs {
		c := Cat{k, len(v)}
		cats = append(cats, c)
	}
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].Cat < cats[j].Cat
	})

	return cats
}