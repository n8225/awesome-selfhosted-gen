package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Entry is the structure of each entry
type Entry struct {
	ID   int `yaml:"ID" json:"ID"`
	Data struct {
		Name    string   `yaml:"Name" json:"Name"`
		Site    string   `yaml:"Site" json:"Site"`
		Descrip string   `yaml:"Descrip,flow" json:"Descrip"`
		Demo    string   `yaml:"Demo" json:"Demo"`
		Source  string   `yaml:"Source" json:"Source"`
		License string   `yaml:"License" json:"License"`
		Lang    string   `yaml:"Lang" json:"Lang"`
		Cat     string   `yaml:"Category" json:"Cat"`
		Tags    []string `yaml:"Tags" json:"Tags"`
		Free    bool     `yaml:"Free" json:"Free"`
		Pdep    bool     `yaml:"PropDep" json:"Pdep"`
	} `yaml:"Entry" json:"Entry"`
}

func main() {
	c := []*Entry{}
	pathPtr := flag.String("path", "", "Path to Readme.md")
	//fileNfPtr := flag.String("nffile", "non-free.md", "Path to Readme.md")
	//optsPtr := flag.String("opts", "", "what would you like to do?")
	flag.Parse()
	//switch to choose functions
	i := 1
	switch i {
	case 1:
		fmt.Println("1: run all :", *pathPtr)
		c = append(c, freeReadMd(*pathPtr+"\\README.md")...)
		toYaml(c)
		//nreadMD(*pathPtr + "\\non-free.md")
	case 2:
		fmt.Println("2: Run to yaml")
	case 3:
		fmt.Println("3: Run to md")
	}
}

func freeReadMd(path string) []*Entry {
	fmt.Println(path)
	inputFile, _ := os.Open(path)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	list := false
	var tag2, tag3, tag4, tagi string
	var i int
	pattern := *regexp.MustCompile("^\\s{0,4}\\Q- [\\E(?P<name>.*?)\\Q](\\E(?P<site>.*?)\\)(?P<pdep>\\Q `⚠` - \\E|\\Q - \\E)(?P<desc>.*?[.])(?:\\s\x60|\\s\\(.*\x60)(?P<license>.*?)\\Q` `\\E(?P<lang>.*?)\\Q`\\E")
	demoP := *regexp.MustCompile("\\Q[Demo](\\E(.*?)\\Q)\\E")
	sourceP := *regexp.MustCompile("\\Q[Source Code](\\E(.*?)\\Q)\\E")
	entries := []*Entry{}

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
				e.Data.Cat = tag2
				//e.Tags = strings.Trim(strings.Join([]string{tag2, tag3, tag4, tagi}, ", "), " , ")
				if tag2 != "" {
					e.Data.Tags = append(e.Data.Tags, tag2)
				}
				if tag3 != "" {
					e.Data.Tags = append(e.Data.Tags, tag3)
				}
				if tag4 != "" {
					e.Data.Tags = append(e.Data.Tags, tag4)
				}
				if tagi != "" {
					e.Data.Tags = append(e.Data.Tags, tagi)
				}
				e.Data.Free = true
				entries = append(entries, e)
				if pattern.MatchString(scanner.Text()) {
					e.ID = i
					i++
					result := pattern.FindAllStringSubmatch(scanner.Text(), -1)
					e.Data.Name = result[0][1]
					e.Data.Site = result[0][2]
					e.Data.Descrip = result[0][4]
					e.Data.License = result[0][5]
					e.Data.Lang = result[0][6]
					e.Data.Pdep = result[0][3] == " `⚠` - "
				}
				if demoP.MatchString(scanner.Text()) {
					result := demoP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Data.Demo = result[0][1]
				} else {
					e.Data.Demo = ""
				}
				if sourceP.MatchString(scanner.Text()) {
					result := sourceP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Data.Source = result[0][1]
				} else {
					e.Data.Source = ""
				}
			}
		}
	}
	return entries
}

func toYaml(entries []*Entry) {
	yamlFile, err := os.Create("./output.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer yamlFile.Close()

	yamlWriter := io.Writer(yamlFile)
	yencoder := yaml.NewEncoder(yamlWriter)

	jsonFile, err := os.Create("./output.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonWriter := io.Writer(jsonFile)
	jencoder := json.NewEncoder(jsonWriter)

	for i := range entries {
		err = yencoder.Encode(&entries[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = jencoder.Encode(&entries[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	jsonFile.Close()
	yamlFile.Close()

}
