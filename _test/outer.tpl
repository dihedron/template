This is the resumé of {{.Name}} {{.Surname}}

Name: {{.Name | blue}}
Surname: {{.Surname | red}}
Phone No.: {{.PhoneNo | purple }}

{{ template "inner.tpl" . }}

{{if .Developer -}} Role: Developer {{- end}}
{{if .SysAdmin}} Role: SysAdmin {{- end}}
{{range .Emails -}} Email: 
    Description: {{ .Name }}
    Address: {{.Address}}
{{end}}

Include script as quote:
{{ include "_test/included.sh" . "> " -}}


Include script with 2-spaces indentation:
{{ include "_test/included.sh" . "  " -}}


Include script as is:
{{ include "_test/included.sh" . }}
