package backend

import "encoding/json"

/* RepoPath */
type RepoPath struct {
	Data []struct {
		ResourceURI  string `json:"resourceURI"`
		RelativePath string `json:"relativePath"`
		Text         string `json:"text"`
		Leaf         bool   `json:"leaf"`
		LastModified string `json:"lastModified"`
		SizeOnDisk   int    `json:"sizeOnDisk"`
	} `json:"data"`
}

func (r RepoPath) Unmarshal(data []byte) (JsonStruct, error) {
	jsonErr := json.Unmarshal(data, &r)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return r, nil
}
