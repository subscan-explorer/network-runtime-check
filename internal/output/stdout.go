package output

import (
	"fmt"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type Stdout struct {
	FormatCharterBase
}

func NewStdout() *Stdout {
	return new(Stdout)
}

func (s Stdout) FormatCompareChart(pallet []string, list []model.NetworkData[string]) error {
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
	totalWidth := utils.TerminalWidth()
	networkMaxWidth := utils.MaxLenArrString(network)
	width := totalWidth - networkMaxWidth - 14
	if data := formatChartErrData(list); len(data) != 0 {
		fmt.Println()
		fmt.Println("Error list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(width)
		fmt.Println(tb.Render("grid"))
	}
	fmt.Println()

	if len(network) == 0 {
		return nil
	}
	fmt.Println("Result list:")
	// not enough width for output
	width = totalWidth - utils.MaxLenArrString(pallet) - 8
	if width < 0 || networkMaxWidth+6 > width {
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
		tb.SetMaxCellSize(networkMaxWidth)
		lastIdx = idx
		fmt.Println(tb.Render("grid"))
	}
	return nil
}

func (s Stdout) FormatChart(list []model.NetworkData[string]) error {
	totalWidth := utils.TerminalWidth()
	networkMaxWidth := networkMaxLen(list)
	width := totalWidth - networkMaxWidth - 14
	if data := formatChartErrData(list); len(data) != 0 {
		fmt.Println()
		fmt.Println("Error list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(width)
		fmt.Println(tb.Render("grid"))
	}
	fmt.Println()
	if data := s.formatChartData(list, width); len(data) != 0 {
		fmt.Println("Result list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet"})
		tb.SetAlign("left")
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(width)
		fmt.Println(tb.Render("grid"))
	}
	return nil
}

func (s Stdout) FormatEventChart(list []model.NetworkData[model.Metadata], c []conf.ParamRule) error {
	totalWidth := utils.TerminalWidth()
	networkMaxWidth := networkMaxLen(list)
	width := totalWidth - networkMaxWidth - 14

	if data := formatChartErrData(list); len(data) != 0 {
		fmt.Println()
		fmt.Println("Error list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Error reason"})
		tb.SetAlign("left")
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(width)
		fmt.Println(tb.Render("grid"))
	}

	fmt.Println()
	if data := s.formatEventChartData(c, list); len(data) != 0 {
		fmt.Println("Event list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet", "Event", "Check", "Note"})
		tb.SetAlign("left")
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(width / 5)
		fmt.Println(tb.Render("grid"))
	}

	fmt.Println()
	if data := s.formatExtrinsicChartData(c, list); len(data) != 0 {
		fmt.Println("Extrinsic list:")
		tb := gotabulate.Create(data)
		tb.SetHeaders([]string{"Network", "Pallet", "Extrinsic", "Check", "Note"})
		tb.SetAlign("left")
		tb.SetWrapStrings(true)
		tb.SetMaxCellSize(width / 5)
		fmt.Println(tb.Render("grid"))
	}
	return nil
}
