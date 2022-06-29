package output

import (
	"fmt"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

type Stdout struct {
	FormatCharterBase
}

func NewStdout() *Stdout {
	return new(Stdout)
}

func (s Stdout) FormatCompareChart(pallet []string, list []subscan.NetworkPallet) error {
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
	totalWidth := utils.TerminalWidth()
	networkMaxWidth := utils.MaxLenArrString(network)
	width := totalWidth - networkMaxWidth - 14
	if data := s.formatChartErrData(list); len(data) != 0 {
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

func (s Stdout) FormatChart(pallet []string, list []subscan.NetworkPallet) error {
	totalWidth := utils.TerminalWidth()
	networkMaxWidth := s.networkMaxLen(list)
	width := totalWidth - networkMaxWidth - 14
	if data := s.formatChartErrData(list); len(data) != 0 {
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
	if data := s.formatChartData(pallet, list); len(data) != 0 {
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
