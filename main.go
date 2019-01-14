package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// List is the total struct
type List struct {
	Entries  []Entry	`json:"Entries"`
	LangList []Langs	`json:"Langs" yaml:"-"`
	//CatList  []Cats		`json:"Cats", yaml:""`
	TagList  []Tags		`json:"Tags" yaml:"-"`
}
// Entry is the structure of each entry
type Entry struct {
	ID      int      `json:"ID" yaml:"ID"`
	Name    string   `json:"N" yaml:"Name"`
	Descrip string   `json:"D" yaml:"Description,flow"`
	Source  string   `json:"Sr" yaml:"Source Code"`
	Demo    string   `json:"Dem,omitempty" yaml:"Demo,omitempty"`
	Clients []string `json:"CL,omitempty" yaml:"Clients,omitempty"`
	Site    string   `json:"Si,omitempty" yaml:"Website,omitempty"`
	License []string `json:"Li" yaml:"License"`
	Lang    []string `json:"La" yaml:"Languages"`
	//Cat     string   `json:"C,omitempty"`
	Tags    []string `json:"T" yaml:"Tags"`
	NonFree    bool  `json:"NF,omitempty" yaml:"NonFree,omitempty"`
	Pdep    bool     `json:"P,omitempty" yaml:"ProprietaryDependency,omitempty"`
	Stars	int		`json:"stars,omitempty" yaml:"-"`
	Created	string	`json:"create,omitempty" yaml:"-"`
	Updated	string	`json:"update,omitempty" yaml:"-"`
}

// Licenses is the struct of licenses
type Langs struct {
	Lang	string `json:"Lang"`
	Count 	int `json:"Count"`
}
//Category struct
type Cats struct {
	Cat   string
	Count int
}
//Tags Struct
type Tags struct {
	Tag   string	`json:"Tag"`
	Count int		`json:"C"`
}

func main() {
	var path string
	const (
		defaultPath = ""
		usage = "Path to Readme.md(On windows wrap path in \""
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

func makeLangs(entries []Entry) []Langs {
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
			langMap[item] +=1
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

func makeTags(entries []Entry) []Tags {
	tagsl := []Tags{}
	t := new(Tags)
	var tmp []string
	for _, e := range entries {
		tmp = append(tmp, e.Tags...)
	}
	tagMap := make(map[string]int)
	for _, item := range tmp {
		_, exist := tagMap[item]
		if exist {
			tagMap[item] +=1
		} else {
			tagMap[item] = 1
		}
	}
	for k, v := range tagMap {
		t.Tag = k
		t.Count = v
		tagsl = append(tagsl, *t)
	}
	sort.Slice(tagsl, func(i, j int) bool {
		return tagsl[i].Tag < tagsl[j].Tag
	})
	return tagsl
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
				if strings.Contains(pdep,"⚠") == true  {
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
					ghur:= strings.TrimSpace(result[0][4])
					e.Stars, e.Updated = getGHRepo(ghur, gh)

				}

				entries = append(entries, *e)
			}
		}
	}
	return entries
}

func getGLRepo (url string) (int, string){
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	type gl struct {
		Stars int `json:"star_count"`
		Created string `json:"created_at"`
		Updated string `json:"last_activity_at"`
		//Node_id int `json:"id"`
	}
	thisgl := gl{}
	err = json.Unmarshal(body, &thisgl)
	if err != nil {
		log.Fatal(err)
	}
	return thisgl.Stars, strings.Split(thisgl.Updated, "T")[0]
}

