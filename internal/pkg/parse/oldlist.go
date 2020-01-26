package parse

import (
	"bufio"
	"regexp"

	"fmt"
	"os"
	"strings"
)

// List is the total struct
type List struct {
	Entries  []Entry `json:"Entries" yaml:"Entries"`
	LangList []Langs `json:"Langs" yaml:"-"`
	LangIDs	 map[string][]int `json:"LangIds" yaml:"-"`
	CatList  []Cat  `json:"Cats" yaml:"-"`
	CatIDs 	 map[string][]int `json:"CatIds" yaml:"-"`
	TagList  []Tags  `json:"Tags" yaml:"-"`
	TagIDs   map[string][]int `json:"TagIds" yaml:"-"`

}

// GetHighestID function to get last ID from entry struct
// func GetHighestID(entries []Entry) int {
// 	max := entries[0]
// 	for _, entries := range entries {
// 		if entries.ID > max.ID {
// 			max = entries
// 		}
// 	}
// 	return max.ID
// }

//MdParser reads README.md and parses data from it.
func oldMdParser(path, gh string) []Entry {
	fmt.Println("Parsing:", path)
	inputFile, _ := os.Open(path)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	list := false
	var tag2, tag3, tag4, tagi string
	var l, i = 0, 0
	entries := []Entry{}

	for scanner.Scan() {
		l++
		if strings.HasPrefix(scanner.Text(), "<!-- BEGIN SOFTWARE LIST -->") {
			list = true
		} else if strings.HasPrefix(scanner.Text(), "<!-- END SOFTWARE LIST -->") {
			list = false
		}
		if list {
			if strings.HasPrefix(scanner.Text(), "## ") {
				tag2, tag3, tag4, tagi = strings.Trim(scanner.Text(), "## "), "", "", ""
			}
			if strings.HasPrefix(scanner.Text(), "### ") {
				tag4, tagi, tag3 = "", "", strings.Trim(scanner.Text(), "### ")
			}
			if strings.HasPrefix(scanner.Text(), "#### ") {
				tagi, tag4 = "", strings.Trim(scanner.Text(), "#### ")
			}
			if strings.HasPrefix(scanner.Text(), "_") {
				tagi = strings.Trim(scanner.Text(), "_")
			}
			if strings.HasPrefix(scanner.Text(), "- [") || strings.HasPrefix(scanner.Text(), "  - [") {
				e := new(Entry)
				e.Cat = tag2

				if tag2 != "" {
					e.Tags = append(e.Tags, Tagmap[tag2]...)
				}
				if tag3 != "" {
					e.Tags = append(e.Tags, Tagmap[tag3]...)
				}
				if tag4 != "" {
					e.Tags = append(e.Tags, Tagmap[tag4]...)
				}
				if tagi != "" {
					e.Tags = append(e.Tags, Tagmap[tagi]...)
				}
				if regexp.MustCompile(Pattern).MatchString(scanner.Text()) {
					i++
					e.Line = l
					e.ID = i
					e.MD = scanner.Text()
					e.Name = GetName(e.MD)
					e.Descrip = GetDescrip(e.MD)
					e.License = GetLicense(e.MD)
					e.Lang = GetLang(e.MD)
					e.Pdep = GetPdep(e.MD)
					e.Demo = GetDemo(e.MD)
					e.Clients = GetClients(e.MD)
					e.Site = GetSite(e.MD)
					e.Source, e.SourceType = GetSource(e.MD)

				} else {
					fmt.Printf("Failed to match pattern, Line: %d : %s", l, scanner.Text())
				}
				entries = append(entries, *e)
			}
		}
	}
	fmt.Printf("Found %d entries\n", len(entries))
	return entries
}
