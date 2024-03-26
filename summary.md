# Version Bump

{{ if .DidFindPR }}{{ if .FinalBump | eq "none" }}Not bumping version as indicated by label.{{ else }}This version will be a {{ .FinalBump }} bump.{{ if .DidFindLabel }} Found the corresponding label on the PR.{{ else }} Using default bump as no applicable label was found.{{ end }}{{ end}}

{{ if len .AllLabels | eq 0 }}No labels were found on the PR.{{ else }}The following labels were detected on the [PR ({{ .PRNumber }})]({{ .PRLink }}):
{{- range .AllLabels }}
- `{{ . -}}`{{ end }}{{ end }}{{ else }}No PRs found. Defaulting to {{ if .FinalBump | eq "none" }}not bump.{{ else }}a {{ .FinalBump }} bump.{{ end }}{{ end }}
