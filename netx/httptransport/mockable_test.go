package httptransport

import (
	"io/ioutil"
	"net/http"
)

type MockableTransport struct {
	Err  error
	Resp *http.Response
}

func (txp MockableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		ioutil.ReadAll(req.Body)
		req.Body.Close()
	}
	if txp.Err != nil {
		return nil, txp.Err
	}
	txp.Resp.Request = req // non thread safe but it doesn't matter
	return txp.Resp, nil
}

func (txp MockableTransport) CloseIdleConnections() {}