package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"flag"

	"gopkg.in/yaml.v2"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
	"github.com/n8225/awesome-selfhosted-gen/pkg/exporter"
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
	var yamlFiles, addedFiles int = 0, 0
	newID := parse.GetHighestID(l.Entries) + 1
	fmt.Println("There are ", len(l.Entries), " entries on the list. Next ID will be: ", newID)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") == true {
			e := parse.Entry{}
			yamlFiles++
			fFile, err := os.Open("./add/" + f.Name())
			if err != nil {
				fmt.Println(err)
			}
			fbytes, _ := ioutil.ReadAll(fFile)
			err = yaml.Unmarshal(fbytes, &e)
			e = parse.CheckEntry(e, l, *ghToken)
			
			if e.Error == true {
				log.Println(e.Errors)
				log.Println(e.Warns)
			} else {
				e.ID = newID
				l.Entries = append(l.Entries, e)
				fmt.Println("Adding ", e.Name, " with ID ", newID)
				addedFiles++
				newID++
			}
		} 
	}
	fmt.Printf("Added %d of %d yaml files found. There are now %d entries on the list.", addedFiles, yamlFiles, len(l.Entries))
	l.TagList = parse.MakeTags(l.Entries)
	l.LangList = parse.MakeLangs(l.Entries)
	l.CatList = parse.MakeCats(l.Entries)

	/*	for _, e := range l.Entries {
		getStars(e, gh)
	}*/

	exporter.ToJSON(l, "list")
	exporter.ToYaml(l, "list")

}
