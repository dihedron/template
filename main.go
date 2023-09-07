package main

import (
	"os"

	"log/slog"

	"github.com/jessevdk/go-flags"
)

func main() {
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
