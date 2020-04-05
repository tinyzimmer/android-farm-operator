package stf

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

// traefikStaticConfigTmpl is the template used for traefik static configurations
var providerTraefikStaticConfigTmpl = template.Must(
	template.New("traefik-config").
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(`
[global]
  checkNewVersion = true
  sendAnonymousUsage = false

[api]
  insecure = true
  dashboard = true

[log]
  level = "INFO"

{{ if .AccessLog -}}
[accessLog]
{{ end -}}

[entryPoints]

  [entrypoints.websocket]
    address = ":8088"
  {{ range $idx, $svc := .Services }}
  {{ range $svcName, $svcAttrs := $svc }}
  [entryPoints.{{ $svcName }}]
    address = ":{{ $svcAttrs.Port }}"
  {{ end }}
  {{ end }}

[providers]
  [providers.file]
    watch = true
    directory = "/etc/configmap/routes"
`))

// traefikDynamicConfigTmpl is the template used for dynamic traefik configurations
var providerTraefikDynamicConfigTmpl = template.Must(
	template.New("traefik-routes").
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(`

  [http.services]
    {{ range $idx, $svc := .Services }}
    {{- range $svcName, $svcAttrs := $svc }}
    [http.services.{{ $svcName }}.loadBalancer]
      {{- range $idx, $ep := $svcAttrs.Endpoints }}
      [[http.services.{{ $svcName }}.loadBalancer.servers]]
        url = "http://{{ $ep }}:{{ $svcAttrs.Port }}/"
      {{- end }}
      [http.services.{{ $svcName }}.loadBalancer.sticky.cookie]
    {{- end }}
    {{ end }}

  [http.routers]
    {{ range $idx, $svc := .Services }}
    {{- range $svcName, $svcAttrs := $svc }}
    [http.routers.{{ $svcName }}]
      rule = "{{ $svcAttrs.Rule }}"
      entryPoints = ["websocket"]
      service = "{{ $svcName }}"
    {{- end }}
    {{ end }}

[tcp]

  [tcp.services]
    {{ range $idx, $svc := .Services }}
    {{- range $svcName, $svcAttrs := $svc -}}
    {{- range $idx, $ep := $svcAttrs.Endpoints }}
    [[tcp.services.{{ $svcName }}.loadBalancer.servers]]
      address = "{{ $ep }}:{{ $svcAttrs.Port }}"
    {{- end }}
    {{- end }}
    {{- end }}

  [tcp.routers]
    {{ $backtick := .Backtick -}}
    {{ range $idx, $svc := .Services }}
    {{- range $svcName, $svcAttrs := $svc -}}
    [tcp.routers.{{ $svcName }}]
      service = "{{ $svcName }}"
      entrypoints = ["{{ $svcName }}"]
      rule = "HostSNI({{ $backtick }}*{{ $backtick }})"
    {{- end }}
    {{ end }}
`))
