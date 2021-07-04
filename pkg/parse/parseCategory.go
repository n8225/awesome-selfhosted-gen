package parse

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type CatList struct {
	Id        int
	Category  string
	StartLine int
	EndLine   int
	Level     int
	SubCat    []int
	Comment   []string
	Entries   []int
}

func parseCategory(line string, state parseState, cats []CatList) (parseState, []CatList) {
	if strings.HasPrefix(line, "## ") {
		setPrev(state.lasts, state.l, cats)
		state.i++
		state.lasts = [3]int{state.i, 0, 0}
		state.lastCat = 0
		cats = append(cats, setCat(line, "## ", state, 1))
	}
	if strings.HasPrefix(line, "### ") {
		setPrev([3]int{0, state.lasts[1], state.lasts[2]}, state.l, cats)
		state.i++
		state.lasts[1], state.lasts[2] = state.i, 0
		cats[state.lasts[0]-1].SubCat = append(cats[state.lasts[0]-1].SubCat, state.i)
		state.lastCat = 1
		cats = append(cats, setCat(line, "### ", state, 2))
	}
	if strings.HasPrefix(line, "#### ") {
		setPrev([3]int{0, 0, state.lasts[2]}, state.l, cats)
		state.i++
		state.lasts[2] = state.i
		cats[state.lasts[1]-1].SubCat = append(cats[state.lasts[1]-1].SubCat, state.i)
		state.lastCat = 2
		cats = append(cats, setCat(line, "#### ", state, 3))
	}
	if strings.HasPrefix(line, "_") {
		cats[state.lasts[state.lastCat]-1].Comment = append(cats[state.lasts[state.lastCat]-1].Comment, strings.Trim(line, "_"))
	}
	return state, cats
}

func setCat(line, p string, state parseState, lvl int) (tmp CatList) {
	tmp.Id = state.i
	tmp.Category = strings.Trim(line, p)
	tmp.StartLine = state.l
	tmp.Level = lvl
	return
}

func setPrev(lasts [3]int, l int, tmpList []CatList) {
	for _, last := range lasts {
		if last != 0 {
			tmpList[last-1].EndLine = l - 1
		}
	}
}

func closeCats(state parseState, cats []CatList) []CatList {
	for _, last := range state.lasts {
		if last != 0 && cats[last-1].EndLine == 0 {
			cats[last-1].EndLine = state.l - 1
		}
	}
	return cats
}

func getCat(id int, state parseState, cats []CatList) ([]string, map[int]int) {
	var cat []string
	e := make(map[int]int)
	for _, l := range state.lasts {
		if l != 0 {
			cat = append(cat, cats[l-1].Category)
		}
	}
	catID := cats[findMax(state.lasts)-1].Id - 1
	e[catID] = id
	return cat, e
}

func findMax(d [3]int) (max int) {
	max = d[0]
	for _, value := range d {
		if value > max {
			max = value
		}
	}
	return
}

func getTags(cats []string) (tags []string) {
	for _, c := range cats {
		tags = append(tags, Tagmap[c]...)
	}
	return
}

func setCats(cat []string) (string, string, string) {
	if len(cat) == 3 {
		return cat[0], cat[1], cat[2]
	} else if len(cat) == 2 {
		return cat[0], cat[1], ""
	}
	return cat[0], "", ""
}

//Category defines the category structure
type Category struct {
	Entries []int               `json:"Entries,omitempty" yaml:"Entries,omitempty"`
	Comment []string            `json:"Comments,omitempty" yaml:"Comments,omitempty"`
	Subcat  map[string]Category `json:"SubCategories,omitempty" yaml:"SubCategories,omitempty"`
}

func toCategories(cats []CatList) {
	newCats := map[string]Category{}
	for _, c := range cats {
		if c.Level == 1 {
			newCats[c.Category] = Category{
				Comment: c.Comment,
				Subcat:  getSub2(cats, c.SubCat),
				Entries: c.Entries,
			}
		}
	}
	catToYaml(&newCats)
}

func catToYaml(cat *map[string]Category) {
	yamlFile, err := os.Create("./output/categories.yaml")
	if err != nil {
		log.Error().Stack().Err(err)
	}
	defer yamlFile.Close()
	YAML, err := yaml.Marshal(&cat)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	_, err = yamlFile.Write(YAML)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	yamlFile.Close()
}

func getSub2(cats []CatList, subcats []int) map[string]Category {
	sub2 := map[string]Category{}
	for _, c := range subcats {
		sub2[cats[c-1].Category] = Category{
			Comment: cats[c-1].Comment,
			Subcat:  getSub3(cats, cats[c-1].SubCat),
			Entries: cats[c-1].Entries,
		}
	}
	return sub2
}

func getSub3(cats []CatList, subcats []int) map[string]Category {
	sub3 := map[string]Category{}
	for _, c := range subcats {
		sub3[cats[c-1].Category] = Category{
			Comment: cats[c-1].Comment,
			Entries: cats[c-1].Entries,
		}
	}
	return sub3
}
