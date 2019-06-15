package util

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
)

func getGLRepo(url string) (int, string) {
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
		Created string `json:"created_at"`
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