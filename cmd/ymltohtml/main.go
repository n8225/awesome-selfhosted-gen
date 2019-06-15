package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"flag"

	"gopkg.in/yaml.v2"

	"github.com/n8225/ash_gen/pkg/parse"


)

func main() {
	inputFile, err := os.Open("./list.yaml")
	var ghToken = flag.String("ghtoken", "", "github oauth token")
	flag.Parse()

	if inputFile != nil {
		defer func() {
			if ferr := inputFile.Close(); ferr != nil {
				err = ferr
			}
		}()
	}
	ybytes, _ := ioutil.ReadAll(inputFile)
	l := parse.List{}
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
			e := parse.Entry{}
			fFile, err := os.Open("./add/" + f.Name())
			if err != nil {
				fmt.Println(err)
			}
			fbytes, _ := ioutil.ReadAll(fFile)
			err = yaml.Unmarshal(fbytes, &e)
			e = parse.CheckYamlAdd(e, l, *ghToken)
			log.Println(e.Errors)
			log.Println(e.Warns)
			if e.Errors == nil {
				l.Entries = append(l.Entries, e)
			}
		}
	}
	l.TagList = parse.MakeTags(l.Entries)
	l.LangList = parse.MakeLangs(l.Entries)

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
