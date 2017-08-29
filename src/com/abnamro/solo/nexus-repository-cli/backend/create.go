package backend

import (
	"bytes"
	"com/abnamro/solo/nexus-repository-cli/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

type CreateRepository struct {
	Data CreateRepositoryData `json:"data"`
}

type CreateRepositoryData struct {
	RepoType              string `json:"repoType"`
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	WritePolicy           string `json:"writePolicy"`
	Browseable            bool   `json:"browseable"`
	Indexable             bool   `json:"indexable"`
	Exposed               bool   `json:"exposed"`
	NotFoundCacheTTL      int    `json:"notFoundCacheTTL"`
	RepoPolicy            string `json:"repoPolicy"`
	Provider              string `json:"provider"`
	ProviderRole          string `json:"providerRole"`
	DownloadRemoteIndexes bool   `json:"downloadRemoteIndexes"`
	ChecksumPolicy        string `json:"checksumPolicy"`
}

type CreateGroupRepository struct {
	Data CreateGroupRepositoryData `json:"data"`
}

type CreateProxyRepositoryData struct {
	ContentResourceURI     string              `json:"contentResourceURI"`
	ID                     string              `json:"id"`
	Name                   string              `json:"name"`
	Provider               string              `json:"provider"`
	ProviderRole           string              `json:"providerRole"`
	Format                 string              `json:"format"`
	RepoType               string              `json:"repoType"`
	Exposed                bool                `json:"exposed"`
	WritePolicy            string              `json:"writePolicy"`
	Browseable             bool                `json:"browseable"`
	Indexable              bool                `json:"indexable"`
	NotFoundCacheTTL       int                 `json:"notFoundCacheTTL"`
	RepoPolicy             string              `json:"repoPolicy"`
	ChecksumPolicy         string              `json:"checksumPolicy"`
	DownloadRemoteIndexes  bool                `json:"downloadRemoteIndexes"`
	DefaultLocalStorageURL string              `json:"defaultLocalStorageUrl"`
	RemoteStorage          createRemoteStorage `json:"remoteStorage"`
	FileTypeValidation     bool                `json:"fileTypeValidation"`
	ArtifactMaxAge         int                 `json:"artifactMaxAge"`
	MetadataMaxAge         int                 `json:"metadataMaxAge"`
	ItemMaxAge             int                 `json:"itemMaxAge"`
	AutoBlockActive        bool                `json:"autoBlockActive"`
}

type createRemoteStorage struct {
	RemoteStorageURL string `json:"remoteStorageUrl"`
}

type CreateProxyRepository struct {
	Data CreateProxyRepositoryData `json:"data"`
}

type CreateGroupRepositoryData struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Format       string   `json:"format"`
	Exposed      bool     `json:"exposed"`
	Provider     string   `json:"provider"`
	Repositories []string `json:"repositories"`
}

func CreateGroup(baseUrl string, user model.User, targetRepo string, provider string, verbose bool) (err error) {

	url := fmt.Sprintf("%s/service/local/repo_groups", baseUrl)
	format := strings.ToLower(provider)

	repository := &CreateGroupRepository{
		Data: CreateGroupRepositoryData{
			ID:           targetRepo,
			Name:         targetRepo,
			Repositories: []string{},
			Format:       format,
			Provider:     provider,
			Exposed:      true,
		},
	}

	body, err := json.Marshal(repository)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user.Username, user.Password)

	if verbose {
		logging, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(logging))
	}

	handleCreateResponse(req, verbose)

	return
}

func createRepo(repoType string, targetRepo string, repoPolicy string, provider string) (repository CreateProxyRepository, err error) {
	repo := CreateProxyRepository{
		Data: CreateProxyRepositoryData{
			RepoType:              repoType,
			ID:                    targetRepo,
			Name:                  targetRepo,
			WritePolicy:           "ALLOW_WRITE",
			Browseable:            true,
			Indexable:             true,
			Exposed:               true,
			NotFoundCacheTTL:      1440,
			RepoPolicy:            repoPolicy,
			Provider:              provider,
			ProviderRole:          "org.sonatype.nexus.proxy.repository.Repository",
			DownloadRemoteIndexes: false,
			ChecksumPolicy:        "IGNORE",
		},
	}
	return repo, nil
}

