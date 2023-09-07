package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/Masterminds/sprig/v3"
	"github.com/dihedron/rawdata"
	"github.com/dihedron/template/formatting"
	"golang.org/x/exp/slog"
)

type Command struct {
	Input      *Input   `short:"i" long:"input" description:"The input data, either as an inline JSON value or as a @file (in JSON or YAML format)." otional:"yes" env:"MOULD_INPUT"`
	Templates  []string `short:"t" long:"template" description:"The paths of all the templates and subtemplates on disk." required:"yes"`
	Output     string   `short:"o" long:"output" description:"The path to the output file." optional:"yes" env:"MOULD_OUTPUT"`
	LogEnabled bool     `short:"l" long:"log-enabled" description:"Whether to log to STDERR." optional:"yes" env:"MOULD_LOG_ENABLED"`
}

type Input struct {
	Data interface{}
}

func (i *Input) UnmarshalFlag(value string) error {
	var err error
	i.Data, err = rawdata.Unmarshal(value)
	if err != nil {
		slog.Error("cannot unmarshal input data", "error", err)
	}
	return err
}

func (cmd *Command) Execute(args []string) error {
	var err error

	// if the input map is nil, then the input data is
	// provided via STDIN, and that's where we take it
	if cmd.Input == nil {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("cannot read data from STDIN", "error", err)
			return fmt.Errorf("error reading input data from STDIN: %w", err)
		}
		cmd.Input = &Input{}
		if err = cmd.Input.UnmarshalFlag(string(input)); err != nil {
			return err
		}
	}

	// prepare the output stream
	var output io.Writer
	if cmd.Output != "" {
		path := filepath.Dir(cmd.Output)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			slog.Error("cannot create output directory", "error", err)
			return fmt.Errorf("error creating output directory %s: %w", path, err)
		}
		output, err = os.Create(cmd.Output)
		if err != nil {
			slog.Error("cannot create output file", "error", err)
			return fmt.Errorf("error creating output file %s: %w", cmd.Output, err)
		}
	} else {
		output = os.Stdout
	}

	// populate the functions map
	functions := template.FuncMap{}
	for k, v := range formatting.FuncMap() {
		functions[k] = v
	}
	for k, v := range sprig.FuncMap() {
		functions[k] = v
	}

	// parse the templates
	main := path.Base(cmd.Templates[0])
	templates, err := template.New(main).Funcs(functions).ParseFiles(cmd.Templates...)
	if err != nil {
		slog.Error("cannot parse template files", "templates", cmd.Templates, "error", err)
		return fmt.Errorf("error parsing template files %v: %w", cmd.Templates, err)
	}

	// execute the template
	if err := templates.ExecuteTemplate(output, main, cmd.Input.Data); err != nil {
		slog.Error("cannot apply data to template", "error", err)
		return fmt.Errorf("error applying data to template: %w", err)
	}
	return nil
}
