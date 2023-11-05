{{ $syntax := .Syntax }}
{{ .Global.Pre }}


{{- range .Steps }}
{{ .Pre }}

``` {{ $syntax }}
{{ .Exec }}
```

{{- if ne .Stdout "" }}
Out:
``` {{ $syntax }}
{{ .Stdout }}
```
{{ end }}

{{- if ne .Stderr "" }}
Err:
``` {{ $syntax }}
{{ .Stderr }}
```
{{ end }}

{{- if ne .Err nil }}
Fatal error, aborting:
``` {{ $syntax }}
{{ .Err }}
```
{{ end }}

{{ .Post }}
{{ end }}
