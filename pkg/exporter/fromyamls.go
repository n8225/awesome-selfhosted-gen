package exporter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"

	"github.com/n8225/awesome-selfhosted-gen/pkg/getexternal"
	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

//ImportYaml imports yaml files into a slice of structs
func ImportYaml(ght string) (l parse.List) {
	files, err := ioutil.ReadDir("./list")
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	var i = 0
	c := make(chan parse.Entry)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			i++

			go processYaml(f, ght, c)
		}
		l.Entries = append(l.Entries, <-c)
	}
	log.Info().Msgf("Added %d of %d yaml files found. There are now %d entries on the list.", i, len(files), len(l.Entries))
	l.TagList = parse.MakeTags(l.Entries)
	l.LangList = parse.MakeLangs(l.Entries)
	l.CatList = parse.MakeCats(l.Entries)

	return
}

func processYaml(f os.FileInfo, ght string, c chan parse.Entry) {
	e := parse.Entry{}

	fFile, err := os.Open("./list/" + f.Name())
	if err != nil {
		log.Error().Stack().Err(err)
	}
	fbytes, _ := ioutil.ReadAll(fFile)
	err = yaml.Unmarshal(fbytes, &e)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	log.Debug().Msgf("%d - %s Imported from %s", e.ID, e.Name, f.Name())
	//g := parse.GetGitdata(e, ght)
	c <- e
}

type ghDataMap struct {
	Data map[string]struct {
		PushedAt        string `json:"pushedAt"`
		UpdatedAt       string `json:"updatedAt"`
		IsArchived      bool   `json:"isArchived"`
		IsDisabled      bool   `json:"isDisabled"`
		NameWithOwner   string `json:"nameWithOwner"`
		StargazerCount  int    `json:"starGazerCount"`
		Cost            int    `json:"cost"`
		Remaining       int    `json:"remaining"`
		ResetAt         string `json:"resetAt"`
		PrimaryLanguage struct {
			Name string `json:"name"`
		} `json:"primaryLanguage"`
		LicenseInfo struct {
			Key string `json:"key"`
		} `json:"licenseInfo"`
	} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Type    string   `json:"type"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

func GetGithubData(l []parse.Entry, ght string) map[int]parse.GithubRepo {
	ghLinks := make(map[int]parse.GithubRepo)
	ghLinkRe := regexp.MustCompile(`https://github\.com/([a-zA-Z\d\-\._]{1,39})/([a-zA-Z\d\-\._]{1,39})(/.*)?`)
	for _, e := range l {
		if ghLinkRe.MatchString(e.Source) {
			res := ghLinkRe.FindStringSubmatch(e.Source)
			ghLinks[e.ID] = parse.GithubRepo{
				Name: res[1],
				Repo: res[2],
			}

		}
	}
	//var jsonStr = []byte(`{"query":"{repository(owner:\"` + r[0] + `\",name:\"` + r[1] + `\"){stargazers{totalCount}defaultBranchRef{target{... on Commit{history(first: 1){edges{node{committedDate}}}}}}}}"}`)
	queryParam := `{pushedAt updatedAt isArchived isDisabled nameWithOwner primaryLanguage{name} licenseInfo{key} stargazerCount} `
	called := 0
	staged := 0
	reqPerCall := 100
	jsonStrBuild := `{"query":"{`
	for k, v := range ghLinks {
		//jsonStrBuild := `{"query":"{rateLimit{cost remaining resetAt}`
		jsonStrBuild += `Q` + strconv.Itoa(k) + `: repository(owner:\"` + v.Name + `\" name:\"` + v.Repo + `\")` + queryParam
		staged++
		called++
		if staged == reqPerCall || called == len(ghLinks) {
			jsonStrBuild += `}"`
			jsonStr := []byte(jsonStrBuild)
			ghData := ghDataMap{}
			err := json.Unmarshal(getexternal.GhClientv4(ght, jsonStr), &ghData)
			if err != nil {
				log.Fatal().Stack().Err(err).Stack().Err(err)
			}
			// if ghData.Errors != nil {
			// 	log.Print(ghData.Errors)
			// }
			//log.Debug().Msgf("%#v/n", ghData)
			//log.Debug().Msgf("%v/n", ghData)
			staged = 0
			jsonStrBuild = `{"query":"{`
			for l, u := range ghData.Data {
				if l == "rateLimit" {
					log.Debug().Msgf("Github API Remaining: %d", u.Remaining)
				} else {
					curID, err := strconv.Atoi(strings.Replace(l, "Q", "", 1))
					if err != nil {
						log.Error().Err(err)
					}
					ghLinks[curID] = parse.GithubRepo{
						Name:          ghLinks[curID].Name,
						Repo:          ghLinks[curID].Repo,
						Stars:         u.StargazerCount,
						NameWithOwner: u.NameWithOwner,
						Lang:          u.PrimaryLanguage.Name,
						Spdx:          u.LicenseInfo.Key,
						PushedAt:      u.PushedAt,
						Archived:      u.IsArchived,
						Disabled:      u.IsDisabled,
					}
				}
			}
		}

	}
	return ghLinks
}

// GHDataToYAML creates a yaml file containing github data
func GhDataToYAML(data map[int]parse.GithubRepo, path string) {
	yamlFile, err := os.Create("./" + path + "/github-data.yaml")
	if err != nil {
		log.Error().Stack().Err(err)
	}
	defer yamlFile.Close()
	YAML, err := yaml.Marshal(data)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	_, err = yamlFile.Write(YAML)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	yamlFile.Close()
}
