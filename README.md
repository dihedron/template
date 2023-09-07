# A Golang template engine

`template`: use Golang templates on the command line and hydrate them with variable values provided via YAML or JSON files.

### Usage

Take a look a the `_test/` directory for an example: 

1. you CAN specify an input file that will contain an arbitrary JSON or YAML structure with the input variables, or alternatively pipe it in via STDIN; in the latter case if the input is in YAML format it should start with `---`, if in JSON format with `{` or `[`; if the input is provided via the `-i` or `--input` argument *as a file*, it must be prefixed with `@`, otheriwse it is interpreted as an inline JSON/YAML;
1. the list of template files must be provided as multiple `--template` arguments pointing to their paths on disk; the first template is considered the main template (the starting point for template expansion);
1. you CAN specify the name of an output file; if left blank the application will write to STDOUT;

For example:

```bash
$> template -i=@_test/input.yaml -t=_test/outer.tpl -t=_test/inner.tpl
```

or 

```bash
$> cat _test/input.json | mason hydrate --template=_test/outer.tpl --template=_test/inner.tpl
```

If no output parameter is specified, `template` will write to STDOUT by default; thus, it can be used with pipes (`|`) where the STDIN is the set of input variables funnelled into `template` and the output goes to SDOUT and can therefore be piped into other commands.

### Custom functions

`template` provides a set of functions to format and process data.

Besides including the whole set of functions made available by the excellent [Sprig library](http://masterminds.github.io/sprig/), it provides some additional functions to enrich the template manipulation capabilities.  

#### Function `include` 

The `include` function can be used when you want to include a sub-template (or any other file) and you would like it to be padded left line by line with a fixed string. 
For instance this provides a way to include some file and have it automatically indented. Look at `_test/outer.tpl` to see how it includes a bash script prepending `> ` to each line.

#### Functions from `color`

The following functions are made available:
* `blue`
* `cyan`
* `green`
* `magenta`
* `purple`
* `red`
* `yellow`
* `white`
* `hiblue`
* `hicyan`
* `higreen`
* `himagenta`
* `hipurple`
* `hired`
* `hiyellow`
* `hiwhite`

To use the function, do as follows in your template:

```
Name: {{.Name | blue}}
Surname: {{.Surname | red}}
Phone No.: {{.PhoneNo | purple }}
```

