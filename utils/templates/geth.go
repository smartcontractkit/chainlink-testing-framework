package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/google/uuid"
)

var InitGethScript = `
#!/bin/bash
if [ ! -d /root/.ethereum/keystore ]; then
	echo "/root/.ethereum/keystore not found, running 'geth init'..."
	geth init /root/genesis.json
	echo "...done!"
fi

geth "$@"
`

var BootnodeScript = `
#!/bin/bash

echo "Starting bootnode"
bootnode --genkey=/root/.ethereum/node.key
echo "Bootnode key generated"

echo "Bootnode address written to file"
bootnode -writeaddress --nodekey=/root/.ethereum/node.key > /root/.ethereum/bootnodes
cat /root/.ethereum/bootnodes

echo "Starting Bootnode"
bootnode --nodekey=/root/.ethereum/node.key --verbosity=6 --addr :30301
`

var InitNonDevGethScript = `
#!/bin/bash

echo "Starting geth"
geth --datadir=/root/.ethereum init /root/genesis.json

bootnode_enode=$(cat /root/.ethereum/keystore/password.txt)
echo "Bootnode enode: $bootnode_enode"
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
	"difficulty": "1",
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

func BuildGenesisJsonForNonDevChain(chainId string, accountAddr []string, extraData string) (string, error) {
	data := struct {
		AccountAddr []string
		ChainId     string
		ExtraData   string
	}{
		AccountAddr: accountAddr,
		ChainId:     chainId,
		ExtraData:   extraData,
	}

	t, err := template.New("genesis-json").Funcs(funcMap).Parse(GenesisJson)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}

type GethPoWGenesisJsonTemplate struct {
	AccountAddr string
	ChainId     string
}

// String representation of the job
func (c GethPoWGenesisJsonTemplate) String() (string, error) {
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
	},
	"nonce": "0x0000000000000042",
	"mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"difficulty": "1",
	"coinbase": "0x3333333333333333333333333333333333333333",
	"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	"extraData": "0x",
	"gasLimit": "8000000000",
	"alloc": {
	  "{{ .AccountAddr }}": {
		"balance": "20000000000000000000000"
	  }
	}
  }`
	return MarshalTemplate(c, uuid.NewString(), tpl)
}
