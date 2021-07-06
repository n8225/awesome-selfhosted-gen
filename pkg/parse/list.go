package parse

import (
	"bufio"
	"regexp"

	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// List is the total struct
type List struct {
	Entries  []Entry  `json:"Entries" yaml:"Entries"`
	LangList []Langs  `json:"Langs" yaml:"-"`
	CatList  []Cats   `json:"Cats" yaml:"-"`
	TagList  []Tags   `json:"Tags" yaml:"-"`
	Header   []string `json:"Header" yaml:"-"`
	Licenses []string `json:"Licenses" yaml:"-"`
	ExtLinks []string `json:"ExtLinks" yaml:"-"`
	Footer   []string `json:"Footer" yaml:"-"`
}
type parseState struct {
	i       int
	l       int
	lasts   [3]int
	lastCat int
	section string
}

//MdParser parses entries and categories from README.me
func MdParser(path string) *List {
	list := List{}
	log.Info().Msgf("Parsing: %s", path)
	inputFile, _ := os.Open(path)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	var i = 0

	state := parseState{}
	cats := []CatList{}

	//catLookup := make(map[int]int)

	for scanner.Scan() {
		state.section = findSection(scanner.Text(), state.section)
		state.l++
		if state.section == "list" {

			if strings.HasPrefix(scanner.Text(), "#") || strings.HasPrefix(scanner.Text(), "_") {
				state, cats = parseCategory(scanner.Text(), state, cats)
			}

			if strings.HasPrefix(scanner.Text(), "- [") || strings.HasPrefix(scanner.Text(), "  - [") {
				if regexp.MustCompile(Pattern).MatchString(scanner.Text()) {
					i++
					e := new(Entry)
					cat, catEntries := getCat(i, state, cats)
					for k, v := range catEntries {
						cats[k].Entries = append(cats[k].Entries, v)
					}
					e.Line = state.l
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
					e.Cat, e.Cat2, e.Cat3 = setCats(cat)
					e.Tags = getTags(cat)
					list.Entries = append(list.Entries, *e)
				} else {
					log.Error().Msgf("Failed to match pattern, Line: %d : %s", state.l, scanner.Text())
				}

			}
		} else if state.section == "licenseList" {
			cats = closeCats(state, cats)
		}

	}
	toCategories(cats)
	return &list
}

func findSection(line string, section string) string {
	switch true {
	case strings.HasPrefix(line, "# Awesome-Selfhosted"):
		return "header"
	case strings.HasPrefix(line, "- List of Software"):
		return "toc"
	case strings.HasPrefix(line, "<!-- BEGIN SOFTWARE LIST -->"):
		return "list"
	case strings.HasPrefix(line, "- List of Licenses"):
		return "licenseList"
	case strings.HasPrefix(line, "## External links"):
		return "extLinks"
	case strings.HasPrefix(line, "## Contributing"):
		return "footer"
	default:
		return section
	}
}
