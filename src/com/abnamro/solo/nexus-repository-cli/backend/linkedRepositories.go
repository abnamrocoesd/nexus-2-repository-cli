package backend

import (
	"com/abnamro/solo/nexus-repository-cli/model"
	"encoding/json"
	"log"
)

type LinkedRepositories struct {
	Data struct {
		ContentResourceURI string `json:"contentResourceURI"`
		ID                 string `json:"id"`
		Name               string `json:"name"`
		Provider           string `json:"provider"`
		Format             string `json:"format"`
		RepoType           string `json:"repoType"`
		Exposed            bool   `json:"exposed"`
		Repositories       []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceURI"`
		} `json:"repositories"`
	} `json:"data"`
}

func (linkedRepositories LinkedRepositories) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &linkedRepositories)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}
	return linkedRepositories, nil
}

func (linkedRepositories LinkedRepositories) ToInternal() (internalRepositories []model.Repository) {
	for _, repo := range linkedRepositories.Data.Repositories {
		repository := model.Repository{
			ID:   repo.ID,
			Name: repo.Name,
		}
		internalRepositories = append(internalRepositories, repository)
	}
	return
}
