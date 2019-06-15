Static website generated from the markdown list at https://github.com/Kickball/awesome-selfhosted

To generate the static files to run this site:

Install go https://golang.org/doc/install

Get your github api token https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/

go run main.go -path "/path/to/awesome-selfhosted/README.md" -ghtoken {your github api token here}

or

go run main.go -path "C:\path\to\awesome-selfhosted\README.md" -ghtoken {your github api token here}

The github api token is required to utilize githubs graphql api.
