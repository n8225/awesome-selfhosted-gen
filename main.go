package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Entry is the structure of each entry
type Entry struct {
	ID      int      `yaml:"ID" json:"ID"`
	Name    string   `yaml:"N" json:"N"`
	Descrip string   `yaml:"D,flow" json:"D"`
	Source  string   `yaml:"Src,omitempty" json:"Sr,omitempty"`
	Demo    string   `yaml:"Demo,omitempty" json:"Dem,omitempty"`
	Site    string   `yaml:"Site,omitempty" json:"Si,omitempty"`
	License string   `yaml:"Lic" json:"Li"`
	Lang    []string `yaml:"Lang" json:"La"`
	Cat     string   `yaml:"Cat" json:"C"`
	Tags    []string `yaml:"T" json:"T"`
	Free    bool     `yaml:"Free,omitempty" json:"F,omitempty"`
	Pdep    bool     `yaml:"Pdep,omitempty" json:"P,omitempty"`
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
		c = append(c, freeReadMd(*pathPtr)...)
		//c = append(c, nreadMD(*pathPtr)...)
		toYaml(c)
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
	var site, pdep string
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
				e.Cat = tag2
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
				//e.Free = true
				entries = append(entries, e)
				if pattern.MatchString(scanner.Text()) {
					e.ID = i
					i++
					result := pattern.FindAllStringSubmatch(scanner.Text(), -1)
					e.Name = strings.TrimSpace(result[0][1])
					site = strings.TrimSpace(result[0][2])
					e.Descrip = strings.TrimSpace(result[0][4])
					e.License = strings.TrimSpace(result[0][5])
					e.Lang = langSplit(strings.TrimSpace(result[0][6]))
					pdep = result[0][3]
				}
				if pdep == " `⚠` - " {
					e.Pdep = true
				}
				if demoP.MatchString(scanner.Text()) {
					result := demoP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Demo = strings.TrimSpace(result[0][1])
				}
				if sourceP.MatchString(scanner.Text()) {
					result := sourceP.FindAllStringSubmatch(scanner.Text(), -1)
					e.Source = strings.TrimSpace(result[0][1])
					e.Site = site
				} else {
					e.Source = site
				}
			}
		}
	}
	return entries
}

func toYaml(entries []*Entry) {
	/*	yamlFile, err := os.Create("./output.yaml")
		if err != nil {
			fmt.Println(err)
		}
		defer yamlFile.Close()
		YAML, err := yaml.Marshal(entries)
		if err != nil{
			fmt.Println("error:", err)
		}
		yamlFile.Write(YAML)
		yamlFile.Close()*/

	jsonFile, err := os.Create("./output.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	//JSON, err := json.MarshalIndent(entries, "", "\t")
	JSON, err := json.Marshal(entries)
	//JSON, err := json.Marshal(entries)
	if err != nil {
		fmt.Println("error:", err)
	}
	jsonFile.Write(JSON)
	jsonFile.Close()
}

func langSplit(lang string) []string {

	if strings.Contains(lang, "/") {
		return strings.Split(lang, "/")
	} else {
		l := make([]string, 1)
		l[0] = lang
		return l
	}
}

var tags = map[string][]string{
	"Analytics": {"Analytics"},
	"Web Analytics": {"Web Analytics"},
	"Archiving and Digital Preservation (DP)": {"archiving", "Digital Preservation"},
	"Automation": {"Automation"},
	"Blogging Platforms":                      {"Blog"},
	"Bookmarks and Link Sharing":              {"Bookmarks", "Links"},
	"Calendaring and Contacts Management":     {"Calendar", "Contacts"},
	"CalDAV or CardDAV servers":               {"CalDAV", "CardDav"},
	"Communication systems":                   {"Communications"},
	"Custom communication systems":            {},
	"Email":                                   {"Email"},
	"Complete solutions":                      {"Complete Email"},
	"Mail Transfer Agents":                    {"MTA"},
	"MTAs / SMTP servers":                     {"SMTP"},
	"Mail Delivery Agents":                    {"MDA"},
	"MDAs - IMAP/POP3 software":               {"IMAP", "POP3"},
	"Mailing lists and Newsletters":           {"Mailng List", "Newsletters"},
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
	"Search Engines": {"Search Engine"},
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
