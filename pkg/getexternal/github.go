package getexternal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Gh struct receives Github v4 api data
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

//GHnoRepo struct receives Github api data when we are guessing
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

//ghv3 is the struct for the v3 github api
type ghv3 struct {
	FullName string `json:"full_name"`
	Stars    int    `json:"stargazers_count"`
	Created  string `json:"created_at"`
	Updated  string `json:"pushed_at"`
	License  struct {
		SpdxID string `json:"spdx_id"`
	} `json:"license"`
	Language string `json:"language"`
}

//getOwnRepo parses the owner and repo from a github url
func getOwnRepo(ghURL, ght string) string {
	u, err := url.Parse(ghURL)
	if err != nil {
		fmt.Println(err)
	}
	return u.Path

}

//GetGH uses the github APIv4 (GRAPHQL) to retrieve star count, last commit, license, and language info. If the repository has moved or there is an error it will fall back to APIv3.
func GetGH(ghURL, ght string, hasErrors []string) (int, string, string, string, []string) {

	ownRepo := strings.TrimPrefix(strings.TrimSuffix(getOwnRepo(ghURL, ght), "/"), "/")

	if strings.Contains(ownRepo, "/") != true {
		ownRepo, hasErrors = chooseRepo(ownRepo, ght)
		if hasErrors != nil {
			return 0, "", "", "", hasErrors
		}
		hasErrors = append(hasErrors, "Repo not provided in Source code URL, guessed to be https://www.github.com/"+ownRepo)
	}
	r := strings.Split(ownRepo, "/")

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
		//log.Printf("Error: https://github.com/%s/%s => %v Trying github api v3.  ", owner, repo, gh.Errors[0].Message)
		return ghv3api(r[0], r[1], ght)
	}
	return gh.Data.Repository.Stargazers.TotalCount, strings.Split(gh.Data.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate, "T")[0], gh.Data.Repository.LicenseInfo.SpdxID, gh.Data.Repository.PrimaryLanguage.Name, hasErrors
}

func ghv3api(u, r, ght string) (stars int, commitDate, license, language string, hasError []string) {
	ghURL := "https://api.github.com/repos/" + u + "/" + r + "?"

	req, err := http.NewRequest("GET", ghURL, nil)
	req.Header.Set("Authorization", "bearer "+ght)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		hasError = append(hasError, "StatusCode: "+strconv.Itoa(res.StatusCode)+". Source link is invalid: "+u+"/"+r)
		return 0, "", "", "", hasError
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("Body: %s\n", body)

	gh := ghv3{}
	err = json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal(err)
	}
	hasError = append(hasError, "Update https://github.com/"+u+"/"+r+" Source code link to: https://github.com/"+gh.FullName)
	if "/"+u+"/"+r != gh.FullName {
		//log.Printf("Retrying Github api v4 with %s", gh.FullName)
		GetGH(gh.FullName, ght, hasError)
	}
	return gh.Stars, strings.Split(gh.Updated, "T")[0], gh.License.SpdxID, gh.Language, hasError
}

// chooserepo guesses the missing repo. It doesn't do a very good job.
func chooseRepo(ur, ght string) (url string, haserr []string) {
	//fmt.Println(ur)
	var jsonStr = []byte(`{"query":"{repositoryOwner(login:\"` + ur + `\"){repositories(first: 100, orderBy: {direction: DESC, field: STARGAZERS}){edges{node{nameWithOwner}}}pinnedRepositories(first: 1, orderBy: {direction: DESC, field: STARGAZERS}){edges{node{nameWithOwner}}}}}"}`)
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
	}

	if len(gh.Data.RepositoryOwner.Repositories.Edges) == 1 {
		return gh.Data.RepositoryOwner.Repositories.Edges[0].Node.Name, nil
	}
	for _, r := range gh.Data.RepositoryOwner.PinnedRepositories.Edges {
		if r.Node.Name == ur {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, "serve") == true {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, "back") == true {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, "api") == true {
			return r.Node.Name, nil
		}
	}
	for _, r := range gh.Data.RepositoryOwner.Repositories.Edges {
		if r.Node.Name == ur {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, "serve") == true {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, "back") == true {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, "api") == true {
			return r.Node.Name, nil
		} else if strings.Contains(r.Node.Name, ur) == true {
			return r.Node.Name, nil
		}
	}
	// if len(gh.Data.RepositoryOwner.PinnedRepositories.Edges) >= 2 {
	// 	return gh.Data.RepositoryOwner.Repositories.Edges[0].Node.Name
	// }

	haserr = append(haserr, ur+" did not match or does not exist.")
	return "", haserr

}
