package parse

import (
	"testing"
	"reflect"
)

func TestMakeCatMaps (t *testing.T) {
	tests := map[string]struct {
		entries []Entry
		want map[string][]int
	}{
		"ex1": {[]Entry{{ID: 1, Cat: "ex1"}}, map[string][]int{"ex1": []int{1}}},
		"ex2": {[]Entry{{ID: 2, Cat: "ex2"},{ID:1, Cat: "ex2"}}, map[string][]int{"ex2": []int{1, 2}}},

	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MakeCatMap(tc.entries)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.entries)
			}
		})
	}
}

func TestMakeCatList (t *testing.T) {
	tests := map[string]struct {
		catIDs map[string][]int
		want []Cat
	}{
		"ex1": {map[string][]int{"ex1": []int{1}}, []Cat{{"ex1", 1}}},
		"ex2": {map[string][]int{"ex2": []int{1, 2}}, []Cat{{"ex2", 2}}},

	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MakeCatList(tc.catIDs)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.catIDs)
			}
		})
	}
}
