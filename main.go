package main

import (
	"os"

	"log/slog"

	"github.com/dihedron/template/metadata"
	"github.com/jessevdk/go-flags"
)

func main() {

	defer cleanup()

	if len(os.Args) == 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		metadata.Print(os.Stdout)
		os.Exit(0)
	} else if len(os.Args) == 3 && os.Args[1] == "version" && (os.Args[2] == "--verbose" || os.Args[2] == "-v") {
		metadata.PrintFull(os.Stdout)
		os.Exit(0)
	}

	command := Command{}
	if args, err := flags.Parse(&command); err != nil {
		os.Exit(1)
	} else {
		if command.LogEnabled {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
		}
		err = command.Execute(args)
		if err != nil {
			os.Exit(1)
		}
	}
}
