package exporter

import (
	"encoding/json"
	"os"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
	"github.com/rs/zerolog/log"
)

// ToJSON exports the list into readable and  minified json files in the static dir
func ToJSON(list parse.List, fileName string) {
	jsonFile, err := os.Create("./" + fileName + ".json")
	if err != nil {
		log.Error().Stack().Err(err)
	}
	defer jsonFile.Close()
	JSON, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		log.Error().Stack().Err(err)
		return
	}

	_, err = jsonFile.Write(JSON)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	jsonFile.Close()

	jsonFileMin, err := os.Create("./" + fileName + ".min.json")
	if err != nil {
		log.Error().Stack().Err(err)
	}
	defer jsonFile.Close()
	JSONmin, err := json.Marshal(list)
	if err != nil {
		log.Error().Stack().Err(err)
		return
	}

	_, err = jsonFileMin.Write(JSONmin)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	jsonFileMin.Close()
}
