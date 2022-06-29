package substrate

import (
	"context"
	"errors"
	"strings"

	"github.com/subscan-explorer/network-runtime-check/internal/api/github"
)

func PalletList(ctx context.Context) ([]string, error) {
	const palletURL = "https://api.github.com/repos/paritytech/substrate/contents/bin/node/runtime/src/lib.rs?ref=master"
	var (
		content []byte
		err     error
	)

	if content, err = github.GetFileContent(ctx, palletURL); err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, errors.New("failed to get substrate pallet")
	}
	list := strings.Split(string(content), "\n")
	for i := 0; i < len(list); i++ {
		if strings.HasPrefix(list[i], "construct_runtime!(") {
			list = list[i+1:]
			break
		}
	}

	for i := 0; i < len(list); i++ {
		if strings.TrimSpace(list[i]) == "{" {
			list = list[i+1:]
			break
		}
	}
	var result []string
	for _, str := range list {
		s := strings.TrimSpace(str)
		if strings.HasPrefix(s, "//") {
			continue
		}
		if s == "}" {
			break
		}
		p := strings.Split(s, ":")
		if len(p) > 1 {
			result = append(result, p[0])
		}
	}
	return result, nil
}
