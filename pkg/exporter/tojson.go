package exporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

// ToJSON exports the list into readable and  minified json files in the static dir
func ToJSON(list parse.List, fileName string) {
	jsonFile, err := os.Create("./static/" + fileName + ".json")
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
	_, err = jsonFile.Write(JSON)
	if err != nil {
		fmt.Println(err)
	}
	jsonFile.Close()

	jsonFileMin, err := os.Create("./static/" + fileName + ".min.json")
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
	_, err = jsonFileMin.Write(JSONmin)
	if err != nil {
		fmt.Println(err)
	}
	jsonFileMin.Close()
}
