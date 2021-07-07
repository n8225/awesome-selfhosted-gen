# Project Archived
See https://github.com/nodiscc/hecat
# Awesome selfhosted

### The "awesome-selfhosted" website generator.


Static website generated from the markdown list at https://github.com/Kickball/awesome-selfhosted


## How

To generate the static files to run this site:

- Install Go
  - https://golang.org/doc/install
- Get your GitHub API token
  - https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
- Run this
  ```
  git clone 
  cd 
  build -o list-gen main.go  //optionally go run main.go
  ```

- Commands
  
  ```
    ./list-gen -help
  USAGE
    list-gen <subcommand> [flags]

  SUBCOMMANDS
    all               Parses README.md into yaml files then generates markdown and json from those yaml files.
    parse             Parses README.md into yaml files. Requires README.md file(-readme_path).
    generate          Generates markdown and json from yaml files. Required flag: -gh_token
    generateMarkdown  Generates markdown from yaml files.

  FLAGS
    -clean false                                  Clean all files from list directory.
    -github_token ...                             Token for github API
    -readme_path ../awesome-selfhosted/README.md  Path to README.md; Default: ../awesome-selfhosted/README.md
  ```
- Most flags may be provided as env vars as well.
   | Flag | env var | description |
   | --- | --- | --- |
   | -github_token | GITHUB_TOKEN | Token for github API |
   | -readme_path | README_PATH | Path to README.md; Default: ../awesome-selfhosted/README.md |
To add an entry from a yaml file:

- Create a `.yaml` file in the `add` directory.
- Then:
  ```
  ./list-gen generate
  ```

- NOTE: The GitHub API token is required to utilize GitHub's graphql API.

- TODO:
  - Update static site UI to look better
  - Rewrite static site .js for faster load times
  - Cleanup a lot of code
