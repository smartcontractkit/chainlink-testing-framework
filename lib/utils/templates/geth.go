package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var InitGethScript = `
#!/bin/bash
if [ ! -d /root/.ethereum/keystore ]; then
	echo "/root/.ethereum/keystore not found, running 'geth init'..."
	geth init --datadir /root/.ethereum/devchain /root/genesis.json
	echo "...done!"
fi

geth "$@"
`

var GenesisJson = `
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
		"period": 2,
		  "epoch": 30000
		}
	},
	"nonce": "0x0000000000000042",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",

	"difficulty": "0x1",
	"coinbase": "0x3333333333333333333333333333333333333333",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"extraData": "{{ .ExtraData }}",
	"gasLimit": "0xE0000000",
	"alloc": {
	{{- $lastIndex := decrement (len $.AccountAddr)}}
	{{- range $i, $addr := .AccountAddr }}
  "{{$addr}}": {
    "balance": "20000000000000000000000"
  }{{ if ne $i $lastIndex }},{{ end }}
{{- end }}
	}
  }`

var funcMap = template.FuncMap{
	// The name "inc" is what the function will be called in the template text.
	"decrement": func(i int) int {
		return i - 1
	},
}

type GenesisConsensus = string

const (
	GethGenesisConsensus_Ethash = "ethash"
	GethGenesisConsensus_Clique = "clique"
)

type GenesisJsonTemplate struct {
	AccountAddr []string
	ChainId     string
	Consensus   GenesisConsensus
	ExtraData   string
}

// String representation of the job
func (c GenesisJsonTemplate) String() (string, error) {
	var consensusStr string
	switch c.Consensus {
	case GethGenesisConsensus_Ethash:
		consensusStr = `,"ethash": {}`
	case GethGenesisConsensus_Clique:
		consensusStr = `,"clique": {"period": 1,"epoch": 30000}`
	default:
		consensusStr = ""
	}

	extraData := c.ExtraData

	if c.Consensus == GethGenesisConsensus_Clique && (c.ExtraData == "" || c.ExtraData == "0x") {
		return "", fmt.Errorf("extraData is required for clique consensus")
	} else if extraData == "" {
		extraData = "0x"
	}

	data := struct {
		AccountAddr []string
		ChainId     string
		ExtraData   string
		Consensus   string
	}{
		AccountAddr: c.AccountAddr,
		ChainId:     c.ChainId,
		ExtraData:   extraData,
		Consensus:   consensusStr,
	}

	tpl := `
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
	  "londonBlock": 0
	  {{ .Consensus }}
	},
	"nonce": "0x0000000000000042",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"difficulty": "0x20000",
	"coinbase": "0x0000000000000000000000000000000000000000",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"extraData": "{{ .ExtraData }}",
	"gasLimit": "8000000000",
	"alloc": {
		{{- $lastIndex := decrement (len $.AccountAddr)}}
		{{- range $i, $addr := .AccountAddr }}
	  "{{$addr}}": {
		"balance": "9000000000000000000000000000"
	  }{{ if ne $i $lastIndex }},{{ end }}
	{{- end }}
	}
  }`

	t, err := template.New("genesis-json").Funcs(funcMap).Parse(tpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}
