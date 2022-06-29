package output

import (
	"os"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
)

type FileOutput struct {
	FormatCharterBase
	path string
}

func NewFileOutput(path string) *FileOutput {
	return &FileOutput{path: path}
}

func (f FileOutput) FormatCompareChart(pallet []string, list []subscan.NetworkPallet) error {
	var fd, err = os.Create(f.path)
	if err != nil {
		return err
	}
	defer fd.Close()
	table := make([][]string, len(pallet))
	palletIdx := make(map[string]int, len(pallet))
	network := make([]string, 0, len(list))

	for i, p := range pallet {
		table[i] = make([]string, len(list)+1)
		table[i][0] = p
		palletIdx[strings.ToLower(p)] = i
	}
	i := 0
	for _, np := range list {
		if np.Err != nil {
			continue
		}
		i++
		network = append(network, np.Network)
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

	if data := f.formatChartErrData(list); len(data) != 0 {
		_, _ = fd.WriteString("Error list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	_, _ = fd.Write([]byte{'\n'})
	_, _ = fd.WriteString("Result list:")
	tb := gotabulate.Create(table)
	tb.SetHeaders(network)
	tb.SetAlign("center")
	_, _ = fd.WriteString(tb.Render("grid"))
	return fd.Sync()
}

func (f FileOutput) FormatChart(pallet []string, list []subscan.NetworkPallet) error {
	var fd, err = os.Create(f.path)
	if err != nil {
		return err
	}
	defer fd.Close()
	if data := f.formatChartErrData(list); len(data) != 0 {
		_, _ = fd.WriteString("Error list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	_, _ = fd.Write([]byte{'\n'})
	if data := f.formatChartData(pallet, list); len(data) != 0 {
		_, _ = fd.WriteString("Result list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	return fd.Sync()
}
