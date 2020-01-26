package parse

import (
	"testing"
	"reflect"
)

func TestMakeTagMap (t *testing.T) {
	tests := map[string]struct {
		entries []Entry
		want map[string][]int
	}{
		"ex1": {[]Entry{{ID: 1, Tags: []string{"ex1"}}}, map[string][]int{"ex1": []int{1}}},
		"ex2": {[]Entry{{ID: 1, Tags: []string{"ex2"}},{ID:2, Tags: []string{"ex2"}}}, map[string][]int{"ex2": []int{1, 2}}},
		"ex3": {[]Entry{{ID: 501, Tags: []string{"ex2"}},{ID:2, Tags: []string{"ex2"}},{ID: 9, Tags: []string{"ex2"}},{ID: 10, Tags: []string{"ex2"}},}, map[string][]int{"ex2": []int{2, 9, 10, 501}}},
		"ex4": {[]Entry{{ID: 508, Tags: []string{"ex2"}},{ID:22, Tags: []string{"ex2"}},{ID: 1, Tags: []string{"ex2"}},{ID: 99, Tags: []string{"ex2"}},{ID: 10, Tags: []string{"ex2"}},}, map[string][]int{"ex2": []int{1, 10, 22, 99, 508}}},

	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MakeTagMap(tc.entries)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.entries)
			}
		})
	}
}

func TestMakeTags (t *testing.T) {
	tests := map[string]struct {
		tagIDs map[string][]int
		want []Tags
	}{
		"ex1": {map[string][]int{"ex1": []int{1}}, []Tags{{"ex1", 1}}},
		"ex2": {map[string][]int{"ex2": []int{1, 2}}, []Tags{{"ex2", 2}}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MakeTags(tc.tagIDs)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.tagIDs)
			}
		})
	}
}

func TestSortTags (t *testing.T) {
	tests := map[string]struct {
		t []Tags
		want []Tags
	}{
		"ex1": {[]Tags{{"A", 1}, {"Z", 2}, {"C", 3}}, []Tags{{"A", 1}, {"C", 3}, {"Z", 2}}},
		"ex2": {[]Tags{{"U", 1}, {"B", 2}, {"P", 3}}, []Tags{{"B", 2}, {"P", 3}, {"U", 1}}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := sortTags(tc.t)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.t)
			}
		})
	}
}

func TestSortIDs (t *testing.T) {
	tests := map[string]struct {
		IDs map[string][]int
		want map[string][]int
	}{
		"ex1": {map[string][]int{"ex1": []int{7, 1, 8, 4, 33, 67, 2}}, map[string][]int{"ex1": []int{1, 2, 4, 7, 8, 33, 67}}},
		"ex2": {map[string][]int{"ex1": []int{458, 7, 1, 52, 8, 4, 6, 33, 67, 2}}, map[string][]int{"ex1": []int{1, 2, 4, 6, 7, 8, 33, 52, 67, 458}}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := SortIDs(tc.IDs)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.IDs)
			}
		})
	}
}