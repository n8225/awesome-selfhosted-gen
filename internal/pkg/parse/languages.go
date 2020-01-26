package parse

import (
	"sort"
	"strings"
)

// Langs is the struct of programming languages
type Langs struct {
	Lang  string `json:"Lang"`
	Count int    `json:"Count"`
}

// MakeLangs creates the language map with ids
func MakeLangs(entries []Entry) map[string][]int {
	langMap := make(map[string][]int)
	for _, e := range entries {
		for _, l := range e.Lang {
			_, exist := langs[l]
			if !exist {
				langs[l] = l
			}
			langMap[langs[l]] = append(langMap[langs[l]], e.ID)
		}
	}
	return SortIDs(langMap)
}
//MakeLangList counts and creates a slice
func MakeLangList(langMap map[string][]int) []Langs {
	langList := []Langs{}
	for k, v := range langMap {
		l := Langs{k, len(v)}
		langList = append(langList, l)
	}
	sort.Slice(langList, func(i, j int) bool {
		return langList[i].Lang < langList[j].Lang
	})
	return langList
}

// LSplit splits language or license string delimited by '/' or '\'
func LSplit(lang string) []string {
	if strings.Contains(lang, "/") {
		return strings.Split(lang, "/")
	} else if strings.Contains(lang, "\\") {
		return strings.Split(lang, "\\")
	} else {
		l := make([]string, 1)
		l[0] = lang
		return l
	}
}

/* // LangSplit creates new language slice from string
func LangSplit(lang string) []string {
	nLangs := LSplit(lang)
	var mLangs []string
	for _, lang := range nLangs {
		mLangs = append(mLangs, langs[lang])
	}
	//fmt.Printf("%v\n", mLangs)
	return mLangs
} */

var langs = map[string]string{
	".NET":          ".NET",
	"Angular":       "HTML5",
	"C":             "C",
	"C#":            "C#",
	"C++":           "C++",
	"CSS":           "HTML5",
	"ClearOS":       "PHP",
	"Clojure":       "Clojure",
	"ClojureScript": "Clojure",
	"CommonLisp":    "CommonLisp",
	"Django":        "Python",
	"Docker":        "Docker",
	"Elixir":        "Elixir",
	"Erlang":        "Erlang",
	"Go":            "Go",
	"GO":            "Go",
	"Golang":        "Golang",
	"HTML":          "HTML5",
	"HTML5":         "HTML5",
	"Haskell":       "Haskell",
	"Java":          "Java",
	"JavaScript":    "HTML5",
	"Javascript":    "HTML5",
	"Kotlin":        "Kotlin",
	"Linux":         "Shell",
	"Lua":           "Lua",
	"Nix":           "Nix",
	"Node.js":       "Nodejs",
	"NodeJS":        "Nodejs",
	"Nodejs":        "Nodejs",
	"OCAML":         "OCaml",
	"OCaml":         "OCaml",
	"Objective-C":   "Objective-C",
	"PHP":           "PHP",
	"PL":            "Perl",
	"Perl":          "Perl",
	"Python":        "Python",
	"Ruby":          "Ruby",
	"Rust":          "Rust",
	"Scala":         "scala",
	"Shell":         "Shell",
	"TypeScript":    "TypeScript",
	"VueJS":         "HTML5",
	"YAML":          "YAML",
	"pgSQL":         "pgSQL",
	"python":        "Python",
	"С++":           "C++",
	"rc":            "rc",
}
