package main

import (
	"bufio"
	"os"
	"strings"
)

// Entry is the structure of each entry
type Entry struct {
ID      int      `yaml:"ID" json:"ID"`
Name    string   `yaml:"N" json:"N"`
Descrip string   `yaml:"D,flow" json:"D"`
Source  string   `yaml:"Src,omitempty" json:"Sr,omitempty"`
Demo    string   `yaml:"Demo,omitempty" json:"Dem,omitempty"`
Site    string   `yaml:"Site,omitempty" json:"Si,omitempty"`
License []string `yaml:"Lic" json:"Li"`
Lang    []string `yaml:"Lang" json:"La"`
Cat     string   `yaml:"Cat" json:"C"`
Tags    []string `yaml:"T" json:"T"`
Free    bool     `yaml:"Free,omitempty" json:"F,omitempty"`
Pdep    bool     `yaml:"Pdep,omitempty" json:"P,omitempty"`
}




func getMD(path string) []string {
	list := false
	lineArray := []string
	inputFile, _ := os.Open(path)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "<!-- BEGIN SOFTWARE LIST -->") {
			list = true
		} else if strings.HasPrefix(scanner.Text(), "<!-- END SOFTWARE LIST -->") {
			list = false
		}
		if list == true {
			lineArray = append(lineArray, scanner.Text())
		}
	}



	return lineArray
}