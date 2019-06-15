package util

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
	"bytes"
)

// Gh struct receives Github api data
type Gh struct {
	Data struct {
		Repository struct {
			Stargazers struct {
				TotalCount int `json:"totalCount"`
			} `json:"stargazers"`
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
		Message string `json:"message"`
		Type    string `json:"type"`
		Path    string `json:"path"`
	} `json:"errors"`
}

func getGHRepo(ur, ght string) (int, string) {
	if strings.Contains(ur, "/") != true {
		log.Printf("No repository provided. Update %s to include a repo.", ur)
		return 0, ""
	}
	r := strings.Split(ur, "/")
	if r[1] == "" {
		log.Printf("No repository provided. Update %s to include a repo.", ur)
		return 0, ""
	}
	var jsonStr = []byte(`{"query":"{repository(owner:\"` + r[0] + `\",name:\"` + r[1] + `\"){stargazers{totalCount}defaultBranchRef{target{... on Commit{history(first: 1){edges{node{committedDate}}}}}}}}"}`)

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
		return ghv3api(r[0], r[1], ght)
	}
	//return *repo.StargazersCount
	return gh.Data.Repository.Stargazers.TotalCount, strings.Split(gh.Data.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate, "T")[0]
}

func ghv3api(u, r, ght string) (int, string) {
	ghURL := "https://api.github.com/repos/" + u + "/" + r + "?" //+ ght

	//fmt.Println(ghURL)
	// request http api
	req, err := http.NewRequest("GET", ghURL, nil)
	req.Header.Set("Authorization", "bearer "+ght)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Printf("StatusCode: %v. Source link is invalid: %v/%v\n", res.StatusCode, u, r)
		return 0, ""
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
	}
	gh := Gh{}
	err = json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Update https://github.com/%s/%s Source code link to: https://github.com/%s\n", u, r, gh.FullName)
	if "/"+u+"/"+r != gh.FullName {
		log.Printf("Retrying Github api v4 with %s", gh.FullName)
		getGHRepo(gh.FullName, ght)
	}
	return gh.Stars, strings.Split(gh.Updated, "T")[0]
}