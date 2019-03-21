package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	User string `json:"user`
	Repo string `json:"repo`
}

// see: https://github.com/google/go-github/blob/master/github/pulls.go
type PullRequest struct {
	ID        *int64     `json:"id,omitempty"`
	Number    *int       `json:"number,omitempty"`
	State     *string    `json:"state,omitempty"`
	Title     *string    `json:"title,omitempty"`
	Body      *string    `json:"body,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	HTMLURL   *string    `json:"html_url,omitempty"`
	Assignee  *User      `csv:"-"`
}

type User struct {
	Login *string `json:"login,omitempty"`
}

func main() {
	LoadEnv()

	row, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := json.Unmarshal(row, &config); err != nil {
		log.Fatal(err)
	}

	// TODO: Refactoring
	for _, c := range config.Repositories {
		// for the time being, max page is 10
		for index := 1; index < 11; index++ {
			// see: https://developer.github.com/enterprise/2.4/v3/pulls/#list-pull-requests
			resp, err := http.Get("https://api.github.com/repos/" + c.User + "/" + c.Repo + "/pulls?access_token=" + os.Getenv("GITHUB_API_TOKEN") + "&per_page=100&state=all&page=" + strconv.Itoa(index))

			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)

			var pullrequest []PullRequest
			if err := json.Unmarshal(body, &pullrequest); err != nil {
				log.Fatal(err)
			}

			// List pull requests API doesn't support Assignee query parameter
			// see: https://developer.github.com/enterprise/2.4/v3/pulls/#list-pull-requests
			result := []PullRequest{}
			for _, p := range pullrequest {
				if p.Assignee != nil {
					if *p.Assignee.Login == os.Getenv("ASSIGNEE_USER") {
						result = append(result, p)
					}
				}
			}

			if len(result) > 1 {
				file, err := os.OpenFile(os.Getenv("PATH_TO_CSV_FILE")+c.Repo+".csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

				if err != nil {
					log.Fatal(err)
				}

				defer file.Close()

				gocsv.MarshalFile(result, file)
			}
		}
	}
}
