{{- define "eth2-common.genesis.formatPreminedAddresses" }}
{{- $addresses := . }}
{{- if $addresses }}
export PREMINE_ADDRS='
{{- range $addr := $addresses }}
  "{{ $addr }}": 1000000000ETH
{{- end }}'
{{- else }}
export PREMINE_ADDRS='{}'
{{- end }}
{{- end }}
