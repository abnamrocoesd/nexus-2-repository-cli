package main

import (
	"com/abnamro/solo/nexus-repository-cli/backend"
	"com/abnamro/solo/nexus-repository-cli/model"
	"com/abnamro/solo/nexus-repository-cli/shared"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
)

func main() {
	nexusUrl := flag.String("nexusUrl", "http://localhost:8081/nexus", "Nexus host name (full base url), e.g. http://localhost:8081/nexus")

	action := flag.String("action", "find", `
	- find: find a repository. Optional arguments: 'type', 'format', 'filter', 'separator', 'linkedRepositoryFilter'
	- add: add repository/repositories to a groupRepository. Required: target, sourceRepo. (optional: separator, if more than one repo is listed in sourceRepo)
	- create: create new repository. Required arguments: username, password, type, provider, target (optional: remote, if repoType is 'proxy')
	- delete: delete a repository. Required arguments: username, password, type, target
	`)
	regexFilter := flag.String("filter", `(.*)`, "Optional: regex filter for the names of repositories")
	repoType := flag.String("type", "hosted", "Optional. Values: hosted, proxy, virtual, group")
	verbose := flag.Bool("verbose", false, "Set the flag if you want verbose output")
	provider := flag.String("provider", "", "The repository provider, such as maven2, npm, nuget, etc.")
	format := flag.String("format", "", "The repository format, such as maven2, npm, nuget, etc.")
	username := flag.String("username", "admin", "Nexus username (required)")
	password := flag.String("password", "admin123", "Nexus username's password (required)")
	target := flag.String("target", "", "Group repository to add the repositories to; in case of create: name of the repository to create")
	sourceRepo := flag.String("sourceRepo", "", "Repository to add the target, if more than one, use separator flag for specifying the separator for splitting the values")
	remote := flag.String("remote", "", "Remote URL for the target if the create's repoType is proxy")
	repoPolicy := flag.String("policy", "snapshot", "snapshot or release")
	separator := flag.String("separator", "\n", "Separator for list of outputs or for inputs (sourceRepo)")
	linkedRepositoryFilter := flag.String("linkedRepositoryFilter", "", "Filter for linked repositories (Not a regex) (optional)")
	notLinked := flag.Bool("notLinked", false, "Set this flag if you do NOT want repos that link with 'linkedRepositoryFilter'")
	flag.Parse()

	if *nexusUrl == "" {
		log.Fatal("nexusUrl is a required parameter")
	}

	switch *action {
	case "create":
		handleCreate(*username, *password, *nexusUrl, *target, *provider, *repoType, strings.ToUpper(*repoPolicy), *remote, *verbose)
	case "delete":
		handleDelete(*username, *password, *nexusUrl, *target, *repoType, *verbose)
	case "find":
		/*
		 * For testing the Regex: http://www.regexplanet.com/advanced/java/index.html
		 * For reading about what each part of the Regex does: http://www.vogella.com/tutorials/JavaRegularExpressions/article.html#regular-expressions
		 *
		 */
		if *verbose {
			fmt.Printf("nexusUrl: %s\nprovider: %s\ntype: %s\nfilter: %s\n", *nexusUrl, *provider, *repoType, *regexFilter)
		}
		filter, err := regexp.Compile(*regexFilter)
		if err != nil {
			log.Fatal("There is a problem with your filter.\n")
			return
		}
		if *repoType == "group" {
			groupRepos, _ := backend.GetGroupRepositories(*nexusUrl)
			groupRepositories, err := shared.FilterGroupRepositories(groupRepos, *format, *filter, *verbose)
			if err != nil {
				fmt.Printf("error filtering grouprepositories: %v\n", err)
				return
			}
			if *linkedRepositoryFilter != "" {
				groupRepositories, err = shared.FilterLinkedRepositories(*nexusUrl, groupRepositories, *linkedRepositoryFilter, !*notLinked, *verbose)
			}
			shared.PrintGroupRepositories(groupRepositories, *separator)
		} else {
			repos, _ := backend.GetRepositories(*nexusUrl)
			repositories, err := backend.FilterUserRepositories(*nexusUrl, repos, *provider, *repoType, *filter, *verbose)
			if err != nil {
				fmt.Printf("error filtering repositories: %v\n", err)
				return
			}
			shared.PrintRepositories(repositories, *separator)
		}

	case "add":
		if *username == "" || *password == "" || *target == "" {
			log.Fatalf("username, password and target required for action %s", action)
		}
		user := model.User{Username: *username, Password: *password}
		if *sourceRepo == "" {
			log.Fatalf("sourceRepo required for action %s", action)
		} else {
			backend.AddRepoToGroupRepo(*nexusUrl, user, *sourceRepo, *separator, *target, *verbose)
		}
	default:
		panic(fmt.Sprintf("Action '%v' not recognized", *action))
	}
}

func handleDelete(username string, password string, nexusUrl string, target string, repoType string, verbose bool) {
	if nexusUrl == "" || target == "" || repoType == "" {
		log.Fatalf("create requires target, nexusUrl and type arguments")
		return
	}
	if username == "" || password == "" {
		log.Fatalf("username, password required for action create")
		return
	}
	user := model.User{Username: username, Password: password}
	backend.DeleteRepository(nexusUrl, user, target, repoType, verbose)
}

func handleCreate(username string, password string, baseUrl string, targetRepo string, provider string, repoType string, repoPolicy string, remote string, verbose bool) {
	if targetRepo == "" || provider == "" || repoType == "" || repoPolicy == "" {
		log.Fatalf("create requires target, provider, type and policy arguments")
		return
	}
	if username == "" || password == "" {
		log.Fatalf("username, password required for action create")
		return
	}
	user := model.User{Username: username, Password: password}

	switch repoType {
	case "hosted":
		backend.Create(baseUrl, user, targetRepo, provider, repoType, repoPolicy, verbose)
	case "group":
		backend.CreateGroup(baseUrl, user, targetRepo, provider, verbose)
	case "proxy":
		backend.CreateProxy(baseUrl, user, targetRepo, provider, repoType, repoPolicy, remote, verbose)
	default:
		panic(fmt.Sprintf("RepoType %v not recognized", repoType))
	}
}
