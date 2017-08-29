package backend

import (
	"com/abnamro/solo/nexus-repository-cli/model"
	"fmt"
	"log"
	"regexp"
	"sort"
)

func GetRepositories(nexusBaseUrl string) (repos Repositories, err error) {
	url := fmt.Sprintf("%v/service/local/repositories", nexusBaseUrl)

	jsonStruct, err := UrlToJsonStruct(url, repos)
	if err != nil {
		log.Fatal(err)
		return
	}

	repos = jsonStruct.(Repositories)
	return
}

func GetGroupRepositories(nexusBaseUrl string) (repositories []model.GroupRepository, err error) {
	url := fmt.Sprintf("%v/service/local/repo_groups", nexusBaseUrl)

	groupRepositories := GroupRepositories{}
	jsonStruct, err := UrlToJsonStruct(url, groupRepositories)

	if err != nil {
		log.Fatal(err)
		return
	}

	groupRepositories = jsonStruct.(GroupRepositories)
	repositories = groupRepositories.ToInternal()

	return
}

func RetrieveLinkedRepositories(baseUrl string, repository model.GroupRepository) (repositories []model.Repository, err error) {
	url := fmt.Sprintf("%s/service/local/repo_groups/%s", baseUrl, repository.ID)
	linkedRepositories := LinkedRepositories{}
	jsonStruct, err := UrlToJsonStruct(url, linkedRepositories)

	if err != nil {
		log.Fatal(err)
		return
	}

	linkedRepositories = jsonStruct.(LinkedRepositories)
	repositories = linkedRepositories.ToInternal()

	return
}

func FilterUserRepositories(baseUrl string, repos Repositories, provider string, repoType string, nameFilter regexp.Regexp, verbose bool) (filteredRepositories []model.Repository, err error) {

	// sort by ID
	sort.Slice(repos.Data[:], func(i, j int) bool {
		return repos.Data[i].ID < repos.Data[j].ID
	})

	for _, repo := range repos.Data {
		if verbose {
			fmt.Printf("looking at repo.Provider: %v, repo.repoType: %v, repo.Name: %v\n", repo.Provider, repo.RepoType, repo.Name)
			fmt.Printf("match: %v\n", nameFilter.MatchString(repo.Name))

		}
		if (provider == "" || repo.Provider == provider) && repo.RepoType == repoType && nameFilter.MatchString(repo.Name) {
			repository := processRepository(baseUrl, repo.ID, repo.Name, verbose)
			filteredRepositories = append(filteredRepositories, repository)
		}

	}
	if verbose {
		fmt.Println("===========================================")
		fmt.Printf("== Found %v filteredRepositories.", len(filteredRepositories))
		fmt.Println("===========================================")
		fmt.Println("===========================================")
	}
	return
}

func processRepository(baseUrl string, id string, name string, verbose bool) model.Repository {
	if verbose {
		fmt.Printf("%v,", id)
	}
	repository := model.Repository{}
	repository.ID = id
	repository.Name = name
	findGroupIds(baseUrl, &repository, 0, "")
	return repository
}

func findGroupIds(baseUrl string, repo *model.Repository, level int, groupIdRoot string) {
	url := fmt.Sprintf("%s/service/local/repositories/%s/content", baseUrl, repo.ID)
	if level > 0 {
		url += fmt.Sprintf("/%s", groupIdRoot)
	}
	url = fmt.Sprintf("%s/?isLocal", url)
	var repoPath RepoPath
	jsonStruct, err := UrlToJsonStruct(url, repoPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	if jsonStruct != nil {
		repoPath = jsonStruct.(RepoPath)
		for _, node := range repoPath.Data {
			if node.Leaf == false && node.SizeOnDisk == -1 {
				if level > 0 {
					groupId := fmt.Sprintf("%s.%s", groupIdRoot, node.Text)
					repo.GroupIds = append(repo.GroupIds, groupId)
				} else {
					findGroupIds(baseUrl, repo, 1, node.Text)
				}
			}
		}
	}
}
