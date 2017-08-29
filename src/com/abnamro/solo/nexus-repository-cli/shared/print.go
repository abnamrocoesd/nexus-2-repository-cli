package shared

import (
	"com/abnamro/solo/nexus-repository-cli/model"
	"fmt"
	"strings"
)

func PrintRepositories(repositories []model.Repository, seperator string) {
	var ids []string
	for _, repository := range repositories {
		ids = append(ids, repository.ID)
	}
	fmt.Println(strings.Join(ids, seperator))
}

func PrintGroupRepositories(groupRepositories []model.GroupRepository, seperator string) {
	var ids []string
	for _, repository := range groupRepositories {
		ids = append(ids, repository.ID)
	}
	fmt.Println(strings.Join(ids, seperator))
}
