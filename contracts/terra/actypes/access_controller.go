package actypes

import "github.com/smartcontractkit/terra.go/msg"

type InstantiateMsg struct{}

type ExecuteAddAccessMsg struct {
	AddAccess ExecuteAddAccessTypeMsg `json:"add_access"`
}

type ExecuteAddAccessTypeMsg struct {
	Address msg.AccAddress `json:"address"`
}

type ExecuteRemoveAccessMsg struct {
	RemoveAccess ExecuteRemoveAccessTypeMsg `json:"remove_access"`
}

type ExecuteRemoveAccessTypeMsg struct {
	Address msg.AccAddress `json:"address"`
}

// Queries

type QueryHasAccessMsg struct {
	HasAccess QueryHasAccessTypeMsg `json:"has_access"`
}

type QueryHasAccessTypeMsg struct {
	Address msg.AccAddress `json:"address"`
}

type QueryHasAccessResponse struct {
	QueryResult bool `json:"query_result"`
}
