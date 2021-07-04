package getexternal

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// GhEdge edges struct
type GhEdge struct {
	Node struct {
		NameWithOwner string `json:"nameWithOwner"`
		Name          string `json:"name"`
		Stargazers    struct {
			TotalCount int `json:"totalcount"`
		} `json:"stargazers"`
	} `json:"node"`
}

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
		User struct {
			PinnedItems struct {
				Edges []GhEdge `json:"edges"`
			} `json:"pinnedItems"`
			Repositories struct {
				Edges []GhEdge `json:"edges"`
			} `json:"repositories"`
		} `json:"user"`
		Organization struct {
			PinnedItems struct {
				Edges []GhEdge `json:"edges"`
			} `json:"pinnedItems"`
			Repositories struct {
				Edges []GhEdge `json:"edges"`
			} `json:"repositories"`
		} `json:"organization"`
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

//GetGH uses the github APIv4 (GRAPHQL) to retrieve star count, last commit, license, and language info. If the repository has moved or there is an error it will fall back to APIv3.
func GetGH(ghURL, ght string, hasErrors []string) (int, string, string, string, []string) {

	ownRepo := strings.TrimPrefix(strings.TrimSuffix(getOwnRepo(ghURL, ght), "/"), "/")

	if !strings.Contains(ownRepo, "/") {
		ownRepo, hasErrors = chooseRepo(ownRepo, ght)
		if hasErrors != nil {
			return 0, "", "", "", hasErrors
		}
		hasErrors = append(hasErrors, "Repo not provided in Source code URL, guessed to be https://www.github.com/"+ownRepo)
	}
	r := strings.Split(ownRepo, "/")
	//var jsonStr = []byte(`{"query":"{repository(owner:\"` + r[0] + `\",name:\"` + r[1] + `\"){stargazers{totalCount}licenseInfo{spdxId}primaryLanguage{name}defaultBranchRef{target{... on Commit{history(first: 1){edges{node{committedDate}}}}}}}}"}`)
	var jsonStr = []byte(`{"query":"{repository(owner:\"` + r[0] + `\",name:\"` + r[1] + `\"){stargazers{totalCount}defaultBranchRef{target{... on Commit{history(first: 1){edges{node{committedDate}}}}}}}}"}`)

	gh := &Gh{}
	err := json.Unmarshal(ghClientv4(ght, jsonStr), &gh)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	if gh.Errors != nil {
		log.Print(gh.Errors)
	}

	//Fallback to github v3 api on error
	if len(gh.Data.Repository.DefaultBranchRef.Target.History.Edges) == 0 {
		log.Info().Msgf("---going github api v3")
		return ghv3api(r[0], r[1], ght)
	}
	return gh.Data.Repository.Stargazers.TotalCount, strings.Split(gh.Data.Repository.DefaultBranchRef.Target.History.Edges[0].Node.CommittedDate, "T")[0], gh.Data.Repository.LicenseInfo.SpdxID, gh.Data.Repository.PrimaryLanguage.Name, hasErrors
}

func ghv3api(u, r, ght string) (stars int, commitDate, license, language string, hasError []string) {
	body, hasError := ghClientv3(ght, u, r)
	if hasError != nil {
		return 0, "", "", "", hasError
	}

	gh := ghv3{}
	err := json.Unmarshal(body, &gh)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	hasError = append(hasError, "Update https://github.com/"+u+"/"+r+" Source code link to: https://github.com/"+gh.FullName)
	if "/"+u+"/"+r != gh.FullName {
		GetGH(gh.FullName, ght, hasError)
	}
	return gh.Stars, strings.Split(gh.Updated, "T")[0], gh.License.SpdxID, gh.Language, hasError
}

