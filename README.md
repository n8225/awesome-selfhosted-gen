# Awesome selfhosted

The "awesome-selfhosted" website generator.


## What

Static website generated from the markdown list at

- https://github.com/Kickball/awesome-selfhosted


## How

To generate the static files to run this site:

- install Go
  - https://golang.org/doc/install
- get your GitHub API token (`your_github_api_token` in the commands below)
  - https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
- run this (Linux, macOS)
  ```
  TOKEN=your_github_api_token \
  go run cmd/mdconvert/main.go -path "/path/to/awesome-selfhosted/README.md" -ghtoken $TOKEN
  ```
  or that (Windows)
  ```
  TOKEN=your_github_api_token \
  go run cmd/mdconvert/main.go -path "C:\path\to\awesome-selfhosted\README.md" -ghtoken $TOKEN
  ```

To add an entry from a yaml file:

- create a `.yaml` file in the `add` directory
- run this
  ```
  TOKEN=your_github_api_token \
  go run cmd/addtoyaml/main.go -ghtoken $TOKEN
  ```

The GitHub API token is required to utilize GitHub's GraphQL API.
