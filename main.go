package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/n8225/awesome-selfhosted-gen/internal/pkg/exporter"
	"github.com/n8225/awesome-selfhosted-gen/internal/pkg/parse"
)

func main() {
	var path string
	const (
		defaultPath = ""
		usage       = "Path to Readme.md(On windows wrap path in \""
	)
	flag.StringVar(&path, "path", defaultPath, usage)
	flag.StringVar(&path, "p", defaultPath, usage)
	var ghToken = flag.String("ghtoken", "", "github oauth token")
	//var glToken = flag.String("gltoken", "", "gitlab oauth token")
	flag.Parse()
	apath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	e := parse.MdParser(apath, *ghToken)
	
	l := new(parse.List)
	
	l.Entries = e
	exporter.ToYamlFiles(*l)
	yl := exporter.ImportYaml(*ghToken)

	exporter.ToJSON(yl, "list")
	exporter.ToYAML(yl, "list")
	exporter.MapToJSON(yl.CatIDs, "catids")
	exporter.MapToJSON(yl.LangIDs, "langids")
	exporter.MapToJSON(yl.TagIDs, "tagids")

}
