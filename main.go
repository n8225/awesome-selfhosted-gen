package main

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/n8225/awesome-selfhosted-gen/pkg/exporter"
	"github.com/n8225/awesome-selfhosted-gen/pkg/parse"
)

func main() {
	initLogs()

	var (
		rootFlagSet  = flag.NewFlagSet("list-gen", flag.ExitOnError)
		github_token = rootFlagSet.String("github_token", "", "Token for github API")
		readme_path  = rootFlagSet.String("readme_path", "../awesome-selfhosted/README.md", "Path to README.md; Default: ../awesome-selfhosted/README.md")
		clean        = rootFlagSet.Bool("clean", false, "Clean all files from list directory.")
	)

	parse := &ffcli.Command{
		Name:       "parse",
		ShortUsage: "list-gen parse [flags]",
		ShortHelp:  "Parses README.md into yaml files. Requires README.md file(-readme_path).",
		LongHelp: "Parses README.md into yaml files. Requires a README.md file as input. " +
			"\nWill look for ../awesome-selfhosted/README.md if no input is provided.",
		FlagSet: rootFlagSet,
		Exec: func(_ context.Context, args []string) error {
			cleanDir(*clean)
			parseFiles(*readme_path)
			return nil
		},
	}

	generate := &ffcli.Command{
		Name:       "generate",
		ShortUsage: "list-gen generate [flags]",
		ShortHelp:  "Generates markdown and json from yaml files. Required flag: -gh_token",
		LongHelp: "Generates markdown and json from yaml files. Requires a github personal access token(-github_token)." +
			"\nhttps://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token",
		FlagSet: rootFlagSet,
		Exec: func(_ context.Context, args []string) error {
			generateFiles(*github_token)
			return nil
		},
	}

	generateMarkdown := &ffcli.Command{
		Name:       "generateMarkdown",
		ShortUsage: "list-gen generateMarkdown [flags]",
		ShortHelp:  "Generates markdown from yaml files.",
		LongHelp:   "Generates markdown from yaml files.",
		FlagSet:    rootFlagSet,
		Exec: func(_ context.Context, args []string) error {
			exporter.ToMD()
			return nil
		},
	}

	all := &ffcli.Command{
		Name:       "all",
		ShortUsage: "list-gen parse [flags]",
		ShortHelp:  "Parses README.md into yaml files then generates markdown and json from those yaml files.",
		LongHelp:   parse.LongHelp + "\n" + generate.LongHelp,
		FlagSet:    rootFlagSet,
		Exec: func(_ context.Context, args []string) error {
			log.Info().Msg("GITHUB_TOKEN: " + *github_token)
			cleanDir(*clean)
			parseFiles(*readme_path)
			generateFiles(*github_token)
			return nil
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "list-gen <subcommand> [flags]",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{all, parse, generate, generateMarkdown},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.Parse(os.Args[1:]); err != nil {
		log.Error().Err(err)
	}

	if err := ff.Parse(rootFlagSet, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		log.Error().Err(err)
	}

	if err := root.Run(context.Background()); err != nil {
		log.Error().Err(err)
	}
}

func parseFiles(path string) {
	apath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal().Stack().Err(err).Stack().Err(err)
	}
	list := new(parse.List)
	list.Entries = parse.MdParser(apath)
	exporter.ToYamlFiles(*list)
}

func generateFiles(ghToken string) {
	if ghToken == "" {
		ghToken = os.Getenv("GITHUB_TOKEN")
	}
	if ghToken == "" {
		log.Fatal().Msg("Github API token required and not provided.")
	}
	yl := exporter.ImportYaml(ghToken)
	exporter.ToJSON(yl, "list")
	exporter.ToYAML(yl, "output/list")
}

func cleanDir(clean bool) {
	if clean {
		log.Info().Msg("Cleaning './list/' directory")
		err := os.RemoveAll("./list/")
		if err != nil {
			log.Error().Stack().Err(err)
		}
	}
}

func initLogs() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't create Log file")
	}
	var logW io.Writer = logFile
	prettyLog := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	jsonFileLog := zerolog.New(logW)
	multiLogger := zerolog.MultiLevelWriter(prettyLog, jsonFileLog)
	log.Logger = zerolog.New(multiLogger).With().Timestamp().Logger()
}