type Gh struct {
	Data struct {
		Repository struct {
			Stargazers struct {
				TotalCount int `json:"totalCount"`
			} `json:"stargazers"`
			DefaultBranchRef struct {
				Target struct {
					History struct {
						Edges [] struct {
							Node struct {
								CommittedDate string `json:"committedDate"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"history"`
				} `json:"target"`
			} `json:"defaultBranchRef"`
		} `json:"repository"`
	} `json:"data"`
 Errors [] struct {
	 Message string `json:"message"`
	 Type string `json:"type"`
	 path string `json:"path"`
 }`json:"errors"`
}

func getGHRepo (ur, ght string) (int, string){
	if strings.Contains(ur, "/") != true {
		log.Printf("No repository provided. Update %s to include a repo.", ur)
		return 0, ""
	}
	r := strings.Split(ur, "/")
	if r[1] == "" {
		log.Printf("No repository provided. Update %s to include a repo.", ur)
		return 0, ""
	}
	var jsonStr = []byte(`{"query":"{repository(owner:\"` + r[0] + `\",name:\"` + r[1] + `\"){stargazers{totalCount}defaultBranchRef{target{... on Commit{history(first: 1){edges{node{committedDate}}}}}}}}"}`)

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "bearer " + ght)
	req.Header.Set("Accept", "application/vnd.github.quicksilver-preview+json")
	req.Header.Set("Content-Type", "application/json")
	//log.Printf("req: ", req.Body)

	client := &http.Client{}
	res, err :=client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	//log.Printf("Body: %s\n", body)

	gh := &Gh{}
		err = json.Unmarshal(body, &gh)
		if err != nil {
		log.Fatal(err)
	}
	//Fallback to github v3 api on error
	if gh.Data.Repository.Stargazers.TotalCount == 0 {
		log.Printf("Error: https://github.com/%s/%s => %v Trying github api v3.  ", r[0], r[1], gh.Errors[0].Message)
		return ghv3api(r[0], r[1], ght)
	}
	//return *repo.StargazersCount
	return gh.Data.Repository.Stargazers.TotalCount, strings.Split(gh.Data.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate, "T")[0]
}

func ghv3api(u, r, ght string) (int, string) {
	ghURL := "https://api.github.com/repos/" + u + "/" + r + "?" //+ ght

	//fmt.Println(ghURL)
	// request http api
	req, err := http.NewRequest("GET", ghURL, nil)
	req.Header.Set("Authorization", "bearer " + ght)
	client := http.Client{}
	res, err :=client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Printf("StatusCode: %v. Source link is invalid: %v/%v\n", res.StatusCode, u, r)
		return 0, ""
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Body: %s\n", body)

	type Gh struct {
		FullName string `json:"full_name"`
		Stars int `json:"stargazers_count"`
		Created string `json:"created_at"`
		Updated string `json:"pushed_at"`
	}
	gh := Gh{}
	err = json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Update https://github.com/%s/%s Source code link to: https://github.com/%s\n", u, r, gh.FullName)
	if "/" + u + "/" + r != gh.FullName {
		log.Printf("Retrying Github api v4 with %s", gh.FullName)
		getGHRepo(gh.FullName, ght)
	}
	return gh.Stars, strings.Split(gh.Updated, "T")[0]
}


func toJson(list List) {
		yamlFile, err := os.Create("./output.yaml")
		if err != nil {
			fmt.Println(err)
		}
		defer yamlFile.Close()
		YAML, err := yaml.Marshal(list)
		if err != nil{
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

var langs = map[string][]string{
	".NET":				{".NET"},
	"Angular":			{"HTML5"},
	"C":				{"C"},
	"C#":				{"C#"},
	"C++":				{"C++"},
	"CSS":				{"HTML5"},
	"ClearOS":			{"PHP"},
	"Clojure":			{"Clojure"},
	"ClojureScript":	{"Clojure"},
	"CommonLisp":		{"CommonLisp"},
	"Django":			{"Python"},
	"Docker":			{"Docker"},
	"Elixir":			{"Elixir"},
	"Erlang":			{"Erlang"},
	"Go":				{"Go"},
	"GO":				{"Go"},
	"Golang":			{"Golang"},
	"HTML":				{"HTML5"},
	"HTML5":			{"HTML5"},
	"Haskell":			{"Haskell"},
	"Java":				{"Java"},
	"JavaScript":		{"HTML5"},
	"Javascript":		{"HTML5"},
	"Kotlin":			{"Kotlin"},
	"Linux":			{"Shell"},
	"Lua":				{"Lua"},
	"Nix":				{"Nix"},
	"Node.js":			{"Nodejs"},
	"NodeJS":			{"Nodejs"},
	"Nodejs":			{"Nodejs"},
	"OCAML":			{"OCaml"},
	"OCaml":			{"OCaml"},
	"Objective-C":		{"Objective-C"},
	"PHP":				{"PHP"},
	"PL":				{"Perl"},
	"Perl":				{"Perl"},
	"Python":			{"Python"},
	"Ruby":				{"Ruby"},
	"Rust":				{"Rust"},
	"Scala":			{"scala"},
	"Shell":			{"Shell"},
	"TypeScript":		{"TypeScript"},
	"VueJS":			{"HTML5"},
	"YAML":				{"YAML"},
	"pgSQL":			{"pgSQL"},
	"python":			{"Python"},
	"С++":				{"C++"},
	"rc":				{"rc"},
}

var tags = map[string][]string{
	"Analytics":     {"Analytics"},
	"Web Analytics": {"Web Analytics"},
	"Archiving and Digital Preservation (DP)": {"archiving", "Digital Preservation"},
	"Automation":                          {"Automation"},
	"Blogging Platforms":                  {"Blog"},
	"Bookmarks and Link Sharing":          {"Bookmarks", "Links"},
	"Calendaring and Contacts Management": {"Calendar", "Contacts"},
	"CalDAV or CardDAV servers":           {"CalDAV", "CardDav"},
	"Communication systems":               {"Communications"},
	"Custom communication systems":        {},
	"Email":                               {"Email"},
	"Complete solutions":                  {"Complete Email"},
	"Mail Transfer Agents":                {"MTA"},
	"MTAs / SMTP servers":                 {"SMTP"},
	"Mail Delivery Agents":                {"MDA"},
	"MDAs - IMAP/POP3 software":           {"IMAP", "POP3"},
	"Mailing lists and Newsletters":       {"Mailng List", "Newsletters"},
	"Mailing lists servers and mass mailing software - one message to many recipients.": {"Mass mail"},
	"Webmail clients": {"Webmail"},
	"IRC":             {"IRC"},
	"[IRC](https://en.wikipedia.org/wiki/Internet_Relay_Chat) communication software": {},
	"SIP": {"SIP"},
	"[SIP](https://en.wikipedia.org/wiki/Session_Initiation_Protocol) telephony software": {},
	"IPBX": {"IPBX"},
	"[IPBX](https://en.wikipedia.org/wiki/IP_PBX) telephony software": {},
	"Social Networks and Forums":                                      {"Social Network", "Forum"},
	"XMPP":                                                            {"XMPP"},
	"XMPP Servers":                                                    {"XMPP Server"},
	"XMPP Web Clients":                                                {"XMPP Webclient"},
	"Conference Management":                                           {"Conference Mgmnt"},
	"Content Management Systems (CMS)":                                {"CMS"},
	"Recipe management":                                               {"Recipe Mgmnt"},
	"E-commerce":                                                      {"E-commerce"},
	"DNS":                                                             {"DNS"},
	"Document Management":                                             {"Doc Mgmnt"},
	"E-books and Integrated Library Systems (ILS)":  {"E-book", "ILS"},
	"Enterprise-class library management software.": {},
	"Feed Readers":                             {"RSS", "Feed Reader"},
	"File Sharing and Synchronization":         {"File Sharing"},
	"File transfer/synchronization":            {"File Transfer", "File Sync"},
	"Peer-to-peer filesharing":                 {"P2P"},
	"Object storage/file servers":              {"Object Storage"},
	"Single-click/drag-n-drop upload":          {"Single-Click Upload", "Drag-n-Drop Upload"},
	"Command-line file upload":                 {"CMD line upload"},
	"Web based file managers":                  {"Web File MGR"},
	"Games":                                    {"Game"},
	"Gateways":                                 {"Gateway"},
	"Groupware":                                {"Groupware"},
	"Human Resources Management (HRM)":         {"HRM", "Human Resources MGMNT"},
	"Internet Of Things (IoT)":                 {"Internet Of Things", "IoT"},
	"Learning and Courses":                     {"Learning", "Courses", "LMS"},
	"Maps and Global Positioning System (GPS)": {"Maps", "GPS"},
	"Media Streaming":                          {"Media", "Streaming"},
	"Multimedia Streaming":                     {"Multimedia"},
	"Audio Streaming":                          {"Audio"},
	"Video Streaming":                          {"Video"},
	"Misc/Other":                               {"Misc", "Other"},
	"Money, Budgeting and Management":          {"Money", "Budget"},
	"Note-taking and Editors":                  {"Note-taking", "Editor"},
	"Office Suites":                            {"Office Suites"},
	"Password Managers":                        {"Password Manager"},
	"Pastebins":                                {"Pastebin"},
	"Personal Dashboards":                      {"Dashboard"},
	"Photo and Video Galleries":                {"Photo and Video Gallery"},
	"Polls and Events":                         {"Polls", "Events"},
	"Proxy":                                    {"Proxy"},
	"Read it Later Lists":                      {"Read it Later Lists"},
	"Resource Planning":                        {"Resource Planning"},
	"Search Engines":                           {"Search Engine"},
	"Enterprise Resource Planning":             {"Enterprise rsrc planning"},
	"Software Development":                     {"Software Dev"},
	"Project Management":                       {"Project Mgmnt"},
	"See also [Ticketing](#ticketing), [Task management/To-do lists](#task-managementto-do-lists)": {},
	"IDE/Tools":                   {"IDE"},
	"Continuous Integration":      {"CI"},
	"FaaS/Serverless":             {"FAAS", "Severless"},
	"API Management":              {"API Mgmnt"},
	"Documentation Generators":    {"Documentation Gen"},
	"Localization":                {"Localiztion"},
	"Task management/To-do lists": {"To-Do", "Task Mgmnt"},
	"Ticketing":                   {"Ticketing"},
	"URL Shorteners":              {"URL Shortener"},
	"Wikis":                       {"Wiki"},
	"Self-hosting Solutions":      {"Self-hosting Solution"},
}
