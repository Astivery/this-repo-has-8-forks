package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	githubToken = "github_token"
	repoID      = "repo_id"
  githubName = "your_github_name"
)

var (
	client = &http.Client{}
	forks  int
	name   string
	test   = false
)

type githubReponse struct {
	Forks int    `json:"forks"`
	Name  string `json:"name"`
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getForks() {
	req, err := http.NewRequest("GET", "https://api.github.com/repositories/"+repoID, nil)
	check(err)
	req.Header.Add("Authorization", "Token "+githubToken)
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	var respBody githubReponse
	err = json.Unmarshal(body, &respBody)
	check(err)
	if !test {
		name = respBody.Name
		test = true
	}
	if forks != respBody.Forks {
		forks = respBody.Forks
		updateRepo()
	}
}

func updateRepo() {
	nameBody := fmt.Sprintf("{\"name\": \"this-repo-has-%v-forks\"}", forks)
	req, err := http.NewRequest("PATCH", "https://api.github.com/repos/"+githubName+"/"+name, strings.NewReader(nameBody))
	check(err)
	req.Header.Add("Authorization", "Token "+githubToken)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	check(err)
	resp.Body.Close()
	name = fmt.Sprintf("this-repo-has-%v-forks", forks)
}

func main() {
	fmt.Println("Bot started!")
	for {
		getForks()
		fmt.Println("Forks:", forks)
		time.Sleep(5 * time.Second)
	}
}
