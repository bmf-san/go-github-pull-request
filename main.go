package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bmf-san/go-github-pull-request/client"
	"github.com/gocarina/gocsv"
)

var token string
var path string
var assignee string
var org string
var owner string

func init() {
	flag.StringVar(&token, "token", "", "GITHUB_API_TOKEN")
	flag.StringVar(&path, "path", "", "PATH_TO_CSV_FILE")
	flag.StringVar(&assignee, "assignee", "", "ASSIGNEE")
	flag.StringVar(&org, "org", "", "ORG")
	flag.StringVar(&owner, "owner", "", "OWNER")
}

const (
	// max is 100
	perPage = 100
	state   = "all"
	// Since there is no option to get all pages, specify 1000 pages for the time being
	page = 1000
)

func main() {
	log.Print("[Start]")

	flag.Parse()

	c := client.NewClient(token)

	var repos []*client.Repository
	for i := 1; i < page; i++ {
		fmt.Printf("GetOrgRepos: page %v\n", i)
		r, err := c.Repos.GetOrgRepos(org, client.GetOrgRepoParams{
			PerPage: perPage,
			Page:    i,
		})
		if err != nil {
			log.Fatal(err)
		}
		if len(r) == 0 {
			break
		}
		repos = append(repos, r...)
	}

	data := map[string][]*client.PullRequest{}
	for _, r := range repos {
		for i := 1; i < page; i++ {
			fmt.Printf("GetPullsByAssignee: owner %v repo %v page %v\n", owner, r.Name, i)
			p, err := c.Pulls.GetPullsByAssignee(assignee, owner, r.Name, client.GetPullParams{
				PerPage: perPage,
				State:   state,
				Page:    i,
			})
			if err != nil {
				log.Fatal(err)
			}
			if len(p) == 0 {
				break
			}
			data[r.Name] = append(data[r.Name], p...)
		}
	}

	for i, v := range data {
		genCSV(fmt.Sprintf("%v.csv", i), v)
	}

	log.Print("[Finish]")
}

func genCSV(fileName string, in interface{}) {
	f, err := os.OpenFile(path+fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	gocsv.MarshalFile(in, f)
}