func createProxyRepo(repoType string, targetRepo string, repoPolicy string, provider string, remote string) (repository CreateProxyRepository, err error) {

	localStorageUrl := fmt.Sprintf("file:/data/sonatype-work/nexus/storage/%s", targetRepo)

	repo := CreateProxyRepository{
		Data: CreateProxyRepositoryData{
			RepoType:               repoType,
			ID:                     targetRepo,
			Name:                   targetRepo,
			WritePolicy:            "ALLOW_WRITE",
			Browseable:             true,
			Indexable:              true,
			Exposed:                true,
			NotFoundCacheTTL:       1440,
			RepoPolicy:             repoPolicy,
			Provider:               provider,
			ProviderRole:           "org.sonatype.nexus.proxy.repository.Repository",
			ChecksumPolicy:         "IGNORE",
			DownloadRemoteIndexes:  true,
			DefaultLocalStorageURL: localStorageUrl,
			RemoteStorage: createRemoteStorage{
				RemoteStorageURL: remote,
			},
			FileTypeValidation: true,
			ArtifactMaxAge:     -1,
			MetadataMaxAge:     1440,
			ItemMaxAge:         1440,
			AutoBlockActive:    true,
		},
	}
	return repo, nil
}

func CreateProxy(baseUrl string, user model.User, targetRepo string, provider string, repoType string, repoPolicy string, remote string, verbose bool) (err error) {
	url := fmt.Sprintf("%s/service/local/repositories", baseUrl)

	localStorageUrl := fmt.Sprintf("file:/data/sonatype-work/nexus/storage/%s", targetRepo)
	repository := CreateProxyRepository{
		Data: CreateProxyRepositoryData{
			RepoType:               repoType,
			ID:                     targetRepo,
			Name:                   targetRepo,
			WritePolicy:            "ALLOW_WRITE",
			Browseable:             true,
			Indexable:              true,
			Exposed:                true,
			NotFoundCacheTTL:       1440,
			RepoPolicy:             repoPolicy,
			Provider:               provider,
			ProviderRole:           "org.sonatype.nexus.proxy.repository.Repository",
			ChecksumPolicy:         "IGNORE",
			DownloadRemoteIndexes:  true,
			DefaultLocalStorageURL: localStorageUrl,
			RemoteStorage: createRemoteStorage{
				RemoteStorageURL: remote,
			},
			FileTypeValidation: true,
			ArtifactMaxAge:     -1,
			MetadataMaxAge:     1440,
			ItemMaxAge:         1440,
			AutoBlockActive:    true,
		},
	}

	body, err := json.Marshal(repository)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user.Username, user.Password)

	if verbose {
		logging, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(logging))
	}

	handleCreateResponse(req, verbose)

	return
}

func Create(baseUrl string, user model.User, targetRepo string, provider string, repoType string, repoPolicy string, verbose bool) (err error) {
	url := fmt.Sprintf("%s/service/local/repositories", baseUrl)
	repository := CreateRepository{
		Data: CreateRepositoryData{
			RepoType:              repoType,
			ID:                    targetRepo,
			Name:                  targetRepo,
			WritePolicy:           "ALLOW_WRITE",
			Browseable:            true,
			Indexable:             true,
			Exposed:               true,
			NotFoundCacheTTL:      1440,
			RepoPolicy:            repoPolicy,
			Provider:              provider,
			ProviderRole:          "org.sonatype.nexus.proxy.repository.Repository",
			DownloadRemoteIndexes: false,
			ChecksumPolicy:        "IGNORE",
		},
	}

	body, err := json.Marshal(repository)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user.Username, user.Password)

	if verbose {
		logging, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(logging))
	}

	handleCreateResponse(req, verbose)

	return
}

func handleCreateResponse(req *http.Request, verbose bool) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if verbose {
		fmt.Println("Request Headers:", req.Header)
		fmt.Println("Response Headers:", resp.Header)
		fmt.Println("Response Status:", resp.Status)
		responseBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response Body:", string(responseBody))
	}

	switch resp.Status {
	case "200 OK":
		fmt.Printf("Success: OK\n")
	case "201 Created":
		fmt.Printf("Success: new resource created\n")
	default:
		panic(fmt.Sprintf("ERROR: call status=%v\n", resp.Status))
	}
}
