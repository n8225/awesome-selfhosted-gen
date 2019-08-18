package exporter

import (
	"encoding/json"
	"os"
	"fmt"

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
	jsonFile.Write(JSON)
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
	jsonFileMin.Write(JSONmin)
	jsonFileMin.Close()
}
