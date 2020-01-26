package parse

import (
	"bufio"

	"os"
	"strings"
)

type tmpCat struct{
	start int
	stop int
	level int
	cat string
}

var prefixLvl = map[string]int{
	"_See": 0,
	"##": 1,
	"###": 2,
	"####": 3,
	"_": 4,
	"-[": 5,
}

var catNeedStop = make(map[int]int)

//MdParser parses the README.md file into lines creates a slice of entries
func MdParser(path, gh string) []Entry {
	inputFile, _ := os.Open(path)
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	l, i := 0, 0
	list := false
	entries := []Entry{}
	CatAarr := []tmpCat{}
	var tmpListMap = make(map[int]string)
	for scanner.Scan() {
		l++
		if strings.HasPrefix(scanner.Text(), "<!-- BEGIN SOFTWARE LIST -->") {
			list = true
		} else if strings.HasPrefix(scanner.Text(), "<!-- END SOFTWARE LIST -->") {
			addStop(CatAarr, catNeedStop, 1, l)
			list = false
		}
		if list {
			switch true {
			case prefixLvl[prefixParse(scanner.Text())] == 5: //adds entry line to a map of entries and their line #s
				tmpListMap[l] = scanner.Text()
			case prefixLvl[prefixParse(scanner.Text())] == 1:
				addStop(CatAarr, catNeedStop, 1, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[prefixParse(scanner.Text())], cat: scanner.Text()[3:len(scanner.Text())]})
				catNeedStop[1] = len(CatAarr)-1
			case prefixLvl[prefixParse(scanner.Text())] == 2:
				addStop(CatAarr, catNeedStop, 2, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[prefixParse(scanner.Text())], cat: scanner.Text()[4:len(scanner.Text())]})
				catNeedStop[2] = len(CatAarr)-1
			case prefixLvl[prefixParse(scanner.Text())] == 3:
				addStop(CatAarr, catNeedStop, 3, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[prefixParse(scanner.Text())], cat: scanner.Text()[5:len(scanner.Text())]})
				catNeedStop[3] = len(CatAarr)-1
			case prefixLvl[prefixParse(scanner.Text())] == 4:
				addStop(CatAarr, catNeedStop, 4, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[prefixParse(scanner.Text())], cat: scanner.Text()[1:len(scanner.Text())-1]})
				catNeedStop[4] = len(CatAarr)-1
			}
		}
	}

	for k, v := range tmpListMap {
		i++
		entries = append(entries, *AddEntry(i, k, v, CatAarr))
	}
	return entries
}

//addStop adds the last line to the category in the category map
func addStop(c []tmpCat, catNeedStop map[int]int, i, l int) {
	for k, v := range catNeedStop {
		if k >= i {
			c[v].stop = l
			delete(catNeedStop, k)
		}
	}
}

//prefixParse parses category to match switch
func prefixParse(s string) string {
	noAZ := func(r rune) rune {
		if (r != '#' && r!= '_' && r!= '-' && r!= '[') {
			return -1
		}
		return r
	}
	if s == "" {
		return s
	} else if s[0:4] == "_See" {
		return ""
	} else {
		return strings.Map(noAZ, s[0:3])
	}
}