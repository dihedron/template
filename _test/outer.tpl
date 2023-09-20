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

We are also running an API call to a REST API:
{{ $response := api "https://openlibrary.org/api/books?bibkeys=ISBN:0201558025,LCCN:93005405&format=json" }} 
* URL: {{ $response.URL }}
* Status Code: {{ $response.Code }}
* Headers:
{{ range $key, $value := $response.Headers -}} 
** {{ $key }} => {{ $value }}
{{ end -}} 
* Payload: {{ range $isdn, $book := $response.Payload }}
** {{ $isdn }}: {{ range $property, $value := $book }}
*** {{ $property }} => {{ $value -}}
{{ end }}{{ end }}