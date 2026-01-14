package leak

import (
	"time"

	f "github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var _ PromQuerier = (*FakeQueryClient)(nil)

type FakeQueryClient struct {
	startRespCalled bool
	endRespCalled   bool
	isStartResp     bool
	startResp       *f.PrometheusQueryResponse
	endResp         *f.PrometheusQueryResponse
}

func NewFakeQueryClient() *FakeQueryClient {
	return &FakeQueryClient{}
}

func (qc *FakeQueryClient) SetResponses(sr *f.PrometheusQueryResponse, er *f.PrometheusQueryResponse) {
	qc.isStartResp = true
	qc.startResp = sr
	qc.endResp = er
}

func (qc *FakeQueryClient) Query(query string, timestamp time.Time) (*f.PrometheusQueryResponse, error) {
	if qc.isStartResp {
		qc.isStartResp = false
		return qc.startResp, nil
	}
	qc.isStartResp = true
	return qc.endResp, nil
}

func PromSingleValueResponse(val string) *f.PrometheusQueryResponse {
	return &f.PrometheusQueryResponse{
		Status: "",
		Data: &f.PromQueryResponseData{
			Result: []f.PromQueryResponseResult{
				{
					Metric: map[string]string{},
					Value:  []interface{}{"", val},
				},
			},
		},
	}
}
