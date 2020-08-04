# Awesome selfhosted

The "awesome-selfhosted" website generator.


## What

Static website generated from the markdown list at

- https://github.com/Kickball/awesome-selfhosted


## How

To generate the static files to run this site:

- Install Go
  - https://golang.org/doc/install
- Get your GitHub API token
  - https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
- Run this
  ```
  go run cmd/mdconvert/main.go -path "/path/to/awesome-selfhosted/README.md" -ghtoken {your github api token here}
  ```
  or that
  ```
  go run cmd/mdconvert/main.go -path "C:\path\to\awesome-selfhosted\README.md" -ghtoken {your github api token here}
  ```

To add an entry from a yaml file:

- Create a `.yaml` file in the `add` directory.
- Run this
  ```
  go run cmd/addtoyaml/main.go -ghtoken {your github api token here}
  ```

The GitHub API token is required to utilize GitHub's graphql API.
