package exporter

import (
	"gopkg.in/yaml.v2"
	"os"
	"fmt"

	"github.com/n8225/ash_gen/pkg/parse"
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