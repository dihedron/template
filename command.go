package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/dihedron/rawdata"
	"github.com/dihedron/template/extensions"
	"github.com/jlaffaye/ftp"
	"golang.org/x/exp/slog"
)

// Command is the main command for the application.
type Command struct {
	Input     *Input   `short:"i" long:"input" description:"The input data, either as an inline JSON value or as a @file (in JSON or YAML format)." optional:"yes" env:"TEMPLATE_INPUT"`
	Templates []string `short:"t" long:"template" description:"The paths of all the templates and subtemplates on disk; the main template must be the first." required:"yes"`
	Output    string   `short:"o" long:"output" description:"The path to the output file." optional:"yes" env:"TEMPLATE_OUTPUT"`
}

// Input is the input data for the template.
type Input struct {
	Data any
}

// UnmarshalFlag unmarshals the input data from a string;
// if the string starts with '@', it is treated as a file path,
// otherwise it is treated as an inline JSON value.
func (i *Input) UnmarshalFlag(value string) error {
	var err error
	i.Data, err = rawdata.Unmarshal(value)
	if err != nil {
		slog.Error("cannot unmarshal input data", "error", err)
	}
	return err
}

// Execute executes the command.
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
	// ensure the output is closed
	if output, ok := output.(io.WriteCloser); ok {
		defer output.Close()
	}

	// populate the functions map
	functions := template.FuncMap{}
	maps.Copy(functions, extensions.FuncMap())
	maps.Copy(functions, sprig.FuncMap())

	// // download the templates if needed
	// for _, t := range cmd.Templates {
	// 	if strings.HasPrefix(t, "http://") || strings.HasPrefix(t, "https://") {
	// 		data, err := downloadHTTP(t)
	// 		if err != nil {
	// 			slog.Error("cannot download template", "error", err)
	// 			return fmt.Errorf("error downloading template %s: %w", t, err)
	// 		}
	// 		cmd.Templates = append(cmd.Templates, string(data))
	// 	} else if strings.HasPrefix(t, "ftp://") || strings.HasPrefix(t, "ftps://") || strings.HasPrefix(t, "sftp://") {
	// 		data, err := downloadFTP(t)
	// 		if err != nil {
	// 			slog.Error("cannot download template", "error", err)
	// 			return fmt.Errorf("error downloading template %s: %w", t, err)
	// 		}
	// 		cmd.Templates = append(cmd.Templates, string(data))
	// 	}
	// }

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

func read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

// downloadHTTP downloads a file from an HTTP or HTTPS server.
// The URL can be in the format:
//
//	http://[<user>[:<password>]@]<host>[:<port>][<path>]
//	https://[<user>[:<password>]@]<host>[:<port>][<path>]
//
// where:
//
// - <user> is the username
// - <password> is the password
// - <host> is the hostname
// - <port> is the port number
// - <path> is the path to the file
func downloadHTTP(rawURL string) ([]byte, error) {
	// send the GET request
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	// ensure the body is closed to prevent memory leaks
	defer resp.Body.Close()

	// check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// read the entire body into a []byte
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// downloadFTP downloads a file from an FTP or SFTP server.
// The URL can be in the format:
//
//	ftp://[<user>[:<password>]@]<host>[:<port>][<path>]
//	ftps://[<user>[:<password>]@]<host>[:<port>][<path>]
//	sftp://[<user>[:<password>]@]<host>[:<port>][<path>]
//
// where:
//
// - <user> is the username
// - <password> is the password
// - <host> is the hostname
// - <port> is the port number
// - <path> is the path to the file
func downloadFTP(rawURL string) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	scheme := strings.ToLower(u.Scheme)

	user := ""
	pass := ""
	if u.User != nil {
		if u.User.Username() != "" && user == "" {
			user = u.User.Username()
		}
		if p, ok := u.User.Password(); ok && pass == "" {
			pass = p
		}
	}

	var useTLS bool
	switch scheme {
	case "ftp":
		useTLS = false
	case "ftps":
		useTLS = true
	default:
		return nil, fmt.Errorf("unsupported protocol for FTP download: %s", scheme)
	}

	host := u.Host
	path := u.Path

	if !strings.Contains(host, ":") {
		if useTLS {
			host += ":990"
		} else {
			host += ":21"
		}
	}

	c, err := ftp.Dial(host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("ftp dial: %w", err)
	}
	defer c.Quit()

	if useTLS {
		err = c.AuthTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, fmt.Errorf("ftp auth tls: %w", err)
		}
	}

	// Login with anonymous as default if user is empty
	if user == "" {
		user = "anonymous"
		pass = "anonymous"
	}

	err = c.Login(user, pass)
	if err != nil {
		return nil, fmt.Errorf("ftp login: %w", err)
	}

	r, err := c.Retr(path)
	if err != nil {
		return nil, fmt.Errorf("ftp retrieve: %w", err)
	}
	defer r.Close()

	return io.ReadAll(r)
}
