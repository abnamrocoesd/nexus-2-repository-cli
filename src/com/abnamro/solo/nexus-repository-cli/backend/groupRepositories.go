package backend

import (
	"com/abnamro/solo/nexus-repository-cli/model"
	"encoding/json"
	"log"
)

type GroupRepositories struct {
	Data []struct {
		ResourceURI        string `json:"resourceURI,omitempty"`
		ContentResourceURI string `json:"contentResourceURI,omitempty"`
		ID                 string `json:"id,omitempty"`
		Name               string `json:"name,omitempty"`
		Format             string `json:"format,omitempty"`
		Exposed            bool   `json:"exposed"`
		UserManaged        bool   `json:"userManaged,omitempty"`
	} `json:"data"`
}

func (groupRepositories GroupRepositories) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &groupRepositories)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}
	return groupRepositories, nil
}

func (groupRepositories GroupRepositories) ToInternal() (internalGroupRepositories []model.GroupRepository) {

	for _, groupRepo := range groupRepositories.Data {
		newRepo := model.GroupRepository{
			ID:     groupRepo.ID,
			Name:   groupRepo.Name,
			Format: groupRepo.Format,
		}
		internalGroupRepositories = append(internalGroupRepositories, newRepo)
	}
	return
}
