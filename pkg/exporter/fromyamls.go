package exporter

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

//ImportYaml imports yaml files into a slice of structs
func ImportYaml(ght string) (l parse.List) {
	files, err := ioutil.ReadDir("./list")
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	var i = 0
	c := make(chan parse.Entry)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			i++

			go processYaml(f, ght, i, c)
		}
		l.Entries = append(l.Entries, <-c)
	}
	log.Info().Msgf("Added %d of %d yaml files found. There are now %d entries on the list.", i, len(files), len(l.Entries))
	l.TagList = parse.MakeTags(l.Entries)
	l.LangList = parse.MakeLangs(l.Entries)
	l.CatList = parse.MakeCats(l.Entries)

	return
}

func processYaml(f os.FileInfo, ght string, i int, c chan parse.Entry) {
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
	e.ID = i
	log.Debug().Msgf("%d - %s Imported from %s", e.ID, e.Name, f.Name())
	g := parse.GetGitdata(e, ght)

	e.Stars = g.Stars
	e.Updated = g.Updated
	e.Gitdata = *g
	c <- e
}
