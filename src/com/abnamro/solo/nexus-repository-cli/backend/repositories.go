package backend

import (
	"encoding/json"
	"log"
)

type Repositories struct {
	Data []struct {
		ResourceURI              string `json:"resourceURI"`
		ContentResourceURI       string `json:"contentResourceURI,omitempty"`
		ID                       string `json:"id"`
		Name                     string `json:"name"`
		RepoType                 string `json:"repoType"`
		RepoPolicy               string `json:"repoPolicy,omitempty"`
		Provider                 string `json:"provider"`
		ProviderRole             string `json:"providerRole"`
		Format                   string `json:"format"`
		UserManaged              bool   `json:"userManaged"`
		Exposed                  bool   `json:"exposed"`
		EffectiveLocalStorageURL string `json:"effectiveLocalStorageUrl"`
		RemoteURI                string `json:"remoteUri,omitempty"`
	} `json:"data"`
}

func (r Repositories) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}
	return r, nil
}
