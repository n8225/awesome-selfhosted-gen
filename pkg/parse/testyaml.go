package parse

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"fmt"
	"unicode"

	"github.com/n8225/ash_gen/pkg/getexternal"
)

// RepoCheckStruct is struct for
type RepoCheckStruct struct {
	License  string
	Language string
	Updated  string
}

// Spdx struct is struct for imported SPDX license list
type Spdx struct {
	LicenseListVersion string `json:"licenseListVersion"`
	Licenses           []struct {
		Reference             string   `json:"reference"`
		IsDeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
		DetailsURL            string   `json:"detailsUrl"`
		ReferenceNumber       string   `json:"referenceNumber"`
		Name                  string   `json:"name"`
		LicenseID             string   `json:"licenseId"`
		SeeAlso               []string `json:"seeAlso"`
		IsOsiApproved         bool     `json:"isOsiApproved"`
	} `json:"licenses"`
	ReleaseDate string `json:"releaseDate"`
}

// CheckEntry runs tests on entries
func CheckEntry(e Entry, l List, ght string) Entry {
	e.Error = false
	ghR := RepoCheckStruct{}

	if e.Name == "" {
		e.Error = true
		e.Errors = append(e.Errors, "Error: Name is null.")
	}
	if dErr, dErrs, dWarn := checkDup(e, l); dErr == true {
		fmt.Println("ln47: ", dErr, dWarn)
		e.Error = true
		e.Errors = append(e.Errors, dErrs)
		e.Warns = append(e.Warns, dWarn)
	}
	if e.Source == "" {
		e.Error = true
		e.Errors = append(e.Errors, "Error: Source is null.")
	} else {
		dErr, dErrs, dWarn := checkLinks(e.Source, "Source Link ")
			if dErr == true {
					e.Errors = append(e.Errors, dErrs)
			}
		fWarn, src := "", ""
		fWarn, src, ghR = checkSource(e.Source, ght)
		if src != "" {
			e.Source = "https://www.github.com/" + src
		}
		e.Warns = append(e.Warns, dWarn, fWarn)

	}
	if e.License == nil {
		e.Error = true
		e.Errors = append(e.Errors, "Error: License is null.")
	} else if ghR.License != "" && ghR.License != e.License[0] {
		e.Warns = append(e.Warns, "Github reported license does not match first license in license array.")
	} else {
		dWarn := checkLicense(e.License)
		e.Warns = append(e.Warns, dWarn...)
	}
	if e.Lang == nil {
		e.Errors = append(e.Errors, "Error: Programming Language is null.")
	} else {
		dWarn := checkLang(e.Lang, l)
		e.Warns = append(e.Warns, dWarn...)
	}
	if e.Tags == nil {
		e.Errors = append(e.Errors, "Error: Tags is null.")
	} else {
		dWarn := checkTag(e.Tags, l)
		e.Warns = append(e.Warns, dWarn...)
	}
	if e.Descrip == "" {
		e.Errors = append(e.Errors, "Error: Description is null.")
	} else {
		dErr, dWarn := "", ""
		dErr, dWarn, e.Descrip = checkDesc(e.Descrip)
		e.Warns = append(e.Warns, dWarn)
		e.Errors = append(e.Errors, dErr)
	}
	if e.Demo != "" {
		dErr, dErrs, dWarn := checkLinks(e.Demo, "Demo Link ")
		if dErr == true {
					e.Errors = append(e.Errors, dErrs)
		}
		e.Warns = append(e.Warns, dWarn)
	}
	if e.Site != "" {
		dErr, dErrs, dWarn := checkLinks(e.Site, "Site Link ")
		if dErr == true {
			e.Errors = append(e.Errors, dErrs)
}
e.Warns = append(e.Warns, dWarn)
	}
	if e.Clients != nil {
		for _, c := range e.Clients {
			dErr, dErrs, dWarn := checkLinks(c, "Clients")
			if dErr == true {
				e.Errors = append(e.Errors, dErrs)
			}
			e.Warns = append(e.Warns, dWarn)
		}
	}
	return e
}
func checkDup(e Entry, l List) (eErr bool, eErrs , eWarn string) {
	for _, es := range l.Entries {
		if e.Source == es.Source || e.Site == es.Source || e.Source == es.Site || e.Site == es.Site {
			eErr = true
			eErrs = "Error: " + e.Name + ", " + e.Source + " is already on the List as" + es.Name + "; "
		}
		if e.Name == es.Name {
			eErr = true
			eWarn = "Warn: " + e.Name + ", is a possible duplicate; "
		}
	}
	return
}
func checkLinks(l, n string) (dErr bool, dErrs, dWarn string) {
	resp, err := http.Get(l)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	} else if resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		dWarn = "Warning: " + n + l + " => " + strconv.Itoa(resp.StatusCode) + http.StatusText(resp.StatusCode)
		return
	} else if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
		dErr = true
		dErrs = "Error: " + n + l + " => " + strconv.Itoa(resp.StatusCode) + http.StatusText(resp.StatusCode)
		return
	}
	return true, "Error: " + n + l + "Not Checked.", ""
}
func checkSource(l, ght string) (dWarn, src string, ghR RepoCheckStruct) {
	u, err := url.Parse(l)
	if err != nil {
		panic(err)
	}
	if u.Host == "github.com" {
		ghR = RepoCheckStruct{}
		gh := strings.TrimFunc(u.Path, func(r rune) bool {return !unicode.IsLetter(r) && !unicode.IsNumber(r)})
		_, ghR.Updated, ghR.License, ghR.Language, src = getexternal.GetGHRepo(gh, ght, gh)
		return "", src, ghR
	} else if u.Host == "gitlab.com" {
		_, updated := getexternal.GetGLRepo(l)
		ghR := RepoCheckStruct{
			Updated: updated,
		}
		return "", "", ghR
	} else {
		ghR := RepoCheckStruct{}
		return "Not github or gitlab, no source checks performed.", "", ghR
	}
}
func checkLicense(lic []string) (dWarn []string) {
	//out, err := os.Create("./spdx.json")
	res, err := http.Get("https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	newSpdx := Spdx{}
	//_, err = io.Copy(out, bytes.NewBuffer(body))
	err = json.Unmarshal(body, &newSpdx)
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range lic {
		if containSpdx(newSpdx, l) != true {
			dWarn = append(dWarn, "Warning: License '"+l+"' is not in the SPDX format or on the SPDX list(https://spdx.org/licenses/).")
		}
	}
	return
}

func containSpdx(s Spdx, e string) bool {
	for _, a := range s.Licenses {
		if a.LicenseID == e {
			return true
		}
	}
	return false
}
func checkLang(langs []string, l List) (dWarn []string) {
	for _, lang := range langs {
		if containLang(MakeLangs(l.Entries), lang) != true {
			dWarn = append(dWarn, "Warning: Programming Language'"+lang+"' is not in the current list(We should standardize these).")
		}
	}
	return dWarn
}
func containLang(l []Langs, lang string) bool {
	for _, a := range l {
		if a.Lang == lang {
			return true
		}
	}
	return false
}
func checkTag(tags []string, l List) (dWarn []string) {
	for _, tag := range tags {
		if containTag(MakeTags(l.Entries), tag) != true {
			dWarn = append(dWarn, "Warning: Tag'"+tag+"' is not in the current list(We should standardize these).")
		}
	}
	return dWarn
}
func containTag(l []Tags, tag string) bool {
	for _, a := range l {
		if a.Tag == tag {
			return true
		}
	}
	return false
}
func checkDesc(d string) (dWarn, dErr, newD string) {
	newD = d
	if strings.HasSuffix(newD, ".") != true {
		dWarn = "Missing fullstop(.) I added it for you."
		newD += "."
	}
	if len(newD) > 250 {
		dErr = "Description too long. Shorten it to 250 or fewer charachters."
	}
	return
}
