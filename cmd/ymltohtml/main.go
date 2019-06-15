package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/ash_gen/pkg/util"

)

func main() {
	inputFile, err := os.Open("./list.yaml")
	if inputFile != nil {
		defer func() {
			if ferr := inputFile.Close(); ferr != nil {
				err = ferr
			}
		}()
	}
	ybytes, _ := ioutil.ReadAll(inputFile)
	l := List{}
	err = yaml.Unmarshal(ybytes, &l)
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir("./add")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix("./add/"+f.Name(), ".yaml") {
			e := Entry{}
			fFile, err := os.Open("./add/" + f.Name())
			if err != nil {
				fmt.Println(err)
			}
			fbytes, _ := ioutil.ReadAll(fFile)
			err = yaml.Unmarshal(fbytes, &e)
			e = CheckYamlAdd(e, l, gh)
			log.Println(e.Errors)
			log.Println(e.Warns)
			if e.Errors == nil {
				l.Entries = append(l.Entries, e)
			}
		}
	}
	l.TagList = makeTags(l.Entries)
	l.LangList = makeLangs(l.Entries)

	/*	for _, e := range l.Entries {
		getStars(e, gh)
	}*/

	jsonFile, err := os.Create("./outputfy.json")
	if jsonFile != nil {
		defer func() {
			if ferr := jsonFile.Close(); ferr != nil {
				err = ferr
			}
		}()
	}
	if err != nil {
		log.Fatal(err)
	}
	JSON, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(string(JSON))
	_, err = jsonFile.Write(JSON)
	if err != nil {
		fmt.Println(err)
	}

}
