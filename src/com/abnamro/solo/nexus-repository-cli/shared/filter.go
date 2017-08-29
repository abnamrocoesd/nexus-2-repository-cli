package shared

import (
	"com/abnamro/solo/nexus-repository-cli/backend"
	"com/abnamro/solo/nexus-repository-cli/model"
	"fmt"
	"regexp"
	"sort"
)

func FilterGroupRepositories(groupRepositories []model.GroupRepository, format string, nameFilter regexp.Regexp, verbose bool) (filteredGroupRepositories []model.GroupRepository, err error) {
	sort.Slice(groupRepositories[:], func(i, j int) bool {
		return groupRepositories[i].ID < groupRepositories[j].ID
	})
	for _, repo := range groupRepositories {
		if verbose {
			fmt.Printf("FilterGroupRepositories: %s", repo.Name)
		}
		if (format == "" || repo.Format == format) && nameFilter.MatchString(repo.Name) {
			if verbose {
				fmt.Printf("-> match")
			}
			filteredGroupRepositories = append(filteredGroupRepositories, repo)
		}
		if verbose {
			fmt.Printf("\n")
		}
	}
	return
}

func FilterLinkedRepositories(baseUrl string, groupRepositories []model.GroupRepository, nameFilter string, include bool, verbose bool) (filteredGroupRepositories []model.GroupRepository, err error) {
	for _, repo := range groupRepositories {
		found := false
		var linkedRepositories []model.Repository
		if verbose {
			fmt.Printf("FilterLinkedRepositories: %s", repo.Name)
		}
		linkedRepositories, err = backend.RetrieveLinkedRepositories(baseUrl, repo)
		if err != nil {
			return
		}

		for _, linkedRepo := range linkedRepositories {
			if verbose {
				fmt.Printf("%s", linkedRepo.Name)
			}
			if linkedRepo.Name == nameFilter {
				if verbose {
					fmt.Printf(" -> match")
				}
				found = true
			}
			if verbose {
				fmt.Printf("\n")
			}
		}
		if (found && include) || (!found && !include) {
			filteredGroupRepositories = append(filteredGroupRepositories, repo)
		}
	}
	return
}
