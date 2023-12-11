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

We can also place a GET call to a REST API (JSON); the response gives us access to the URL, the headers and to the payload as an object:
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

Environment:
HTTP_PROXY={{ env "HTTP_PROXY" }}

Directory:
{{$files := listDir "."}}{{range $file := $files}}
 - {{if isFile $file}}{{$file}} (size: {{fileSize $file}} bytes){{end}}{{ end}}

