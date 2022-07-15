package output

import (
	"strings"

	"github.com/subscan-explorer/network-runtime-check/conf"
	"github.com/subscan-explorer/network-runtime-check/internal/model"
	"github.com/subscan-explorer/network-runtime-check/internal/utils"
)

const (
	Exist    = "O"
	NotExist = "X"
)

type FormatCompareCharter interface {
	FormatCompareChart([]string, []model.NetworkData[string]) error
}

type FormatCharter interface {
	FormatChart([]model.NetworkData[string]) error
}

type FormatEventCharter interface {
	FormatEventChart([]model.NetworkData[model.Metadata], []conf.NetworkRule) error
}

type FormatCharterBase struct{}

func (FormatCharterBase) formatChartData(list []model.NetworkData[string], maxWidth int) [][]string {
	var tableData [][]string
	for _, np := range list {
		if np.Err != nil {
			continue
		}
		var support, resultPallet []string
		support = append(support, np.Network)
		resultPallet = np.Data
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

func formatChartErrData[T any](list []model.NetworkData[T]) [][]string {
	var tableData [][]string
	for _, p := range list {
		if p.Err == nil {
			continue
		}
		tableData = append(tableData, []string{p.Network, utils.ErrorReduction(p.Err)})
	}
	return tableData
}

func networkMaxLen[T any](list []model.NetworkData[T]) (maxLen int) {
	for _, pallet := range list {
		l := len(pallet.Network)
		if l > maxLen {
			maxLen = l
		}
	}
	return
}

func (FormatCharterBase) formatExtrinsicChartData(nodes []conf.NetworkRule, list []model.NetworkData[model.Metadata]) [][]string {
	nodeMap := make(map[string]conf.NetworkRule)
	for _, node := range nodes {
		nodeMap[node.Name] = node
	}
	var tableData [][]string
	for _, meta := range list {
		if meta.Err != nil {
			continue
		}
		peMap := make(map[string]map[string][]string)
		if c, ok := nodeMap[meta.Network]; ok {
			for _, p := range c.Pallet {
				m := make(map[string][]string)
				for _, e := range p.Extrinsic {
					m[strings.ToLower(e.Name)] = e.Param
				}
				peMap[strings.ToLower(p.Name)] = m
			}
		} else {
			continue
		}
		for _, me := range meta.Data {
			if palletMap, ok := peMap[strings.ToLower(me.Name)]; ok {
				for _, call := range me.Calls {
					if params, ok := palletMap[strings.ToLower(call.Name)]; ok {
						row := []string{meta.Network, me.Name, call.Name}
						args := make([]string, 0, len(call.Args))
						for _, arg := range call.Args {
							args = append(args, arg.TypeName)
						}
						if utils.SliceEqual(params, args) {
							row = append(row, Exist)
						} else {
							row = append(row, NotExist)
							str := strings.Builder{}
							str.WriteString("Exception: [")
							str.WriteString(strings.Join(params, ","))
							str.WriteString("] ")
							str.WriteString("Actual: [")
							str.WriteString(strings.Join(args, ","))
							str.WriteString("]")
							row = append(row, str.String())
						}
						delete(palletMap, strings.ToLower(call.Name))
						tableData = append(tableData, row)
					}
				}
			}
		}
		for p, m := range peMap {
			for e := range m {
				row := []string{meta.Network, p, e, NotExist, "Not Found"}
				tableData = append(tableData, row)
			}
		}
	}
	return tableData
}

func (FormatCharterBase) formatEventChartData(nodes []conf.NetworkRule, list []model.NetworkData[model.Metadata]) [][]string {
	nodeMap := make(map[string]conf.NetworkRule)
	for _, node := range nodes {
		nodeMap[node.Name] = node
	}
	var tableData [][]string
	for _, meta := range list {
		if meta.Err != nil {
			continue
		}
		peMap := make(map[string]map[string][]string)
		if c, ok := nodeMap[meta.Network]; ok {
			for _, p := range c.Pallet {
				m := make(map[string][]string)
				for _, e := range p.Event {
					m[strings.ToLower(e.Name)] = e.Param
				}
				peMap[strings.ToLower(p.Name)] = m
			}
		} else {
			continue
		}

		for _, me := range meta.Data {
			if palletMap, ok := peMap[strings.ToLower(me.Name)]; ok {
				for _, e := range me.Events {
					if params, ok := palletMap[strings.ToLower(e.Name)]; ok {
						var row = []string{meta.Network, me.Name, e.Name}
						var args []string
						if len(e.ArgsTypeName) != 0 {
							args = e.ArgsTypeName
						} else if len(e.Args) != 0 {
							args = e.Args
						}
						if utils.SliceEqual(params, args) {
							row = append(row, Exist)
						} else {
							row = append(row, NotExist)
							str := strings.Builder{}
							str.WriteString("Exception: [")
							str.WriteString(strings.Join(params, ","))
							str.WriteString("] ")
							str.WriteString("Actual: [")
							str.WriteString(strings.Join(args, ","))
							str.WriteString("]")
							row = append(row, str.String())
						}
						delete(palletMap, strings.ToLower(e.Name))
						tableData = append(tableData, row)
					}
				}
			}
		}

		for p, m := range peMap {
			for e := range m {
				row := []string{meta.Network, p, e, NotExist, "Not Found"}
				tableData = append(tableData, row)
			}
		}
	}
	return tableData
}
