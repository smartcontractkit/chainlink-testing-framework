package client

type ChainlinkConfig struct {
	URL      string
	Email    string
	Password string
}

type ResponseArray struct {
	Data []map[string]interface{}
}

type Response struct {
	Data map[string]interface{}
}

func NewResponse() Response {
	return Response{
		Data: map[string]interface{}{},
	}
}

type BridgeType struct {
	Data BridgeTypeData `json:"data"`
}

type BridgeTypeData struct {
	Attributes BridgeTypeAttributes `json:"attributes"`
}

type BridgeTypeAttributes struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OCRKeys struct {
	Data []OCRKeyData `json:"data"`
}

type OCRKey struct {
	Data OCRKeyData `json:"data"`
}

type OCRKeyData struct {
	Attributes OCRKeyAttributes `json:"attributes"`
}

type OCRKeyAttributes struct {
	ID                    string `json:"id"`
	ConfigPublicKey       string `json:"configPublicKey"`
	OffChainPublicKey     string `json:"offChainPublicKey"`
	OnChainSigningAddress string `json:"onChainSigningAddress"`
}

type P2PKeys struct {
	Data []P2PKeyData `json:"data"`
}

type P2PKey struct {
	Data P2PKeyData `json:"data"`
}

type P2PKeyData struct {
	Attributes P2PKeyAttributes `json:"attributes"`
}

type P2PKeyAttributes struct {
	ID        int    `json:"id"`
	PeerID    string `json:"peerId"`
	PublicKey string `json:"publicKey"`
}

type ETHKeys struct {
	Data []ETHKeyData `json:"data"`
}

type ETHKey struct {
	Data ETHKeyData `json:"data"`
}

type ETHKeyData struct {
	Attributes ETHKeyAttributes `json:"attributes"`
}

type ETHKeyAttributes struct {
	Address string `json:"address"`
}

type SpecForm struct {
	TOML string `json:"toml"`
}

type Spec struct {
	Data SpecData `json:"data"`
}

type SpecData struct {
	ID string `json:"id"`
}

type JobForm struct {
	TOML string `json:"toml"`
}

type Job struct {
	Data JobData `json:"data"`
}

type JobData struct {
	ID string `json:"id"`
}
