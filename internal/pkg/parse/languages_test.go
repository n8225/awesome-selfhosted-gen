package parse

import (
	"testing"
	"reflect"
)

func TestMakeLangs (t *testing.T) {
	tests := map[string]struct {
		entries []Entry
		want map[string][]int
	}{
		"ex1": {[]Entry{{ID: 1, Lang: []string{"GO"}}}, map[string][]int{"Go": []int{1}}},
		"ex2": {[]Entry{{ID: 2, Lang: []string{"Lua"}},{ID:1, Lang: []string{"Lua"}}}, map[string][]int{"Lua": []int{1, 2}}},
		"ex3": {[]Entry{{ID: 5, Lang: []string{"UNK"}},{ID:3, Lang: []string{"UNK"}}}, map[string][]int{"UNK": []int{3, 5}}},
		"ex4": {[]Entry{{ID: 111, Lang: []string{"UNK", "NA"}},{ID:21, Lang: []string{"UNK"}}}, map[string][]int{"UNK": []int{21, 111},"NA": []int{111}}},

	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MakeLangs(tc.entries)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.entries)
			}
		})
	}
}

func TestMakeLangList (t *testing.T) {
	tests := map[string]struct {
		langMap map[string][]int
		want []Langs
	}{
		"ex1": {map[string][]int{"Go": []int{1}}, []Langs{{"Go", 1}}},
		"ex2": {map[string][]int{"Lua": []int{1, 2}}, []Langs{{"Lua", 2}}},

	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := MakeLangList(tc.langMap)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.langMap)
			}
		})
	}
}

func TestLSplit (t *testing.T) {
	tests := map[string]struct {
		s string
		want []string
	}{
		"ex1": {"One/Two", []string{"One", "Two"}},
		"ex2": {"Three\\Four", []string{"Three", "Four"}},
		"ex3": {"Five", []string{"Five"}},

	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := LSplit(tc.s)
			if reflect.DeepEqual(got,tc.want) != true {
				t.Errorf("\nGot: %+v \nWant: %+v \nGiven: %+v", got, tc.want, tc.s)
			}
		})
	}
}
