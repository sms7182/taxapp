package terminal

import (
	"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"tax-management/pkg"
)

type ClientLoggerExtensionImpl struct {
	GatewayRepository pkg.Repository
}

func (h ClientLoggerExtensionImpl) Do(taxRawId *uint, taxProcessId *uint, requestUniqueId string, request *http.Request, apiname string) (*http.Response, error) {
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
			requestUniqueId,
			apiname,
			request.URL.String(),
			-1,
			requestBody,
			nil,
			&msg,
		)
		return resp, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res := err.Error()
		h.GatewayRepository.LogReqRes(
			taxRawId,
			taxProcessId,
			requestUniqueId,
			apiname,
			request.URL.String(),
			resp.StatusCode,
			requestBody,
			&res,
			&res,
		)
		return resp, err
	}
	resStr := string(responseBody)
	if resp.StatusCode == http.StatusOK {
		_ = h.GatewayRepository.LogReqRes(
			taxRawId,
			taxProcessId,
			requestUniqueId,
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
			requestUniqueId,
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