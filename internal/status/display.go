package status

import (
	"fmt"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
	ColorGray   = "\033[90m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// Display 负责展示状态信息
type Display struct {
	summary *Summary
}

// NewDisplay 创建展示器
func NewDisplay(summary *Summary) *Display {
	return &Display{summary: summary}
}

// Show 展示完整的状态报告
func (d *Display) Show() {
	fmt.Println()
	d.showHeader()
	fmt.Println()
	d.showOverallProgress()
	fmt.Println()
	d.showKeyMetrics()
	fmt.Println()
	d.showStatusDistribution()
	fmt.Println()
	d.showInsights()

	if len(d.summary.BlockedFeatures) > 0 {
		fmt.Println()
		d.showBlockedFeatures()
	}

	if len(d.summary.StaleFeatures) > 0 {
		fmt.Println()
		d.showStaleFeatures()
	}

	fmt.Println()
	d.showDetailedFeatureList()
	fmt.Println()
}

// showHeader 显示标题
func (d *Display) showHeader() {
	title := "📊 Project Status Report"
	separator := strings.Repeat("═", 70)

	fmt.Println(ColorCyan + "  ╔" + separator + "╗" + ColorReset)
	padding := (70 - len(title)) / 2
	fmt.Print(ColorCyan + "  ║ " + ColorReset)
	fmt.Print(strings.Repeat(" ", padding))
	fmt.Print(ColorBold + title + ColorReset)
	fmt.Print(strings.Repeat(" ", 70-len(title)-padding))
	fmt.Println(ColorCyan + " ║" + ColorReset)
	fmt.Println(ColorCyan + "  ╚" + separator + "╝" + ColorReset)
}

// showOverallProgress 显示总体进度
func (d *Display) showOverallProgress() {
	fmt.Println(ColorBold + "  Overall Progress" + ColorReset)
	fmt.Println()

	progress := d.summary.OverallProgress
	barWidth := 50
	filled := (progress * barWidth) / 100

	// 选择进度条颜色
	barColor := ColorGreen
	if progress < 30 {
		barColor = ColorRed
	} else if progress < 70 {
		barColor = ColorYellow
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("  %s%s%s %s%d%%%s\n", barColor, bar, ColorReset, ColorBold, progress, ColorReset)

	// 显示阶段指示器
	fmt.Println()
	d.showPhaseIndicator()
}

// showPhaseIndicator 显示阶段指示器
func (d *Display) showPhaseIndicator() {
	phases := []struct {
		name  string
		count int
		icon  string
	}{
		{"Not Started", d.summary.NotStartedCount, "📝"},
		{"In Progress", d.summary.InProgressCount, "⚡"},
		{"Completed", d.summary.CompletedCount, "✅"},
		{"Blocked", len(d.summary.BlockedFeatures), "🚫"},
	}

	fmt.Print("  ")
	for i, phase := range phases {
		if i > 0 {
			fmt.Print(ColorDim + " → " + ColorReset)
		}

		color := ColorGray
		if phase.count > 0 {
			switch phase.name {
			case "Completed":
				color = ColorGreen
			case "In Progress":
				color = ColorBlue
			case "Blocked":
				color = ColorRed
			case "Not Started":
				color = ColorYellow
			}
		}

		fmt.Printf("%s%s %s (%d)%s", color, phase.icon, phase.name, phase.count, ColorReset)
	}
	fmt.Println()
}

// showKeyMetrics 显示关键指标
func (d *Display) showKeyMetrics() {
	fmt.Println(ColorBold + "  Key Metrics" + ColorReset)
	fmt.Println()

	metrics := []struct {
		label string
		value string
		color string
	}{
		{"Total Features", fmt.Sprintf("%d", d.summary.TotalFeatures), ColorCyan},
		{"Completed", fmt.Sprintf("%d (%.0f%%)", d.summary.CompletedCount, d.getPercentage(d.summary.CompletedCount)), ColorGreen},
		{"In Progress", fmt.Sprintf("%d (%.0f%%)", d.summary.InProgressCount, d.getPercentage(d.summary.InProgressCount)), ColorBlue},
		{"Not Started", fmt.Sprintf("%d (%.0f%%)", d.summary.NotStartedCount, d.getPercentage(d.summary.NotStartedCount)), ColorYellow},
		{"Blocked", fmt.Sprintf("%d (%.0f%%)", len(d.summary.BlockedFeatures), d.getPercentage(len(d.summary.BlockedFeatures))), ColorRed},
	}

	for _, metric := range metrics {
		fmt.Printf("  %s%-15s%s %s%s%s\n", ColorDim, metric.label+":", ColorReset, metric.color, metric.value, ColorReset)
	}
}

// showStatusDistribution 显示状态分布
func (d *Display) showStatusDistribution() {
	fmt.Println(ColorBold + "  Status Distribution" + ColorReset)
	fmt.Println()

	maxCount := 0
	for _, count := range d.summary.StatusCounts {
		if count > maxCount {
			maxCount = count
		}
	}

	if maxCount == 0 {
		fmt.Println(ColorDim + "  No features found" + ColorReset)
		return
	}

	for _, status := range AllStatuses {
		count := d.summary.StatusCounts[status]
		if count == 0 {
			continue
		}

		percentage := d.getPercentage(count)
		barWidth := 30
		filled := 0
		if maxCount > 0 {
			filled = (count * barWidth) / maxCount
		}

		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		color := GetStatusColor(status)

		statusName := strings.ReplaceAll(string(status), "_", " ")
		fmt.Printf("  %s%-18s%s %s%s%s %3d (%.0f%%)\n",
			ColorDim, statusName, ColorReset,
			color, bar, ColorReset,
			count, percentage)
	}
}

// showInsights 显示关键洞察
func (d *Display) showInsights() {
	insights := d.summary.GetTopInsights()
	if len(insights) == 0 {
		return
	}

	fmt.Println(ColorBold + "  💡 Insights" + ColorReset)
	fmt.Println()

	for _, insight := range insights {
		fmt.Println("  " + insight)
	}
}

// showBlockedFeatures 显示阻塞的 features
func (d *Display) showBlockedFeatures() {
	fmt.Println(ColorBold + ColorRed + "  ⚠️  Blocked Features (Need Attention)" + ColorReset)
	fmt.Println()

	for _, feature := range d.summary.BlockedFeatures {
		reason := feature.Reason
		if reason == "" {
			reason = "No reason provided"
		}

		owner := feature.Owner
		if owner == "" {
			owner = "Unassigned"
		}

		fmt.Printf("  %s%-30s%s %s[%s]%s %s%s%s\n",
			ColorBold, feature.Name, ColorReset,
			ColorDim, owner, ColorReset,
			ColorYellow, reason, ColorReset)
	}
}

// showStaleFeatures 显示过期的 features
func (d *Display) showStaleFeatures() {
	fmt.Println(ColorBold + ColorYellow + "  ⏰ Stale Features (Not Updated in 30+ Days)" + ColorReset)
	fmt.Println()

	for _, feature := range d.summary.StaleFeatures {
		statusName := strings.ReplaceAll(string(feature.Status), "_", " ")
		fmt.Printf("  %s%-30s%s %s%-18s%s %sLast: %s%s\n",
			ColorBold, feature.Name, ColorReset,
			ColorDim, statusName, ColorReset,
			ColorDim, feature.LastUpdated, ColorReset)
	}
}

// showDetailedFeatureList 显示详细的 feature 列表（按状态分组）
func (d *Display) showDetailedFeatureList() {
	fmt.Println(ColorBold + "  📋 Features by Status" + ColorReset)
	fmt.Println()

	for _, status := range AllStatuses {
		features := d.summary.FeaturesByStatus[status]
		if len(features) == 0 {
			continue
		}

		statusName := strings.ReplaceAll(string(status), "_", " ")
		color := GetStatusColor(status)

		fmt.Printf("  %s%s%s (%d)\n", color+ColorBold, statusName, ColorReset, len(features))

		for _, feature := range features {
			owner := feature.Owner
			if owner == "" {
				owner = ColorDim + "unassigned" + ColorReset
			}

			updated := feature.LastUpdated
			if updated == "" || updated == "YYYY-MM-DD" {
				updated = ColorDim + "not set" + ColorReset
			}

			fmt.Printf("    • %-30s %s[%s]%s  %sUpdated: %s%s\n",
				feature.Name,
				ColorDim, owner, ColorReset,
				ColorDim, updated, ColorReset)
		}
		fmt.Println()
	}
}

// getPercentage 计算百分比
func (d *Display) getPercentage(count int) float64 {
	if d.summary.TotalFeatures == 0 {
		return 0
	}
	return float64(count*100) / float64(d.summary.TotalFeatures)
}

// ShowCompact 显示紧凑版本的状态报告
func (d *Display) ShowCompact() {
	fmt.Println()
	fmt.Printf("%s📊 Project Status:%s %d features | %d%% complete | %d blocked\n",
		ColorBold, ColorReset,
		d.summary.TotalFeatures,
		d.summary.OverallProgress,
		len(d.summary.BlockedFeatures))

	// 简单的进度条
	progress := d.summary.OverallProgress
	barWidth := 40
	filled := (progress * barWidth) / 100

	barColor := ColorGreen
	if progress < 30 {
		barColor = ColorRed
	} else if progress < 70 {
		barColor = ColorYellow
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	fmt.Printf("%s%s%s %d%%\n", barColor, bar, ColorReset, progress)
	fmt.Println()
}
