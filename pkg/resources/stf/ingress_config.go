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
  [entrypoints.proxy]
		address = ":8880"
  {{ if .UseSSL }}
  [entryPoints.websecure]
    address = ":443"
    [entryPoints.websecure.http.tls]
  {{ else }}
  [entryPoints.web]
    address = ":80"
  {{ end }}
	{{- range $svcName, $svcAttrs := .Services }}
	{{- if $svcAttrs.IsProvider }}
  [entryPoints.{{ $svcName }}]
    address = ":{{ $svcAttrs.Port }}"
	{{ end -}}
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
    {{ range $svcName, $svcAttrs := .Services }}
    [http.services.{{ $svcName }}.loadBalancer]
			{{- range $idx, $ep := $svcAttrs.Endpoints }}
      [[http.services.{{ $svcName }}.loadBalancer.servers]]
        url = "http://{{ $ep }}:{{ $svcAttrs.Port }}/"
			{{- end }}
			[http.services.{{ $svcName }}.loadBalancer.sticky.cookie]
    {{ end }}

  [http.routers]
    {{ if .DashboardEnabled }}
    [http.routers.dashboard]
      rule = "{{ .DashboardRule }}"
      entryPoints = ["{{ if .UseSSL }}websecure{{ else }}web{{ end }}"]
      service = "api@internal"
      {{ if .UseSSL -}}
      middlewares = ["whitelist"]
      {{ end -}}
    {{ end -}}
    {{ $ssl := .UseSSL -}}
    {{ range $svcName, $svcAttrs := .Services }}
    [http.routers.{{ $svcName }}]
      rule = "{{ $svcAttrs.Rule }}"
      entryPoints = ["{{ if $ssl }}websecure{{ else }}web{{ end }}"]
      service = "{{ $svcName }}"
      priority = {{ $svcAttrs.Priority }}
      {{ if $svcAttrs.Middlewares -}}
      middlewares = {{ toJson $svcAttrs.Middlewares }}
      {{ end -}}
    {{ end }}
    {{ range $proxyName, $proxyAttrs := .Proxies }}
		[http.routers.{{ $proxyName }}-proxy]
      rule = "{{ $proxyAttrs.Rule }}"
      entryPoints = ["proxy"]
      service = "{{ $proxyName }}"
      priority = {{ $proxyAttrs.Priority }}
		{{ end }}

[tcp]

  [tcp.services]
    {{ range $svcName, $svcAttrs := .Services -}}
    {{- if $svcAttrs.IsProvider }}
		{{- range $idx, $ep := $svcAttrs.Endpoints }}
		[[tcp.services.{{ $svcName }}.loadBalancer.servers]]
			address = "{{ $ep }}:{{ $svcAttrs.Port }}"
		{{- end }}
    {{- end }}
    {{- end }}

  [tcp.routers]
		{{ $backtick := .Backtick -}}
		{{ range $svcName, $svcAttrs := .Services -}}
		{{ if $svcAttrs.IsProvider }}
    [tcp.routers.{{ $svcName }}]
      service = "{{ $svcName }}"
      entrypoints = ["{{ $svcName }}"]
      rule = "HostSNI({{ $backtick }}*{{ $backtick }})"
		{{ end -}}
		{{ end }}

`))
