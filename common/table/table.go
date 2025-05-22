package table

import (
	"strings"

	"github.com/go-xuan/quanx/base/fmtx"
	"github.com/go-xuan/quanx/types/stringx"
)

var (
	headerColor = fmtx.Magenta
	valueColor  = fmtx.Green
)

// SetHeaderColor 设置表头颜色
func SetHeaderColor(color fmtx.Color) {
	headerColor = color
}

// SetValueColor 设置值颜色
func SetValueColor(color fmtx.Color) {
	valueColor = color
}

// Printer 打印器
type Printer interface {
	GetHeaders() []string                 // 获取表头
	GetValuesAndWides() ([]string, []int) // 获取值和最大宽度
}

// Print 打印
func Print[P Printer](data []P) {
	if len(data) == 0 {
		return
	}
	headers := data[0].GetHeaders()
	var maxWides []int
	for _, header := range headers {
		maxWides = append(maxWides, len(header))
	}
	var rows [][]string
	rows = append(rows, headers)
	for _, datum := range data {
		values, wides := datum.GetValuesAndWides()
		rows = append(rows, values)
		maxWides = getMaxWides(maxWides, wides)
	}
	for i, values := range rows {
		sb := strings.Builder{}
		for j, value := range values {
			sb.WriteString(stringx.Fill(value, " ", maxWides[j]+2))
		}
		if i == 0 {
			headerColor.Println(sb.String())
		} else {
			valueColor.Println(sb.String())
		}
	}
}

// 获取最大宽度集（取两个数据集每个位置的最大值）
func getMaxWides(wides1, wides2 []int) []int {
	result := make([]int, len(wides2))
	for i := 0; i < len(wides2); i++ {
		if wides1[i] > wides2[i] {
			result[i] = wides1[i]
		} else {
			result[i] = wides2[i]
		}
	}
	return result
}
