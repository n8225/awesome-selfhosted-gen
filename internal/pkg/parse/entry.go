package parse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/n8225/awesome-selfhosted-gen/internal/pkg/getexternal"
)

//Entry is the structure of each entry
type Entry struct {
	ID         int      `json:"ID" yaml:"ID"`
	Name       string   `json:"N" yaml:"Name"`
	Descrip    string   `json:"D" yaml:"Description,flow"`
	Source     string   `json:"Sr" yaml:"Source Code"`
	Demo       string   `json:"Dem,omitempty" yaml:"Demo,omitempty"`
	Clients    []string `json:"CL,omitempty" yaml:"Clients,omitempty"`
	Site       string   `json:"Si,omitempty" yaml:"Website,omitempty"`
	License    []string `json:"Li" yaml:"License"`
	Lang       []string `json:"La" yaml:"Languages"`
	Cat        string   `json:"C,omitempty" yaml:"C"`
	Tags       []string `json:"T" yaml:"Tags"`
	Pdep       bool     `json:"P,omitempty" yaml:"ProprietaryDependency,omitempty"`
	MD         string   `json:"-" yaml:"MD"`
	SourceType string   `json:"SourceType" yaml:"SourceType,omitempty"`
	Line       int      `json:"Line" yaml:"Line"`
	Stars      int      `json:"stars,omitempty" yaml:"stars,omitempty"`
	Updated    string   `json:"update,omitempty" yaml:"update,omitempty"`
	NonFree    bool     `json:"NF,omitempty" yaml:"NonFree,omitempty"`
	Gitdata    Gitdata  `json:"-" yaml:"Gitdata,omitempty"`
}

//Gitdata holds data retrieved from provider apis.
type Gitdata struct {
	License  string
	Language string
	Archived bool
	Source   string
	Errors   []string
	Stars    int
	Updated  string
}

//AddEntry creates an Entry
func AddEntry(i, l int, t string, catAarr []tmpCat) *Entry {
	e := new(Entry)
	e.Line = l
	e.ID = i
	e.MD = t
	e.Name = GetName(e.MD)
	e.Descrip = GetDescrip(e.MD)
	e.License = GetLicense(e.MD)
	e.Lang = GetLang(e.MD)
	e.Pdep = GetPdep(e.MD)
	e.Demo = GetDemo(e.MD)
	e.Clients = GetClients(e.MD)
	e.Site = GetSite(e.MD)
	e.Source, e.SourceType = GetSource(e.MD)
	e.Cat, e.Tags = getCat(l, catAarr)
	return e
}

//GetCat get categories and tags from the temporary Arrays
func getCat(l int, catAarr []tmpCat) (cat string, tags []string) {
	for _, c := range catAarr{
		if (c.level == 1 && l > c.start && l < c.stop) {
			cat = c.cat
			tags = append(tags, c.cat)
		} else if (c.start < l && l < c.stop) {
			tags = append(tags, c.cat)
		}
	}
	return
}

//Pattern to parse data from markdown entry.
const Pattern string = "^\\s{0,4}\\Q- [\\E(?P<name>.*?)\\Q](\\E(?P<site>.*?)\\)(?P<pdep>\\Q `⚠` - \\E|\\Q -  `⚠`\\E|\\Q - \\E)(?P<desc>.*?[.])(?:\\s\x60|\\s\\(.*\x60)(?P<license>.*?)\\Q` `\\E(?P<lang>.*?)\\Q`\\E"

//GetName parses the name from the markdown entry.
func GetName(e string) string {
	return strings.TrimSpace(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][1])
}

//GetDescrip parses the name from the markdown entry.
func GetDescrip(e string) string {
	return strings.TrimSpace(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][4])
}

//GetLicense parses the license from the markdown entry and separates multiple licenses into a slice.
func GetLicense(e string) []string {
	return LSplit(strings.TrimSpace(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][5]))
}

//GetLang parses the programming language from the markdown entry and separates multiple licenses into a slice.
func GetLang(e string) []string {
	return LSplit(strings.TrimSpace(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][6]))
}

//GetPdep determines whether an entry has a proprietary dependency.
func GetPdep(e string) bool {
	return strings.Contains(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][3], "⚠") 
}

//GetDemo parses the URL for the Demo site.
func GetDemo(e string) string {
	const demop string = "\\Q[Demo](\\E(.*?)\\Q)\\E"
	if regexp.MustCompile(demop).MatchString(e) {
		return strings.TrimSpace(regexp.MustCompile(demop).FindAllStringSubmatch(e, -1)[0][1])
	}
	return ""
}

//GetClients parses the URL for the Client site.//TODO This needs to parse multiple client sites.
func GetClients(e string) []string {
	const clientp string = "\\Q[Clients](\\E(.*?)\\Q)\\E"
	var clients []string
	if regexp.MustCompile(clientp).MatchString(e) {
		return append(clients, strings.TrimSpace(regexp.MustCompile(clientp).FindAllStringSubmatch(e, -1)[0][1]))
	}
	return nil
}

//GetSource parses the URL for the Source Code and might determine the source code hosting site.//TODO this should use url.Parse to check against host name
func GetSource(e string) (u, s string) {
	const sourcep string = "\\Q[Source Code](\\E(.*?)\\Q)\\E"
	if regexp.MustCompile(sourcep).MatchString(e) {
		u = strings.TrimSpace(regexp.MustCompile(sourcep).FindAllStringSubmatch(e, -1)[0][1])
	} else {
		u = strings.TrimSpace(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][2])
	}
	switch true {
	case strings.Contains(u, "github.com"):
		s = "Github"
	case strings.Contains(u, "gitlab.com"):
		s = "Gitlab"
	case strings.Contains(u, "bitbucket.com"):
		s = "Bitbucket"
	default:
		s = ""
	}
	return
}

//GetSite parses the URL for the Web site.
func GetSite(e string) string {
	const sourcep string = "\\Q[Source Code](\\E(.*?)\\Q)\\E"
	if regexp.MustCompile(sourcep).MatchString(e) {
		return strings.TrimSpace(regexp.MustCompile(Pattern).FindAllStringSubmatch(e, -1)[0][2])
	}
	return ""
}

//GetGitdata pulls data from the providers api
func GetGitdata(e Entry, ght string) *Gitdata {
	const glp string = "^(http.://)(www.){0,1}(gitlab.com)/(.*)$"
	const bbp string = "^(http.://)(www.){0,1}(bitbucket.org)/(.*)/(.*)$"
	const ghp string = "^(http.://)(www.){0,1}(github.com)/(.*)$"

	g := new(Gitdata)
	switch e.SourceType {
	case "Gitlab":
		result := regexp.MustCompile(glp).FindAllStringSubmatch(e.Source, -1)
		g.Stars, g.Updated = getexternal.GetGLRepo(result[0][4])
	case "Bitbucket":
		g.Stars, g.Updated = getexternal.GetBbRepo(e.Source)
	case "Github":
		g.Stars, g.Updated, g.License, g.Language, g.Errors = getexternal.GetGH(e.Source, ght, nil)
		if g.Errors != nil {
			fmt.Printf("Github API errors, Line %d: %v\n", e.Line, g.Errors)
		}
	default:

	}
	return g
}
