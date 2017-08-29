package backend

import (
	"com/abnamro/solo/nexus-repository-cli/model"
	"net/http"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
)

func DeleteRepository (baseUrl string, user model.User, targetRepo string, repoType string, verbose bool) (err error) {
	url := fmt.Sprintf("%s/service/local/repositories/%s", baseUrl, targetRepo)
	if repoType == "group" {
		url = fmt.Sprintf("%s/service/local/repo_groups/%s", baseUrl, targetRepo)
	}
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(user.Username, user.Password)

	if verbose {
		logging, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(logging))
	}

	handleDeleteResponse(req, verbose)
	return
}

func handleDeleteResponse(req *http.Request, verbose bool) {
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
	case "204 No Content":
		fmt.Printf("Success: deleted\n")
	case "404 Not Found":
		fmt.Printf("Warning: repository not found\n")
	default:
		panic(fmt.Sprintf("ERROR: call status=%v\n", resp.Status))
	}
}