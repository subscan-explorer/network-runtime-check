package rumtime

import (
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const HostUrl = "https://%s.api.subscan.io/api/scan/"

var HttpCli *http.Client
var APIKey string

func init() {
	HttpCli = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Second * 10,
				KeepAlive: time.Second * 30,
			}).DialContext,
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     time.Second * 90,
			TLSHandshakeTimeout: time.Second * 10,
			ForceAttemptHTTP2:   true,
		},
	}
	APIKey = strings.TrimSpace(os.Getenv("SUBSCAN_API_KEY"))
}
