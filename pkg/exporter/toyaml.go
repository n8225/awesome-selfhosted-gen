package exporter

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

// ToYaml creates/replaces the yaml list in the main dir
func ToYaml(list parse.List, fileName string) {
	yamlFile, err := os.Create("./" + fileName + ".yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer yamlFile.Close()
	YAML, err := yaml.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
	}
	yamlFile.Write(YAML)
	yamlFile.Close()
}

// ToYamlFiles creates directories and individual yaml files named from source url
func ToYamlFiles(list parse.List) {
	// for _, d := range list.CatList {
	// 	var dname = d.Cat
	// 	dname = strings.ReplaceAll(dname, "/", "&")
	// 	dname = strings.ReplaceAll(dname, " ", "")
	// 	if _, err := os.Stat("list/" + dname); os.IsNotExist(err) {
	// 		err = os.MkdirAll("list/"+dname, 0755)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }
	for _, e := range list.Entries {
		var fname = e.Source
		var dname = e.Cat
		fname = strings.TrimSuffix(fname, "/")
		fname = strings.TrimPrefix(fname, "http://")
		fname = strings.TrimPrefix(fname, "https://")
		fname = strings.TrimPrefix(fname, "www.")
		fname = strings.ReplaceAll(fname, " ", "")
		fname = strings.ReplaceAll(fname, "/", "-")
		dname = strings.ReplaceAll(dname, "/", "&")
		dname = strings.ReplaceAll(dname, " ", "")

		yamlFile, err := os.Create("list/" + fname + ".yaml")
		if err != nil {
			fmt.Println(err)
		}
		defer yamlFile.Close()
		YAML, err := yaml.Marshal(e)
		if err != nil {
			fmt.Println("error:", err)
		}
		//fmt.Printf("%s\n\n", string(YAML))
		yamlFile.Write(YAML)
		yamlFile.Close()
	}

}
