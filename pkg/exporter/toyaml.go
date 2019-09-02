package exporter

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"

	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

// ToYAML creates/replaces the yaml list in the main dir
func ToYAML(list parse.List, fileName string) {
	yamlFile, err := os.Create("./" + fileName + ".yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer yamlFile.Close()
	YAML, err := yaml.Marshal(list)
	if err != nil {
		fmt.Println("error:", err)
	}
	_, err = yamlFile.Write(YAML)
	if err != nil {
		fmt.Println(err)
	}
	yamlFile.Close()
}

// ToYamlFiles creates directories and individual yaml files named from source url
func ToYamlFiles(list parse.List) {
	var i = 0
	for _, e := range list.Entries {
		// u, err := url.Parse(e.Source)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fname := strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(strings.TrimSuffix(u.Host+"_"+strings.TrimPrefix(u.EscapedPath(), "/"), "/"), "_"), "/", "_"), ":", "_")
		fname := fileName(e.Name)
		if fileExists("list/"+fname+".yaml") == true {
			fname = fname + "2"
			fmt.Printf("%d: %s already exists. Renaming to %s.\n", e.Line, e.Name, fname)
		}
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
		_, err = yamlFile.Write(YAML)
		if err != nil {
			fmt.Println(err)
		}
		yamlFile.Close()
		if fileExists("list/"+fname+".yaml") != true {
			fmt.Printf("Failed to write %d: %s\n", e.Line, e.Name)
		} else {
			i++
		}
	}
	fmt.Printf("Wrote %d of %d files\n", i, len(list.Entries))
}

func fileName(f string) string {

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(f, "")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
