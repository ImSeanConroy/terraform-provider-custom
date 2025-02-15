# Manage example {{snakeCase .Name}}.
resource "custom_{{snakeCase .Name}}" "example" {
{{- range  .Attributes}}
  {{.TfName}} = {{if eq .Type "String"}}"{{.Example}}"{{else}}{{.Example}}{{end}}
{{- end}}
}