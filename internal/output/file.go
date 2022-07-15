package output

import (
	"os"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
)

type FileOutput struct {
	FormatCharterBase
	path string
}

func NewFileOutput(path string) *FileOutput {
	return &FileOutput{path: path}
}

func (f FileOutput) FormatCompareChart(pallet []string, list []model.NetworkData[string]) error {
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
		for _, p := range np.Data {
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

	if data := formatChartErrData(list); len(data) != 0 {
		_, _ = fd.WriteString("Error list:\n")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	_, _ = fd.Write([]byte{'\n'})
	_, _ = fd.WriteString("Result list:\n")
	tb := gotabulate.Create(table)
	tb.SetHeaders(network)
	tb.SetAlign("center")
	_, _ = fd.WriteString(tb.Render("grid"))
	return fd.Sync()
}

func (f FileOutput) FormatChart(list []model.NetworkData[string]) error {
	var fd, err = os.Create(f.path)
	if err != nil {
		return err
	}
	defer fd.Close()
	if data := formatChartErrData(list); len(data) != 0 {
		_, _ = fd.WriteString("Error list:\n")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	_, _ = fd.Write([]byte{'\n'})
	const maxWidth = 180
	if data := f.formatChartData(list, maxWidth); len(data) != 0 {
		_, _ = fd.WriteString("Result list:\n")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet"})
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(maxWidth)
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	return fd.Sync()
}

func (f FileOutput) FormatEventChart(list []model.NetworkData[model.Metadata], c []conf.ParamRule) error {
	var fd, err = os.Create(f.path)
	if err != nil {
		return err
	}
	defer fd.Close()

	if data := formatChartErrData(list); len(data) != 0 {
		_, _ = fd.WriteString("Error list:\n")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	_, _ = fd.Write([]byte{'\n'})

	if data := f.formatEventChartData(c, list); len(data) != 0 {
		_, _ = fd.WriteString("Event list:\n")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet", "Event", "Check", "Note"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}

	_, _ = fd.Write([]byte{'\n'})

	if data := f.formatExtrinsicChartData(c, list); len(data) != 0 {
		_, _ = fd.WriteString("Extrinsic list:\n")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet", "Extrinsic", "Check", "Note"})
		tb.SetAlign("left")
		_, _ = fd.WriteString(tb.Render("grid"))
	}
	return fd.Sync()
}
