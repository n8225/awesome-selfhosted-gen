package parse

import (
	"sort"
)

//Tags Struct
type Tags struct {
	Tag   string `json:"Tag"`
	Count int    `json:"C"`
}

// MakeTags creates the tag struct with counts
func MakeTags(entries []Entry) []Tags {
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
			tagMap[item]++
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

//Tagmap is used to create tags from category labels
var Tagmap = map[string][]string{
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
	"Knowledge Management Tools":               {"Knowledge Mgmnt"},
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
