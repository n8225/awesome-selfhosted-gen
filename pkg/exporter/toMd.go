package exporter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

type MdList struct {
	FrontMatter     []string
	TableOfContents MdTOC
}

type MdTOC struct {
	Cat []Cat
}
type Cat struct {
	Name     string
	ChildCat []Cat
}

func createMd() (cats map[string]parse.Category) {
	catFile, err := os.Open("./output/categories.yaml")
	if err != nil {
		log.Error().Stack().Err(err)
	}
	fbytes, _ := ioutil.ReadAll(catFile)
	err = yaml.Unmarshal(fbytes, &cats)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	return
}

func importEntries() []parse.Entry {
	files, err := ioutil.ReadDir("./list")
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	l := make([]parse.Entry, len(files)+1)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			e := parse.Entry{}
			fFile, err := os.Open("./list/" + f.Name())
			if err != nil {
				log.Error().Stack().Err(err)
			}
			fbytes, _ := ioutil.ReadAll(fFile)
			err = yaml.Unmarshal(fbytes, &e)
			if err != nil {
				log.Error().Stack().Err(err)
			}
			l[e.ID] = e
		}
	}
	return l
}

func ToMD() {
	t := template.Must(template.New("mdList.tmpl").Funcs(template.FuncMap{
		"mdLink": func(s string) string {
			return fmt.Sprintf("[%s](#%s)", s, makeTOClink(s))
		},
		"mdEntry": func(e parse.Entry) string {
			link, source := mainLink(e.Site, e.Source)
			demo := linkSyntax(e.Demo, "Demo")
			clients := clientLinks(e.Clients)
			linkString := links([3]string{demo, source, clients})
			return fmt.Sprintf("- [%s](%s) %s- %s%s `%s` `%s`", e.Name, link, pDep(e.Pdep), e.Descrip, linkString, e.License, e.Lang)
		},
	}).ParseFiles("./templates/mdList.tmpl"))

	f, err := os.Create("./output/README.md")
	if err != nil {
		log.Error().Err(err)
	}

	err = t.Execute(f, struct {
		Cats    map[string]parse.Category
		Entries []parse.Entry
	}{
		Cats:    createMd(),
		Entries: importEntries(),
	})
	if err != nil {
		log.Error().Err(err)
	}
	f.Close()
}

func makeTOClink(s string) string {
	r := map[string]string{
		" ": "-",
		"(": "",
		")": "",
	}
	for k, v := range r {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}
