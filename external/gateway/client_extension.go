package gateway

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"tax-management/pkg"
)

type ClientLoggerExtensionImpl struct {
	GatewayRepository pkg.ClientRepository
}

func (h ClientLoggerExtensionImpl) Do(taxRawId uint, taxProcessId uint, request *http.Request, apiname string) (*http.Response, error) {
	requestBody := "{}"
	if request.Body != nil {
		reqBody, err := ioutil.ReadAll(request.Body)
		request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
		if err != nil {
			return nil, errors.New("")
		}
		requestBody = string(reqBody)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		msg := err.Error()
		_ = h.GatewayRepository.LogReqRes(
			taxRawId,
			taxProcessId,
			apiname,
			request.URL.String(),
			-1,
			requestBody,
			nil,
			&msg,
		)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(request.Body)
	resStr := string(requestBody)
	if resp.StatusCode == http.StatusOK {
		_ = h.GatewayRepository.LogReqRes(
			taxRawId,
			taxProcessId,
			apiname,
			request.URL.String(),
			resp.StatusCode,
			requestBody,
			&resStr,
			nil,
		)
	} else {
		errMsg := resp.Status
		_ = h.GatewayRepository.LogReqRes(
			taxRawId,
			taxProcessId,
			apiname,
			request.URL.String(),
			resp.StatusCode,
			requestBody,
			&resStr,
			&errMsg,
		)
	}
	_ = resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewReader(responseBody))
	return resp, nil
}
