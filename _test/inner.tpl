Prior Experience:
{{ range .Experiences -}} 
  * since {{ .From | yellow }} until {{ .To | yellow }} as {{ .Description | green }}  
{{ end }}