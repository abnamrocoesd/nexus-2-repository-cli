nexus-2-repository-cli is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

nexus-2-repository-cli is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with nexus-2-repository-cli. If not, see http://www.gnu.org/licenses/.

# nexus-2-repository-cli
A CLI tool for managing Nexus 2 Repositories.

***Only works with Nexus Repository Version 2.x***

## Context

This contains a library that can interact with Nexus Repository. It is written in Go.
See the documentation of the 'action' parameter for a list of possible operations that can be performed.

## Usage

Documentation from the commandline tool (cli):
```
Usage of ./nexus-repository-cli:
  -action string
	- find: find a repository. Optional arguments: 'type', 'format', 'filter', 'separator', 'linkedRepositoryFilter'
	- add: add repository/repositories to a groupRepository. Required: target, sourceRepo. (optional: separator, if more than one repo is listed in sourceRepo)
	- create: create new repository. Required arguments: username, password, type, provider, target (optional: remote, if repoType is 'proxy')
	- delete: delete a repository. Required arguments: username, password, type, target
	 (default "find")
  -filter string
    	Optional: regex filter for the names of repositories (default "(.*)")
  -format string
    	The repository format, such as maven2, npm, nuget, etc.
  -linkedRepositoryFilter string
    	Filter for linked repositories (Not a regex) (optional)
  -nexusUrl string
    	Nexus host name (full base url), e.g. http://localhost:8081/nexus
  -notLinked
    	Set this flag if you do NOT want repos that link with 'linkedRepositoryFilter'
  -password string
    	Nexus username's password (required) (default "admin123")
  -policy string
    	snapshot or release (default "snapshot")
  -provider string
    	The repository provider, such as maven2, npm, nuget, etc.
  -remote string
    	Remote URL for the target if the create's repoType is proxy
  -separator string
    	Separator for list of outputs or for inputs (sourceRepo) (default "\n")
  -sourceRepo string
    	Repository to add the target, if more than one, use separator flag for specifying the separator for splitting the values
  -target string
    	Group repository to add the repositories to; in case of create: name of the repository to create
  -type string
    	Optional. Values: hosted, proxy, virtual, group (default "hosted")
  -username string
    	Nexus username (required)(default "admin")
  -verbose
    	Set the flag if you want verbose output
```

### Examples

```bash
./nexus-repository-cli -action create -type proxy -target random-proxy -provider maven2 -remote https://repo.jenkins-ci.org/releases/ -policy release
Success: new resource created
```

```bash
./nexus-repository-cli -action create -type group -target random-group -provider maven2
Success: new resource created
```

```bash
./nexus-repository-cli -action create -type hosted -target random-release -provider maven2 -policy release
Success: new resource created
```

```bash
./nexus-repository-cli -action add -target random-group -sourceRepo random-release,random-proxy -separator ,
Success: linked repositories 
```

```bash
./nexus-repository-cli -action find -separator ,  -filter "([a-zA-Z]*)"
random-release,releases,snapshots,thirdparty
``` 

```bash
./nexus-repository-cli -action find -separator , -type group  -linkedRepositoryFilter random-proxy         
random-group
```


## Building

To build the project yourself, do the following in the project root-folder
```
# bash-only
export GOPATH=`pwd`
go build com/abnamro/solo/nexus-repository-cli
```
