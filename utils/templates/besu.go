package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var BesuGenesisJson = `
{
	"config": {
	  "chainId": {{ .ChainId }},
	  "homesteadBlock": 0,
	  "eip150Block": 0,
	  "eip155Block": 0,
	  "eip158Block": 0,
	  "eip160Block": 0,
	  "byzantiumBlock": 0,
	  "constantinopleBlock": 0,
	  "petersburgBlock": 0,
	  "istanbulBlock": 0,
	  "muirGlacierBlock": 0,
	  "berlinBlock": 0,
	  "londonBlock": 0,
	  "clique": {
		"blockperiodseconds": 2,
		"epochlength": 30000
      },
	  "ethash": {
		"fixeddifficulty": 100
	  }
	},
	"nonce": "0x42",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"difficulty": "0x10000",
	"coinbase": "0x0000000000000000000000000000000000000000",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"extraData": "{{ .ExtraData }}",
	"gasLimit": "0x1fffffffffffff",
	"alloc": {
	{{- $lastIndex := decrement (len $.AccountAddr)}}
	{{- range $i, $addr := .AccountAddr }}
  "{{$addr}}": {
    "balance": "20000000000000000000000"
  }{{ if ne $i $lastIndex }},{{ end }}
{{- end }}
	}
  }`

func BuildBesuGenesisJsonForNonDevChain(chainId string, accountAddr []string, extraData string) (string, error) {
	data := struct {
		AccountAddr []string
		ChainId     string
		ExtraData   string
	}{
		AccountAddr: accountAddr,
		ChainId:     chainId,
		ExtraData:   extraData,
	}

	t, err := template.New("genesis-json").Funcs(funcMap).Parse(BesuGenesisJson)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}
