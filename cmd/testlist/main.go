package main

import (
	"io/ioutil"
	"log"
	"os"
	"flag"

	"gopkg.in/yaml.v2"

	"github.com/n8225/ash_gen/pkg/parse"
)
func main() {
	inputFile, err := os.Open("./list.yaml")
	//var ghToken = flag.String("ghtoken", "", "github oauth token")
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
}