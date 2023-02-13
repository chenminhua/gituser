package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"

	"github.com/brettski/go-termtables"
)

// Repo alias
type Repo map[string]interface{}

// Repos alias
type Repos []Repo
// Len for sort
func (repos Repos) Len() int {
  return len(repos)
}
// Swap for sort
func (repos Repos) Swap(i,j int){
  repos[i], repos[j] = repos[j], repos[i]
}
// Less for sort
func (repos Repos) Less(i,j int) bool{
  return repos[i]["stargazers_count"].(float64) > repos[j]["stargazers_count"].(float64)
}


//["name","email","location","follower","following", "created_at"]
func getRepos(url string, ch chan []Repo){
  resp,_ := http.Get(url)
  body,_ := ioutil.ReadAll(resp.Body)
  var res []Repo
  json.Unmarshal(body, &res)

  ch <- res
}

func main(){
  c := make(chan []Repo)
  if len(os.Args) <= 1{
    fmt.Println("usage: gituser [username]")
    fmt.Println("for example: gituser chenminhua")
    os.Exit(1)
  }
  username := os.Args[1]
  fmt.Printf("retrieving %s's info ...\n", username)
  resp, _ := http.Get(fmt.Sprintf("https://api.github.com/users/%s", username))
  body,_ := ioutil.ReadAll(resp.Body)
  var res map[string]interface{}
  json.Unmarshal(body, &res)

  ut := termtables.CreateTable()
  ut.AddHeaders("name","email","location","follower","following", "repos", "created_at")
  ut.AddRow(res["login"], res["email"], res["location"], res["followers"], res["following"], res["public_repos"],res["created_at"])
  fmt.Println(ut.Render())

  totalWorkers := int(math.Ceil(res["public_repos"].(float64) / 30))
  //pages := make([]string, totalWorkers)
  repos := []Repo{}
  for i :=1; i<totalWorkers+1;i++{
    go getRepos(fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d", username, i), c)
    repos = append(repos, <-c ...)
  }
  sort.Sort(Repos(repos))
  totalStars := 0
  rt := termtables.CreateTable()
  rt.AddHeaders("name","star","fork","html","language","description")
  for i:=0; i<len(repos);i++{
    if repos[i]["fork"].(bool){
      continue
    }
    if repos[i]["stargazers_count"].(float64) == 0{
      break
    }
    totalStars += int(repos[i]["stargazers_count"].(float64))
    rt.AddRow(repos[i]["name"], repos[i]["stargazers_count"],repos[i]["forks"],repos[i]["html_url"], repos[i]["language"],repos[i]["description"])
  }
  fmt.Println(rt.Render())
  fmt.Println(fmt.Sprintf("%s has get %d stars", username, totalStars))
}
