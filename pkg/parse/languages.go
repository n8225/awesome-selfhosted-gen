package parse

import (
	"sort"
	"strings"
	"fmt"
)

// Langs is the struct of programming languages
type Langs struct {
	Lang  string `json:"Lang"`
	Count int    `json:"Count"`
}

func MakeLangs(entries []Entry) []Langs {
	langsl := []Langs{}
	l := new(Langs)
	var tmp []string
	for _, e := range entries {
		tmp = append(tmp, e.Lang...)
	}
	langMap := make(map[string]int)
	for _, item := range tmp {
		_, exist := langMap[item]
		if exist {
			langMap[item]++
		} else {
			langMap[item] = 1
		}
	}
	for k, v := range langMap {
		l.Lang = k
		l.Count = v
		langsl = append(langsl, *l)
	}
	sort.Slice(langsl, func(i, j int) bool {
		return langsl[i].Lang < langsl[j].Lang
	})
	return langsl
}

func LSplit(lang string) []string {
	if strings.Contains(lang, "/") {
		return strings.Split(lang, "/")
	} else if strings.Contains(lang, "\\") {
		fmt.Println(strings.Split(lang, "\\"))
		return strings.Split(lang, "\\")
	} else {
		l := make([]string, 1)
		l[0] = lang
		return l
	}
}

func LangSplit(lang string) []string {
	nLangs := LSplit(lang)
	var mLangs []string
	for _, lang := range nLangs {
		mLangs = append(mLangs, langs[lang]...)
	}
	return mLangs
}

var langs = map[string][]string{
	".NET":          {".NET"},
	"Angular":       {"HTML5"},
	"C":             {"C"},
	"C#":            {"C#"},
	"C++":           {"C++"},
	"CSS":           {"HTML5"},
	"ClearOS":       {"PHP"},
	"Clojure":       {"Clojure"},
	"ClojureScript": {"Clojure"},
	"CommonLisp":    {"CommonLisp"},
	"Django":        {"Python"},
	"Docker":        {"Docker"},
	"Elixir":        {"Elixir"},
	"Erlang":        {"Erlang"},
	"Go":            {"Go"},
	"GO":            {"Go"},
	"Golang":        {"Golang"},
	"HTML":          {"HTML5"},
	"HTML5":         {"HTML5"},
	"Haskell":       {"Haskell"},
	"Java":          {"Java"},
	"JavaScript":    {"HTML5"},
	"Javascript":    {"HTML5"},
	"Kotlin":        {"Kotlin"},
	"Linux":         {"Shell"},
	"Lua":           {"Lua"},
	"Nix":           {"Nix"},
	"Node.js":       {"Nodejs"},
	"NodeJS":        {"Nodejs"},
	"Nodejs":        {"Nodejs"},
	"OCAML":         {"OCaml"},
	"OCaml":         {"OCaml"},
	"Objective-C":   {"Objective-C"},
	"PHP":           {"PHP"},
	"PL":            {"Perl"},
	"Perl":          {"Perl"},
	"Python":        {"Python"},
	"Ruby":          {"Ruby"},
	"Rust":          {"Rust"},
	"Scala":         {"scala"},
	"Shell":         {"Shell"},
	"TypeScript":    {"TypeScript"},
	"VueJS":         {"HTML5"},
	"YAML":          {"YAML"},
	"pgSQL":         {"pgSQL"},
	"python":        {"Python"},
	"С++":           {"C++"},
	"rc":            {"rc"},
}