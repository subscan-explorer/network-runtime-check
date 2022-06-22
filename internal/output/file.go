package output

import (
	"os"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type FileOutput struct {
	path string
}

func NewFileOutput(path string) *FileOutput {
	return &FileOutput{path: path}
}

func (f FileOutput) Output(pallet []string, list []subscan.NetworkPallet) error {
	table := make([][]string, len(pallet))
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
	tb := gotabulate.Create(table)
	tb.SetHeaders(network)
	tb.SetAlign("center")
	tb.SetEmptyString("None")
	return os.WriteFile(f.path, []byte(tb.Render("grid")), os.FileMode(0644))
}
