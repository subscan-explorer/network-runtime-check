package output

import (
	"fmt"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type Stdout struct {
}

func NewStdout() *Stdout {
	return new(Stdout)
}

func (Stdout) Output(pallet []string, list []subscan.NetworkPallet) error {
	var table = make([][]string, len(pallet))
	palletIdx := make(map[string]int)
	for i, p := range pallet {
		table[i] = make([]string, len(list)+1)
		table[i][0] = p
		palletIdx[strings.ToLower(p)] = i
	}
	var network []string
	for i, np := range list {
		i = i + 1
		network = append(network, np.Network)
		if np.Err != nil {
			table[0][i] = utils.ErrorReduction(np.Err)
			continue
		}
		for _, p := range np.Pallet {
			if idx, ok := palletIdx[strings.ToLower(p)]; ok {
				table[idx][i] = Exist
			}
		}
		for j := 0; j < len(pallet); j++ {
			if table[j][i] == "" {
				table[j][i] = NotExist
			}
		}
	}
	if len(table) == 0 {
		return nil
	}

	// not enough width for output
	width := utils.TerminalWidth() - utils.MaxLenArrString(pallet) - 8
	maxListLen := utils.MaxLenArrString(network)

	if width < 0 || maxListLen+7 > width {
		tb := gotabulate.Create(table)
		tb.SetAlign("center")
		tb.SetEmptyString("None")
		tb.SetHeaders(network)
		fmt.Println(tb.Render("grid"))
		return nil
	}
	// Adaptive window size
	var idx = 0
	var lastIdx = 0
	for idx < len(network) {
		sub := width
		for idx < len(network) {
			if len(network[idx]) > sub {
				break
			}
			sub = sub - len(network[idx]) - 8
			idx++
		}
		var subTable [][]string
		for _, t := range table {
			var result []string
			result = append(result, t[0])
			result = append(result, t[lastIdx+1:idx+1]...)
			subTable = append(subTable, result)
		}
		tb := gotabulate.Create(subTable)
		tb.SetAlign("center")
		tb.SetEmptyString("None")
		tb.SetHeaders(network[lastIdx:idx])
		tb.SetMaxCellSize(maxListLen)
		lastIdx = idx
		fmt.Println(tb.Render("grid"))
	}
	return nil
}
