package backend

import (
	"encoding/json"
	"log"
)

type GroupRepositoryRepo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ResourceURI string `json:"resourceURI"`
}

type GroupRepository struct {
	Data struct {
		ContentResourceURI string                `json:"contentResourceURI"`
		ID                 string                `json:"id"`
		Name               string                `json:"name"`
		Provider           string                `json:"provider"`
		Format             string                `json:"format"`
		RepoType           string                `json:"repoType"`
		Exposed            bool                  `json:"exposed"`
		Repositories       []GroupRepositoryRepo `json:"repositories"`
	} `json:"data"`
}

func (g GroupRepository) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &g)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}
	return g, nil
}
