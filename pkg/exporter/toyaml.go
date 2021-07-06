package exporter

import (
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

// ToYAML creates/replaces the yaml list in the main dir
func ToYAML(list parse.List, fileName string) {
	yamlFile, err := os.Create("./" + fileName + ".yaml")
	if err != nil {
		log.Error().Stack().Err(err)
	}
	defer yamlFile.Close()
	YAML, err := yaml.Marshal(list)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	_, err = yamlFile.Write(YAML)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	yamlFile.Close()
}

// ToYamlFiles creates directories and individual yaml files named from source url
func ToYamlFiles(list parse.List) {
	var i = 0
	err := os.Mkdir("list", 0755)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	for _, e := range list.Entries {
		u, err := url.Parse(e.Source)
		if err != nil {
			log.Error().Stack().Err(err)
		}
		fname := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(strings.TrimSuffix(u.Host+"_"+strings.TrimPrefix(u.EscapedPath(), "/"), "/"), "_"), "/", "_"), ":", "_")
		//fname := fileName(e.Name)
		if fileExists("list/" + fname + ".yaml") {
			fname = fname + "2"
			log.Info().Msgf("%d: %s already exists. Renaming to %s.", e.Line, e.Name, fname)
		}
		yamlFile, err := os.Create("list/" + fname + ".yaml")
		if err != nil {
			log.Error().Stack().Err(err)
		}
		defer yamlFile.Close()
		YAML, err := yaml.Marshal(e)
		if err != nil {
			log.Error().Stack().Err(err)
		}
		_, err = yamlFile.Write(YAML)
		if err != nil {
			log.Error().Stack().Err(err)
		}
		yamlFile.Close()
		if !fileExists("list/" + fname + ".yaml") {
			log.Info().Msgf("Failed to write %d: %s", e.Line, e.Name)
		} else {
			i++
		}
	}
	log.Info().Msgf("Wrote %d of %d files", i, len(list.Entries))
}

// func fileName(f string) string {

// 	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return reg.ReplaceAllString(f, "")
// }

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func UpdateExtData(entries *[]parse.Entry, ghData map[int]parse.GithubRepo) {
	for i, e := range *entries {
		(*entries)[i].Stars = ghData[e.ID].Stars
		(*entries)[i].Updated = ghData[e.ID].PushedAt
	}
}
