package backend

import (
	"bytes"
	"com/abnamro/solo/nexus-repository-cli/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"net/http/httputil"
)

func addReposToGroupRepo(nexusBaseUrl string, user model.User, repoId string, groupRepoIds []string, verbose bool) {
	// retrieve groupRepo
	url := fmt.Sprintf("%v/service/local/repo_groups/%v", nexusBaseUrl, repoId)
	var targetRepo GroupRepository
	jsonStruct, err := UrlToJsonStruct(url, targetRepo)
	if err != nil {
		log.Fatal(err)
		return
	}
	// parse json to struct
	targetRepo = jsonStruct.(GroupRepository)

	var reposToAdd []string
	for _, repoIdToAdd := range groupRepoIds {
		// check if repoId exists in repo-group-member
		repoExists := false
		for _, repo := range targetRepo.Data.Repositories {
			if repo.ID == repoIdToAdd {
				repoExists = true
			}
		}
		if !repoExists {
			repoToAdd := GroupRepositoryRepo{ID: repoIdToAdd}
			targetRepo.Data.Repositories = append(targetRepo.Data.Repositories, repoToAdd)
			reposToAdd = append(reposToAdd, repoIdToAdd)
		}
	}

	fmt.Printf("[FOUND %v Repositories to add to %v]\n", len(reposToAdd), targetRepo.Data.ID)

	if len(reposToAdd) <= 0 {
		fmt.Printf("[No repositories to add, will not update %v]\n\n", targetRepo.Data.ID)
	} else {
		fmt.Printf("[Will add \n%v]\n\n", reposToAdd)
		updateGroupRepo(url, user, targetRepo, verbose)
	}
}

func AddRepoToGroupRepo(nexusBaseUrl string, user model.User, repoId string, separator string,  groupRepoId string, verbose bool) {
	// retrieve groupRepo
	url := fmt.Sprintf("%v/service/local/repo_groups/%v", nexusBaseUrl, groupRepoId)
	var targetRepo GroupRepository
	jsonStruct, err := UrlToJsonStruct(url, targetRepo)
	if err != nil {
		log.Fatal(err)
		return
	}
	// parse json to struct
	targetRepo = jsonStruct.(GroupRepository)

	var reposToAdd []string
	if separator != "\n" && strings.Contains(repoId, separator) {
		reposList := strings.Split(repoId, separator)
		for _,repo := range reposList {
			reposToAdd= append(reposToAdd, repo)
		}
	} else if separator != "\n" {
		fmt.Printf("Non-empty separator found but separator (%s) not found in repoId (%s)", separator, repoId)
		return
	} else {
		reposToAdd= append(reposToAdd, repoId)
	}

	// check if repoId exists in repo-group-member
	for _,repoToAdd := range reposToAdd {
		repoExists := false
		for _, repo := range targetRepo.Data.Repositories {
			if repo.ID == repoToAdd {
				repoExists = true
			}
		}
		if repoExists {
			fmt.Printf("Repo %v is already part of group repo %v\n", repoToAdd, groupRepoId)
			return
		} else {
			repoToAddStruct := GroupRepositoryRepo{ID: repoToAdd}
			targetRepo.Data.Repositories = append(targetRepo.Data.Repositories, repoToAddStruct)
		}
	}

	updateGroupRepo(url, user, targetRepo, verbose)
}

func updateGroupRepo(url string, user model.User, targetRepo GroupRepository, verbose bool) {

	body, err := json.Marshal(targetRepo)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user.Username, user.Password)

	if verbose {
		logging, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(logging))
	}

	handleUpdateGroupRepoResponse(req, verbose)
}

func handleUpdateGroupRepoResponse(req *http.Request, verbose bool) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Println("Request Headers:", req.Header)
		fmt.Println("Response Headers:", resp.Header)
		responseBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response Body:", string(responseBody))
	}

	switch resp.Status {
	case "200 OK":
		fmt.Printf("Success: linked repositories\n")
	default:
		panic(fmt.Sprintf("ERROR: call status=%v\n", resp.Status))
	}

}
