package self

import (
	"context"

	"github.com/subscan-explorer/network-runtime-check/internal/api/github"
)

const configURL = "https://api.github.com/repos/subscan-explorer/network-runtime-check/contents/conf/config.yaml?ref=main"

func GetConfigData(ctx context.Context) ([]byte, error) {
	return github.GetFileContent(ctx, configURL)
}
