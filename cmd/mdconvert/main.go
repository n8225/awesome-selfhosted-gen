package main

import (
	"bufio"

	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/n8225/ash_gen/pkg/exporter"
	"github.com/n8225/ash_gen/pkg/getexternal"
	"github.com/n8225/ash_gen/pkg/parse"
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
	l := new(parse.List)
	l.Entries = e
	l.TagList = parse.MakeTags(l.Entries)
	l.CatList = parse.MakeCats(l.Entries)
	l.LangList = parse.MakeLangs(l.Entries)
	exporter.ToJSON(*l, "list")
	exporter.ToYaml(*l, "list")
}

func freeReadMd(path, gh string) []parse.Entry {
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
	entries := []parse.Entry{}
	glregex := regexp.MustCompile("^(http.://)(www.){0,1}(gitlab.com)/(.*)/(.*)$")
	ghregex := regexp.MustCompile("^(http.://)(www.){0,1}(github.com)/(.*)$")
	bbregex := regexp.MustCompile("^(http.://)(www.){0,1}(bitbucket.org)/(.*)/(.*)$")
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
				e := new(parse.Entry)
				e.Cat = tag2
				//e.Tags = strings.Trim(strings.Join([]string{tag2, tag3, tag4, tagi}, ", "), " , ")
				if tag2 != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tag2))
					e.Tags = append(e.Tags, parse.Tagmap[tag2]...)
				}
				if tag3 != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tag3))
					e.Tags = append(e.Tags, parse.Tagmap[tag3]...)
				}
				if tag4 != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tag4))
					e.Tags = append(e.Tags, parse.Tagmap[tag4]...)
				}
				if tagi != "" {
					//e.Tags = append(e.Tags, strings.TrimSpace(tagi))
					e.Tags = append(e.Tags, parse.Tagmap[tagi]...)
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
					e.License = parse.LSplit(strings.TrimSpace(result[0][5]))
					e.Lang = parse.LangSplit(strings.TrimSpace(result[0][6]))
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
					glAPI := "https://gitlab.com/api/v4/projects/" + result[0][4] + "%2F" + result[0][5]
					e.Stars, e.Updated = getexternal.GetGLRepo(glAPI)

				} else if ghregex.MatchString(e.Source) {
					src := ""
					result := ghregex.FindAllStringSubmatch(e.Source, -1)
					ghur := strings.TrimSpace(result[0][4])
					e.Stars, e.Updated, _, _, src = getexternal.GetGHRepo(ghur, gh, src)
					if src != "" {
						fmt.Println("Updated " + e.Source + " to https://www.github.com/" + src)
						e.Source = "https://www.github.com/" + src
					}
				} else if bbregex.MatchString(e.Source) {
					//result := bbregex.FindAllStringSubmatch(e.Source, -1)
					//bbAPI := "https://bitbucket.org/api/2.0/repositories/" + result[0][4] + "%2F" + result[0][5]
					e.Stars, e.Updated = getexternal.GetBbRepo(e.Source)
				}
				entries = append(entries, *e)
			}
		}
	}
	return entries
}
