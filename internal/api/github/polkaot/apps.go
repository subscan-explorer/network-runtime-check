package polkaot

import (
	"context"
	"log"
	"strings"

	"github.com/subscan-explorer/network-runtime-check/internal/api/github"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type NetworkNode struct {
	Name      string
	Endpoints []string
}

func NetworkList(ctx context.Context) map[string][]string {
	var (
		networkNode []NetworkNode
		networkMap  = make(map[string][]string)
		err         error
	)

	if networkNode, err = LiveNetworkList(ctx); err != nil {
		log.Println(err.Error())
	}
	for _, node := range networkNode {
		networkMap[node.Name] = node.Endpoints
	}
	if networkNode, err = KusamaNetworkList(ctx); err != nil {
		log.Println(err.Error())
	}
	for _, node := range networkNode {
		if v, ok := networkMap[node.Name]; !ok {
			networkMap[node.Name] = node.Endpoints
		} else {
			v = append(v, node.Endpoints...)
			networkMap[node.Name] = v
		}
	}
	if networkNode, err = PolkadotNetworkList(ctx); err != nil {
		log.Println(err.Error())
	}
	for _, node := range networkNode {
		if v, ok := networkMap[node.Name]; !ok {
			networkMap[node.Name] = node.Endpoints
		} else {
			v = append(v, node.Endpoints...)
			networkMap[node.Name] = v
		}
	}
	return networkMap
}

func KusamaNetworkList(ctx context.Context) ([]NetworkNode, error) {
	const addr = "https://api.github.com/repos/polkadot-js/apps/contents/packages/apps-config/src/endpoints/productionRelayKusama.ts?ref=master"
	return networkList(ctx, addr)
}

func PolkadotNetworkList(ctx context.Context) ([]NetworkNode, error) {
	const addr = "https://api.github.com/repos/polkadot-js/apps/contents/packages/apps-config/src/endpoints/productionRelayPolkadot.ts?ref=master"
	return networkList(ctx, addr)
}

func LiveNetworkList(ctx context.Context) ([]NetworkNode, error) {
	const addr = "https://api.github.com/repos/polkadot-js/apps/contents/packages/apps-config/src/endpoints/production.ts?ref=master"
	return networkList(ctx, addr)
}

func networkList(ctx context.Context, addr string) ([]NetworkNode, error) {
	var (
		content []byte
		err     error
	)
	if content, err = github.GetFileContent(ctx, addr); err != nil {
		return nil, err
	}
	return parseNetwork(string(content)), nil
}

func parseNetwork(content string) []NetworkNode {
	var (
		nodeList []NetworkNode
		list     = strings.Split(content, "\n")
	)
	for i := 0; i < len(list); i++ {
		str := strings.TrimSpace(list[i])
		if strings.HasPrefix(str, "//") {
			continue
		}
		if strings.HasPrefix(str, "info") {
			m := strings.Split(str, ",")
			if len(m) == 0 || len(m[0]) == 0 {
				continue
			}
			info := strings.Split(str, ":")
			if len(info) != 2 {
				continue
			}
			if nodes := parseProviders(calculateIndex(list[i:])); len(nodes) != 0 {
				nodeList = append(nodeList, NetworkNode{
					Name:      strings.Trim(strings.TrimSpace(strings.Trim(info[1], ",")), "'"),
					Endpoints: nodes,
				})
			}
		}
	}
	return nodeList
}

func calculateIndex(data []string) []string {
	var start, end int
	for i, d := range data {
		str := strings.TrimSpace(d)
		if strings.HasPrefix(str, "providers: {") {
			start = i + 1
			continue
		}
		if str == "}" || str == "}," {
			end = i
			break
		}
	}
	if start >= end {
		return nil
	}
	return data[start:end]
}

/// example data
/// 'Moonbeam Foundation': 'wss://wss.api.moonbeam.network',
/// Blast: 'wss://moonbeam.public.blastapi.io', // something
/// Dwellir: 'wss://moonbeam-rpc.dwellir.com',
/// OnFinality: 'wss://moonbeam.api.onfinality.io/public-ws',
/// Pinknode: 'wss://public-rpc.pinknode.io/moonbeam'

func parseProviders(data []string) (result []string) {
	for i := 0; i < len(data); i++ {
		str := strings.TrimSpace(data[i])
		if strings.HasPrefix(str, "//") {
			continue
		}
		m := strings.Split(str, ",")
		if len(m) == 0 || len(m[0]) == 0 {
			continue
		}
		ws := strings.Split(m[0], ": ")
		if len(ws) != 2 {
			continue
		}
		node := strings.TrimSpace(strings.Trim(ws[1], "'"))
		if strings.HasPrefix(node, "wss://") || strings.HasPrefix(node, "ws://") {
			result = append(result, node)
		}
		utils.Reverse(result)
	}
	return
}
