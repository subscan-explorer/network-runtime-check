package output

import "github.com/subscan-explorer/network-runtime-check/internal/api/subscan"

const (
	Exist    = "O"
	NotExist = "X"
)

type Output interface {
	Output([]string, []subscan.NetworkPallet) error
}
