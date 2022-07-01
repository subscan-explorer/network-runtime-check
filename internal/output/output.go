package output

import (
	"strings"

	"github.com/subscan-explorer/network-runtime-check/internal/api/subscan"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

const (
	Exist    = "O"
	NotExist = "X"
)

type FormatCompareCharter interface {
	FormatCompareChart([]string, []subscan.NetworkPallet) error
}

type FormatCharter interface {
	FormatChart([]string, []subscan.NetworkPallet) error
}

type FormatCharterBase struct{}

func (FormatCharterBase) formatChartData(pallet []string, list []subscan.NetworkPallet, maxWidth int) [][]string {
	var palletSet map[string]struct{}
	if len(pallet) != 0 {
		palletSet = make(map[string]struct{})
		for _, p := range pallet {
			palletSet[strings.ToLower(p)] = struct{}{}
		}
	}
	var tableData [][]string
	for _, np := range list {
		if np.Err != nil {
			continue
		}
		var support, resultPallet []string
		support = append(support, np.Network)
		if palletSet != nil {
			for _, p := range np.Pallet {
				if _, ok := palletSet[strings.ToLower(p)]; ok {
					resultPallet = append(resultPallet, p)
				}
			}
		} else {
			resultPallet = np.Pallet
		}
		str := strings.Builder{}
		remainWidth := maxWidth
		for i := 0; i < len(resultPallet); {
			sl := len(resultPallet[i])
			if sl <= remainWidth {
				str.WriteString(resultPallet[i])
				remainWidth -= sl
				if remainWidth >= 2 {
					str.WriteString("  ")
					remainWidth -= 2
				} else {
					str.WriteString(strings.Repeat(" ", remainWidth))
					remainWidth = maxWidth // reset
				}
				i++
			} else {
				if sl >= maxWidth { // too max
					str.Reset()
					break
				}
				str.WriteString(strings.Repeat(" ", remainWidth))
				remainWidth = maxWidth // reset
			}
		}
		if str.Len() == 0 {
			support = append(support, strings.Join(resultPallet, "  "))
		} else {
			support = append(support, str.String())
		}
		tableData = append(tableData, support)
	}
	return tableData
}

func (FormatCharterBase) formatChartErrData(list []subscan.NetworkPallet) [][]string {
	var tableData [][]string
	for _, p := range list {
		if p.Err == nil {
			continue
		}
		tableData = append(tableData, []string{p.Network, utils.ErrorReduction(p.Err)})
	}
	return tableData
}

func (FormatCharterBase) networkMaxLen(list []subscan.NetworkPallet) (maxLen int) {
	for _, pallet := range list {
		l := len(pallet.Network)
		if l > maxLen {
			maxLen = l
		}
	}
	return
}
