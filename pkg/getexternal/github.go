package getexternal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Gh struct receives Github api data
type Gh struct {
	Data struct {
		Repository struct {
			Stargazers struct {
				TotalCount int `json:"totalCount"`
			} `json:"stargazers"`
			PrimaryLanguage struct {
				Name string `json:"name"`
			} `json:"primaryLanguage"`
			LicenseInfo struct {
				SpdxID string `json:"spdxId"`
			} `json:"licenseInfo"`
			DefaultBranchRef struct {
				Target struct {
					History struct {
						Edges []struct {
							Node struct {
								CommittedDate string `json:"committedDate"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"history"`
				} `json:"target"`
			} `json:"defaultBranchRef"`
		} `json:"repository"`
	} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Type    string   `json:"type"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type GHnoRepo struct {
	Data struct {
		RepositoryOwner struct {
			PinnedRepositories struct {
				Edges []struct {
					Node struct {
						Name string `json:"nameWithOwner"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"pinnedRepositories"`
			Repositories struct {
				Edges []struct {
					Node struct {
						Name string `json:"nameWithOwner"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"Repositories"`
		} `json:"repositoryOwner"`
	} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Type    string   `json:"type"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

// GetGHRepo uses the github APIv4 (GRAPHQL) to retrieve star count, last commit, license, and language info. If the repository has moved or there is an error it will fall back to APIv3.
func GetGHRepo(ur, ght, src string) (stars int, commitDate, license, language, srcUpdate string) {
	ur = strings.TrimPrefix(strings.TrimSuffix(ur, "/"), "/")
	if strings.Contains(ur, "/") != true {
		//log.Printf("No repository provided. Update %s to include a repo.", ur)
		fmt.Println("Sending " + ur + " to repo chooser.")
		ur = chooseRepo(ur, ght)
		fmt.Println("fixed: " + ur)
		//return 0, "", "", "", src
	}
	r := strings.Split(ur, "/")
	//fmt.Println(r)
	if r[1] == "" {
		//log.Printf("No repository provided. Update %s to include a repo.", ur)
		//return 0, "", "", "", src
		fmt.Println("Sending " + ur + " to repo chooser.")
		r = strings.Split(chooseRepo(ur, ght), "/")
	}

	var jsonStr = []byte(`{"query":"{repository(owner:\"` + r[0] + `\",name:\"` + r[1] + `\"){stargazers{totalCount}licenseInfo{spdxId}primaryLanguage{name}defaultBranchRef{target{... on Commit{history(first: 1){edges{node{committedDate}}}}}}}}"}`)
	//log.Printf("Body: %s\n", jsonStr)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "bearer "+ght)
	req.Header.Set("Accept", "application/vnd.github.quicksilver-preview+json")
	req.Header.Set("Content-Type", "application/json")
	//log.Printf("req: ", req.Body)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	//log.Printf("Body: %s\n", body)

	gh := &Gh{}
	err = json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal(err)
	}
	//Fallback to github v3 api on error
	if gh.Data.Repository.Stargazers.TotalCount == 0 {
		log.Printf("Error: https://github.com/%s/%s => %v Trying github api v3.  ", r[0], r[1], gh.Errors[0].Message)
		return ghv3api(r[0], r[1], ght, srcUpdate)
	}
	return gh.Data.Repository.Stargazers.TotalCount, strings.Split(gh.Data.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate, "T")[0], gh.Data.Repository.LicenseInfo.SpdxID, gh.Data.Repository.PrimaryLanguage.Name, src
}

func ghv3api(u, r, ght, src string) (stars int, commitDate, license, language, srcUpdate string) {
	ghURL := "https://api.github.com/repos/" + u + "/" + r + "?"

	req, err := http.NewRequest("GET", ghURL, nil)
	req.Header.Set("Authorization", "bearer "+ght)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Printf("StatusCode: %v. Source link is invalid: %v/%v\n", res.StatusCode, u, r)
		return 0, "", "", "", srcUpdate
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Body: %s\n", body)

	type Gh struct {
		FullName string `json:"full_name"`
		Stars    int    `json:"stargazers_count"`
		Created  string `json:"created_at"`
		Updated  string `json:"pushed_at"`
		License  struct {
			SpdxID string `json:"spdx_id"`
		} `json:"license"`
		Language string `json:"language"`
	}
	gh := Gh{}
	err = json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Update https://github.com/%s/%s Source code link to: https://github.com/%s\n", u, r, gh.FullName)
	if "/"+u+"/"+r != gh.FullName {
		log.Printf("Retrying Github api v4 with %s", gh.FullName)
		GetGHRepo(gh.FullName, ght, gh.FullName)
	}
	return gh.Stars, strings.Split(gh.Updated, "T")[0], gh.License.SpdxID, gh.Language, gh.FullName
}

func chooseRepo(ur, ght string) (url string) {
	//fmt.Println(ur)
	var jsonStr = []byte(`{"query":"{repositoryOwner(login:\"` + ur + `\"){pinnedRepositories(first: 1, orderBy: {direction: DESC, field: STARGAZERS}){edges{node{nameWithOwner}}}repositories(last:1, orderBy: {direction: ASC, field: STARGAZERS}){edges{node{nameWithOwner}}}}}"}`)
	//log.Printf("Body: %s\n", jsonStr)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "bearer "+ght)
	req.Header.Set("Accept", "application/vnd.github.quicksilver-preview+json")
	req.Header.Set("Content-Type", "application/json")
	//log.Printf("req: ", req.Body)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	//log.Printf("Body: %s\n", body)

	gh := &GHnoRepo{}
	err = json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	if len(gh.Data.RepositoryOwner.PinnedRepositories.Edges) == 0 {
		fmt.Println("ln24 true " + ur)
		return gh.Data.RepositoryOwner.Repositories.Edges[0].Node.Name
	}
	return gh.Data.RepositoryOwner.PinnedRepositories.Edges[0].Node.Name
}
