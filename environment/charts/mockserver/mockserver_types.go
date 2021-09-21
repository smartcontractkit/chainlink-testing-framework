package mockserver

type HttpRequest struct {
	Path string `json:"path"`
}

type HttpResponse struct {
	Body interface{} `json:"body"`
}

type HttpInitializer struct {
	Request  HttpRequest  `json:"httpRequest"`
	Response HttpResponse `json:"httpResponse"`
}

// For OTPE - weiwatchers

type NodeInfoJSON struct {
	ID          string   `json:"id"`
	NodeAddress []string `json:"nodeAddress"`
}

type ContractInfoJSON struct {
	ContractAddress string `json:"contractAddress"`
	ContractVersion int    `json:"contractVersion"`
	Path            string `json:"path"`
	Status          string `json:"status"`
}
