package Progress

import (
	"fmt"
	"strings"
)

// ProgressBar 结构体定义进度条属性
type ProgressBar struct {
	Total       int    // 总进度
	Current     int    // 当前进度
	BarLength   int    // 进度条长度
	DisplayChar string // 进度条显示的字符
}

// NewProgressBar 创建一个新的进度条实例
func NewProgressBar(total int, barLength int) *ProgressBar {
	return &ProgressBar{
		Total:       total,
		BarLength:   barLength,
		DisplayChar: "=",
	}
}

// Show 更新并显示进度条
func (pb *ProgressBar) Show() {
	percent := float64(pb.Current) / float64(pb.Total)
	filledLength := int(percent * float64(pb.BarLength))

	// 生成进度条字符串
	bar := pb.getBar(filledLength)

	// 光标回到行首并打印进度条
	fmt.Printf("\r[%-50s] %6.2f%% Completed", bar, percent*100)
}

// getBar 根据完成的长度获得进度条状态
func (pb *ProgressBar) getBar(filledLength int) string {
	// 如果有填充则构造进度条，用'>'表示当前位置
	if filledLength > 0 {
		return strings.Repeat(pb.DisplayChar, filledLength-1) + ">" + strings.Repeat(" ", pb.BarLength-filledLength)
	}
	// 没有填充直接返回空白进度条
	return strings.Repeat(" ", pb.BarLength)
}
