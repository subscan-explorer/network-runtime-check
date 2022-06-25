package self

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/subscan-explorer/network-runtime-check/internal/api"
)

type configResp struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

const configURL = "https://api.github.com/repos/subscan-explorer/network-runtime-check/contents/conf/config.yaml?ref=master"

//const configURL = "https://api.github.com/repos/subscan-explorer/network-runtime-check/contents/conf/config.yaml?ref=feat/pallet-compare"

func GetConfigData(ctx context.Context) ([]byte, error) {
	var (
		req     *http.Request
		rsp     *http.Response
		rspData configResp
		err     error
	)
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, configURL, nil); err != nil {
		return nil, err
	}
	if rsp, err = api.HTTPCli.Do(req); err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		_, _ = io.Copy(ioutil.Discard, rsp.Body)
		return nil, errors.New(rsp.Status)
	}
	if err = json.NewDecoder(rsp.Body).Decode(&rspData); err != nil {
		return nil, err
	}
	var content []byte
	if content, err = base64.StdEncoding.DecodeString(rspData.Content); err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("failed to get substrate pallet")
	}
	return content, nil
}
