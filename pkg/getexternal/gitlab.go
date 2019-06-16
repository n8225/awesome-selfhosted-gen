package getexternal

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
)

// GetGLRepo retrieves star count and last activity.  This should probable be update to retrieve last commit and possible use GRAPHQL.
func GetGLRepo(url string) (int, string) {
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