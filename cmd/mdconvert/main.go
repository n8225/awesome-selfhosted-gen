package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ash_gen/pkg/util"
)

func main() {
	var path string
	const (
		defaultPath = ""
		usage       = "Path to Readme.md(On windows wrap path in \""
	)
	flag.StringVar(&path, "path", defaultPath, usage)
	flag.StringVar(&path, "p", defaultPath, usage)
	var ghToken = flag.String("ghtoken", "", "github oauth token")
	//var glToken = flag.String("gltoken", "", "gitlab oauth token")
	flag.Parse()
	apath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	e := freeReadMd(apath, *ghToken)
	l := new(List)
	l.Entries = e
	l.TagList = makeTags(e)
	//l.CatList = makeCats(e)
	l.LangList = makeLangs(e)
	toJson(*l)
}





/*func makeCats(entries []Entry) []Cats {
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
			catMap[item] +=1
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
}*/

func freeReadMd(path, gh string) []Entry {
	fmt.Println("Parsing:", path)
	inputFile, _ := os.Open(path)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	list := false
	var tag2, tag3, tag4, tagi string
	var i int
	var site, pdep string

	//pattern := *regexp.MustCompile(`(?m)^\s{0,4}- \[(?P<name>.*?)\Q](\E(?P<site>.*?)\)(?P<pdep>\s-\s\s\x60⚠\x60|\s\x60⚠\x60\s-\s|\x60⚠\x60|\s-\s)(?P<desc>.*?[.])(?:\s\x60|\s\(.*\x60)(?P<license>.*?)\x60\s\x60(?P<lang>.*?)\x60$`)
	pattern := *regexp.MustCompile("^\\s{0,4}\\Q- [\\E(?P<name>.*?)\\Q](\\E(?P<site>.*?)\\)(?P<pdep>\\Q `⚠` - \\E|\\Q -  `⚠`\\E|\\Q - \\E)(?P<desc>.*?[.])(?:\\s\x60|\\s\\(.*\x60)(?P<license>.*?)\\Q` `\\E(?P<lang>.*?)\\Q`\\E")
	demoP := *regexp.MustCompile("\\Q[Demo](\\E(.*?)\\Q)\\E")
	sourceP := *regexp.MustCompile("\\Q[Source Code](\\E(.*?)\\Q)\\E")
	clientP := *regexp.MustCompile("\\Q[Clients](\\E(.*?)\\Q)\\E")
	entries := []Entry{}
	glregex := regexp.MustCompile("^(http.://)(www.){0,1}(gitlab.com)/(.*)/(.*)$")
	ghregex := regexp.MustCompile("^(http.://)(www.){0,1}(github.com)/(.*)$")
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "<!-- BEGIN SOFTWARE LIST -->") {
			list = true
		} else if strings.HasPrefix(scanner.Text(), "<!-- END SOFTWARE LIST -->") {
			list = false
		}
		if list == true {
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
				//e.Cat = tag2
				//e.Tags = strings.Trim(strings.Join([]string{tag2, tag3, tag4, tagi}, ", "), " , ")
				if tag2 != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tag2))
					e.Tags = append(e.Tags, tags[tag2]...)
				}
				if tag3 != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tag3))
					e.Tags = append(e.Tags, tags[tag3]...)
				}
				if tag4 != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tag4))
					e.Tags = append(e.Tags, tags[tag4]...)
				}
				if tagi != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tagi))
					e.Tags = append(e.Tags, tags[tagi]...)
				}

				if pattern.MatchString(scanner.Text()) {
					e.ID = i
					i++
					result := pattern.FindAllStringSubmatch(scanner.Text(), -1)
					e.Name = strings.TrimSpace(result[0][1])
					//Set Test entry to nonfree
					if result[0][1] == "TEST" {
						e.NonFree = true
					}
					site = strings.TrimSpace(result[0][2])
					e.Descrip = strings.TrimSpace(result[0][4])
					e.License = lSplit(strings.TrimSpace(result[0][5]))
					e.Lang = langSplit(strings.TrimSpace(result[0][6]))
					pdep = result[0][3]
				}
				if strings.Contains(pdep, "⚠") == true {
					e.Pdep = true
				}
				if demoP.MatchString(scanner.Text()) {
					result := demoP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Demo = strings.TrimSpace(result[0][1])
				}
				if clientP.MatchString(scanner.Text()) {
					result := clientP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Clients = append(e.Clients, strings.TrimSpace(result[0][1]))
				}
				if sourceP.MatchString(scanner.Text()) {
					result := sourceP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Source = strings.TrimSpace(result[0][1])
					e.Site = site
				} else {
					e.Source = site
				}
				if glregex.MatchString(e.Source) {
					result := glregex.FindAllStringSubmatch(e.Source, -1)
					glApi := "https://gitlab.com/api/v4/projects/" + result[0][4] + "%2F" + result[0][5]
					e.Stars, e.Updated = getGLRepo(glApi)

				} else if ghregex.MatchString(e.Source) {
					result := ghregex.FindAllStringSubmatch(e.Source, -1)
					ghur := strings.TrimSpace(result[0][4])
					e.Stars, e.Updated = getGHRepo(ghur, gh)

				}

				entries = append(entries, *e)
			}
		}
	}
	return entries
}





func toJson(list List) {
	yamlFile, err := os.Create("./output.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer yamlFile.Close()
	YAML, err := yaml.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
	}
	yamlFile.Write(YAML)
	yamlFile.Close()

	jsonFile, err := os.Create("./output.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	JSON, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	//fmt.Println(string(JSON))
	jsonFile.Write(JSON)
	jsonFile.Close()

	jsonFileMin, err := os.Create("./output.min.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	JSONmin, err := json.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	//fmt.Println(string(JSON))
	jsonFileMin.Write(JSONmin)
	jsonFileMin.Close()
}

func lSplit(lang string) []string {
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

func langSplit(lang string) []string {
	nLangs := lSplit(lang)
	var mLangs []string
	for _, lang := range nLangs {
		mLangs = append(mLangs, langs[lang]...)
	}
	return mLangs
}




