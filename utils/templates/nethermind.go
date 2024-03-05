package templates

import "github.com/google/uuid"

type NethermindPoAGenesisJsonTemplate struct {
	AccountAddr string
	ChainId     string
	ExtraData   string
}

func (c NethermindPoAGenesisJsonTemplate) String() (string, error) {
	tpl := `
	{
		"name": "Devnet",
		"engine": {
			"clique": {
				"params": {
					"period": 1,
					"epoch": 30000
				}
			}
		},
		"params": {
		  "networkID" : "{{ .ChainId }}",
		  "gasLimitBoundDivisor": "0x400",
		  "registrar": "0x0000000000000000000000000000000000000000",
		  "accountStartNonce": "0x0",
		  "maximumExtraDataSize": "0xffff",
		  "minGasLimit": "0x1388",
		  "maxCodeSize": "0x6000",
		  "maxCodeSizeTransition": "0x0",
		  "eip150Transition": "0x0",
		  "eip158Transition": "0x0",
		  "eip160Transition": "0x0",
		  "eip161abcTransition": "0x0",
		  "eip161dTransition": "0x0",
		  "eip155Transition": "0x0",
		  "eip140Transition": "0x0",
		  "eip211Transition": "0x0",
		  "eip214Transition": "0x0",
		  "eip658Transition": "0x0",
		  "eip145Transition": "0x0",
		  "eip1014Transition": "0x0",
		  "eip1052Transition": "0x0",
		  "eip1283Transition": "0x0",
		  "eip1283DisableTransition": "0x0",
		  "eip152Transition": "0x0",
		  "eip1108Transition": "0x0",
		  "eip1344Transition": "0x0",
		  "eip1884Transition": "0x0",
		  "eip2028Transition": "0x0",
		  "eip2200Transition": "0x0",
		  "eip2565Transition": "0x0",
		  "eip2929Transition": "0x0",
		  "eip2930Transition": "0x0",
		  "eip1559Transition": "0x0",
		  "eip3198Transition": "0x0",
		  "eip3529Transition": "0x0",
		  "eip3541Transition": "0x0",
		},
		"genesis": {
		  "seal": {
			"ethereum": {
			  "nonce": "0x0000000000000042",
			  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
			}
		  },
		  "difficulty": "0x01",
		  "author": "0x0000000000000000000000000000000000000000",
		  "timestamp": "0x00",
		  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "extraData": "{{ .ExtraData }}",
		  "gasLimit": "0x17D7840"
		},
		"accounts": {
		  "{{ .AccountAddr }}": {
			"balance": "1000000000000000000000000000"
		  }
		},
	  }
	  `
	return MarshalTemplate(c, uuid.NewString(), tpl)
}
