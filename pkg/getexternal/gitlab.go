package getexternal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//GetGLRepo retrieves star count and last activity.  //TODO This should probable be updated to retrieve last commit and possible use GRAPHQL.
func GetGLRepo(ur string) (int, string) {
	if strings.Contains(ur, "/") != true {
		return 0, ""
	}
	r := strings.Split(ur, "/")

	url := "https://gitlab.com/api/v4/projects/" + r[0] + "%2f" + r[1]
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	type gl struct {
		Stars   int    `json:"star_count"`
		Updated string `json:"last_activity_at"`
		//Node_id int `json:"id"`
	}
	thisgl := gl{}
	err = json.Unmarshal(body, &thisgl)
	if err != nil {
		log.Fatal(err)
	}
	return thisgl.Stars, strings.Split(thisgl.Updated, "T")[0]
}
