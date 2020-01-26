package parse

import (
	"bufio"

	"os"
	"strings"
	"fmt"
)



type tmpCat struct{
	start int
	stop int
	level int
	cat string
}

var prefixLvl = map[string]int{
	"##": 1,
	"###": 2,
	"####": 3,
	"_": 4,
	"-[": 5,
}

var catNeedStop = make(map[int]int)


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
			case prefixLvl[PrefixParse(scanner.Text())] == 5:
				tmpListMap[l] = scanner.Text()
			case prefixLvl[PrefixParse(scanner.Text())] == 1:
				addStop(CatAarr, catNeedStop, 1, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[PrefixParse(scanner.Text())], cat: strings.Trim(scanner.Text(), "## ")})
				catNeedStop[1] = len(CatAarr)-1
			case prefixLvl[PrefixParse(scanner.Text())] == 2:
				addStop(CatAarr, catNeedStop, 2, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[PrefixParse(scanner.Text())], cat: strings.Trim(scanner.Text(), "### ")})
				catNeedStop[2] = len(CatAarr)-1
			case prefixLvl[PrefixParse(scanner.Text())] == 3:
				addStop(CatAarr, catNeedStop, 3, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[PrefixParse(scanner.Text())], cat: strings.Trim(scanner.Text(), "### ")})
				catNeedStop[3] = len(CatAarr)-1
			case prefixLvl[PrefixParse(scanner.Text())] == 4:
				addStop(CatAarr, catNeedStop, 4, l)
				CatAarr = append(CatAarr, tmpCat{start: l, level: prefixLvl[PrefixParse(scanner.Text())], cat: strings.Trim(scanner.Text(), "_")})
				catNeedStop[4] = len(CatAarr)-1
			}
		}
	}

	for k, v := range tmpListMap {
		fmt.Printf("%d: %v\n", k, v)
		entries = append(entries, *addEntry(i, k, v, CatAarr))
		i++
	}

	fmt.Printf("%+v\n", CatAarr)
	return entries
}

func addEntry(i, l int, t string, catAarr []tmpCat) *Entry {
	e := new(Entry)
	e.Line = l
	e.ID = i
	e.MD = t
	e.Name = GetName(e.MD)
	e.Descrip = GetDescrip(e.MD)
	e.License = GetLicense(e.MD)
	e.Lang = GetLang(e.MD)
	e.Pdep = GetPdep(e.MD)
	e.Demo = GetDemo(e.MD)
	e.Clients = GetClients(e.MD)
	e.Site = GetSite(e.MD)
	e.Source, e.SourceType = GetSource(e.MD)
	e.Cat, e.Tags = getCat(l, catAarr)
	return e
}

func getCat(l int, catAarr []tmpCat) (cat string, tags []string) {
	for _, c := range catAarr{
		if (c.level == 1 && l > c.start && l < c.stop) {
			cat = c.cat
			tags = append(tags, c.cat)
		} else if (c.start < l && l < c.stop) {
			tags = append(tags, c.cat)
		}
	}
	return
}

func addStop(c []tmpCat, catNeedStop map[int]int, i, l int) {
	for k, v := range catNeedStop {
		if k >= i {
			c[v].stop = l
			delete(catNeedStop, k)
		}
	}
}

func PrefixParse(s string) string {
	noAZ := func(r rune) rune {
		if (r != '#' && r!= '_' && r!= '-' && r!= '[') {
			return -1
		}
		return r
	}
	if s == "" {
		return s
	} else {
		return strings.Map(noAZ, s[0:3])
	}
	
}