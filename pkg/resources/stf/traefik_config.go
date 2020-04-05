package stf

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

// traefikStaticConfigTmpl is the template used for traefik static configurations
var traefikStaticConfigTmpl = template.Must(
	template.New("traefik-config").
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(`
[global]
  checkNewVersion = true
  sendAnonymousUsage = false

[api]
  insecure = {{ if .UseSSL }}false{{ else }}true{{ end }}
  dashboard = true

[log]
  level = "INFO"

{{ if .AccessLog -}}
[accessLog]
{{ end -}}

[entryPoints]

  [entrypoints.stf-internal-proxy]
		address = ":8880"

  [entryPoints.stf-web]
    address = "{{ if .UseSSL }}:8443{{ else }}:8088{{ end }}"
    {{ if .UseSSL }}[entryPoints.stf-web.http.tls]{{ end }}

  {{ range $idx, $svc := .Services }}
  {{- range $svcName, $svcAttrs := $svc }}
  {{- if $svcAttrs.IsProvider }}
  [entryPoints.{{ $svcName }}]
    address = ":{{ $svcAttrs.Port }}"
  {{- end }}
  {{- end }}
  {{ end }}

[providers]
  [providers.file]
    watch = true
    directory = "/etc/configmap/routes"
`))

// traefikDynamicConfigTmpl is the template used for dynamic traefik configurations
var traefikDynamicConfigTmpl = template.Must(
	template.New("traefik-routes").
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(`
{{ if and .UseSSL (not .UseSelfSigned) }}
[tls]
  [tls.stores]
    [tls.stores.default]
      [tls.stores.default.defaultCertificate]
        certFile = "/etc/traefik/ssl/tls.crt"
        keyFile  = "/etc/traefik/ssl/tls.key"
{{ end }}
[http]

  [http.middlewares]
    [http.middlewares.strip-rethinkdb.stripPrefix]
      prefixes = ["/rethinkdb"]
  {{ if and .UseSSL .DashboardEnabled }}
    [http.middlewares.whitelist.ipWhiteList]
    sourceRange = {{ .DashboardWhitelist }}
  {{ end }}

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
    {{ if .DashboardEnabled }}
    [http.routers.dashboard]
      rule = "{{ .DashboardRule }}"
      entryPoints = ["stf-web"]
      service = "api@internal"
      {{ if .UseSSL -}}
      middlewares = ["whitelist"]
      {{ end -}}
    {{ end -}}
    {{ range $idx, $svc := .Services }}
    {{ range $svcName, $svcAttrs := $svc }}
    [http.routers.{{ $svcName }}]
      rule = "{{ $svcAttrs.Rule }}"
      entryPoints = ["stf-web"]
      service = "{{ $svcName }}"
      priority = {{ $svcAttrs.Priority }}
      {{ if $svcAttrs.Middlewares -}}
      middlewares = {{ toJson $svcAttrs.Middlewares }}
      {{ end -}}
    {{- end }}
    {{ end }}
    {{ range $idx, $proxy := .Proxies }}
    {{ range $proxyName, $proxyAttrs := $proxy }}
    [http.routers.{{ $proxyName }}-proxy]
      rule = "{{ $proxyAttrs.Rule }}"
      entryPoints = ["stf-internal-proxy"]
      service = "{{ $proxyName }}"
      priority = {{ $proxyAttrs.Priority }}
		{{- end }}
    {{ end }}

[tcp]

  [tcp.services]
    {{ range $idx, $svc := .Services }}
    {{- range $svcName, $svcAttrs := $svc -}}
    {{- if $svcAttrs.IsProvider }}
    {{- range $idx, $ep := $svcAttrs.Endpoints }}
    [[tcp.services.{{ $svcName }}.loadBalancer.servers]]
      address = "{{ $ep }}:{{ $svcAttrs.Port }}"
    {{- end }}
    {{- end }}
    {{- end }}
    {{- end }}

  [tcp.routers]
    {{ $backtick := .Backtick -}}
    {{ range $idx, $svc := .Services }}
    {{ range $svcName, $svcAttrs := $svc -}}
    {{ if $svcAttrs.IsProvider }}
    [tcp.routers.{{ $svcName }}]
      service = "{{ $svcName }}"
      entrypoints = ["{{ $svcName }}"]
      rule = "HostSNI({{ $backtick }}*{{ $backtick }})"
    {{ end -}}
    {{ end -}}
    {{ end }}

`))
