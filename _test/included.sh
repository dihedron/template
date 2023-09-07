#!/bin/bash

test_me() {
    echo "hallo, {{ if .Name -}}{{ .Name }}{{ else }}world{{ end }}!"
}