// chooserepo guesses the missing repo. It doesn't do a very good job.
func chooseRepo(ur, ght string) (url string, haserr []string) {
	var jsonStr = []byte(`{"query": "{ user(login: \"` + ur + `\") { pinnedItems(first: 6, types: REPOSITORY) { edges { node { ... on Repository { nameWithOwner stargazers { totalCount } } } } } repositories(first: 30, orderBy: {field: STARGAZERS, direction: DESC}) { edges { node { nameWithOwner stargazers { totalCount } } } } } organization(login: \"` + ur + `\") { pinnedItems(first: 6) { edges { node { ... on Repository { nameWithOwner stargazers { totalCount } } } } } repositories(first: 30, orderBy: {field: STARGAZERS, direction: DESC}) { edges { node { nameWithOwner stargazers { totalCount } } } } } }"}`)

	gh := &Gh{}
	err := json.Unmarshal(ghClientv4(ght, jsonStr), &gh)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	/* 	if len(gh.Data.User.Repositories.Edges) == 1 {
		return gh.Data.RepositoryOwner.Repositories.Edges[0].Node.Name, nil
	} */
	if gh.Errors != nil {
		log.Print(gh.Errors)
	}
	switch 1 {
	case len(gh.Data.User.PinnedItems.Edges):
		return gh.Data.User.PinnedItems.Edges[0].Node.NameWithOwner, nil
	case len(gh.Data.Organization.PinnedItems.Edges):
		return gh.Data.Organization.PinnedItems.Edges[0].Node.NameWithOwner, nil
	}

	log.Info().Msg("chooserepo: " + ur)
	var res string
	if len(gh.Data.User.Repositories.Edges) > 0 {
		if len(gh.Data.User.PinnedItems.Edges) > 0 {
			res = ghRepoPicker(ur, gh.Data.User.PinnedItems.Edges)
			if res == "" {
				res = gh.Data.User.PinnedItems.Edges[0].Node.NameWithOwner
			}
		} else {
			res = ghRepoPicker(ur, gh.Data.User.Repositories.Edges)
			if res == "" {
				res = gh.Data.User.Repositories.Edges[0].Node.NameWithOwner
			}
		}
		return res, nil
	}
	if len(gh.Data.Organization.Repositories.Edges) > 0 {
		if len(gh.Data.Organization.PinnedItems.Edges) > 0 {
			res = ghRepoPicker(ur, gh.Data.Organization.PinnedItems.Edges)
			if res == "" {
				res = gh.Data.Organization.PinnedItems.Edges[0].Node.NameWithOwner
			}
		} else {
			res = ghRepoPicker(ur, gh.Data.Organization.Repositories.Edges)
			if res == "" {
				res = gh.Data.Organization.Repositories.Edges[0].Node.NameWithOwner
			}
		}
		return res, nil
	}
	// if len(gh.Data.RepositoryOwner.PinnedRepositories.Edges) >= 2 {
	// 	return gh.Data.RepositoryOwner.Repositories.Edges[0].Node.Name
	// }
	haserr = append(haserr, ur+" did not match or does not exist.")
	return "", haserr
}

func ghRepoPicker(ur string, repos []GhEdge) string {
	list := [5]string{"server", "serve", "back", "stack", "api"}
	for _, r := range repos {
		n := r.Node.NameWithOwner
		if strings.EqualFold(r.Node.Name, ur) {
			return n
		}
		for _, l := range list {
			if strings.Contains(strings.ToLower(n), strings.ToLower(l)) {
				return n
			}
		}
	}
	return ""
}

func ghClientv4(ght string, jsonStr []byte) []byte {
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	req.Header.Set("Authorization", "bearer "+ght)
	req.Header.Set("Accept", "application/vnd.github.quicksilver-preview+json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	defer res.Body.Close()
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	return body
}

func ghClientv3(ght, u, r string) ([]byte, []string) {
	url := "https://api.github.com/repos/" + u + "/" + r + "?"
	var hasError []string
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	req.Header.Set("Authorization", "bearer "+ght)
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	if res.StatusCode != 200 {
		hasError = append(hasError, "StatusCode: "+strconv.Itoa(res.StatusCode)+". Source link is invalid: "+u+"/"+r)
		return nil, hasError
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	return body, hasError
}

//getOwnRepo parses the owner and repo from a github url
func getOwnRepo(ghURL, ght string) string {
	u, err := url.Parse(ghURL)
	if err != nil {
		log.Error().Stack().Err(err)
	}
	return u.Path
}
