package status

import (
	"fmt"
	"sort"
)

// Summary 状态汇总信息
type Summary struct {
	TotalFeatures     int
	StatusCounts      map[FeatureStatus]int
	BlockedFeatures   []Feature
	InProgressCount   int
	CompletedCount    int
	NotStartedCount   int
	OverallProgress   int
	FeaturesByStatus  map[FeatureStatus][]Feature
	StaleFeatures     []Feature // 超过30天未更新的 features
}

// Aggregator 状态聚合器
type Aggregator struct {
	features []Feature
}

// NewAggregator 创建聚合器
func NewAggregator(features []Feature) *Aggregator {
	return &Aggregator{features: features}
}

// Aggregate 聚合统计信息
func (a *Aggregator) Aggregate() *Summary {
	summary := &Summary{
		TotalFeatures:    len(a.features),
		StatusCounts:     make(map[FeatureStatus]int),
		BlockedFeatures:  []Feature{},
		FeaturesByStatus: make(map[FeatureStatus][]Feature),
		StaleFeatures:    []Feature{},
	}

	if len(a.features) == 0 {
		return summary
	}

	// 初始化所有状态计数
	for _, status := range AllStatuses {
		summary.StatusCounts[status] = 0
		summary.FeaturesByStatus[status] = []Feature{}
	}

	// 统计每个状态的数量
	totalProgress := 0
	for _, feature := range a.features {
		status := feature.Status
		summary.StatusCounts[status]++
		summary.FeaturesByStatus[status] = append(summary.FeaturesByStatus[status], feature)

		// 分类统计
		switch status {
		case StatusBlocked:
			summary.BlockedFeatures = append(summary.BlockedFeatures, feature)
		case StatusFinished:
			summary.CompletedCount++
		case StatusNotReviewed, StatusUnknown:
			summary.NotStartedCount++
		default:
			summary.InProgressCount++
		}

		// 累加进度
		totalProgress += GetStatusProgress(status)

		// 检查是否过期（超过30天未更新）
		if feature.IsOld(30) {
			summary.StaleFeatures = append(summary.StaleFeatures, feature)
		}
	}

	// 计算总体进度
	if len(a.features) > 0 {
		summary.OverallProgress = totalProgress / len(a.features)
	}

	return summary
}

// GetTopInsights 获取关键洞察
func (s *Summary) GetTopInsights() []string {
	insights := []string{}

	// 1. 完成度洞察
	if s.TotalFeatures > 0 {
		completionRate := (s.CompletedCount * 100) / s.TotalFeatures
		if completionRate >= 80 {
			insights = append(insights, fmt.Sprintf("🎉 Great progress! %d%% of features are completed", completionRate))
		} else if completionRate >= 50 {
			insights = append(insights, fmt.Sprintf("📊 Halfway there! %d%% of features are completed", completionRate))
		} else if completionRate > 0 {
			insights = append(insights, fmt.Sprintf("🚀 Getting started: %d%% of features are completed", completionRate))
		}
	}

	// 2. 阻塞项警告
	if len(s.BlockedFeatures) > 0 {
		insights = append(insights, fmt.Sprintf("⚠️  %d feature(s) are BLOCKED and need attention", len(s.BlockedFeatures)))
	}

	// 3. 进行中的工作
	if s.InProgressCount > 0 {
		insights = append(insights, fmt.Sprintf("⚡ %d feature(s) are actively in progress", s.InProgressCount))
	}

	// 4. 待开始的工作
	if s.NotStartedCount > 0 {
		insights = append(insights, fmt.Sprintf("📝 %d feature(s) are waiting to be reviewed", s.NotStartedCount))
	}

	// 5. 过期 features
	if len(s.StaleFeatures) > 0 {
		insights = append(insights, fmt.Sprintf("⏰ %d feature(s) haven't been updated in 30+ days", len(s.StaleFeatures)))
	}

	// 6. 设计阶段提醒
	designPhaseCount := s.StatusCounts[StatusUnderDesign] + s.StatusCounts[StatusDesigned] + s.StatusCounts[StatusSpecReady]
	if designPhaseCount > s.TotalFeatures/2 {
		insights = append(insights, "💡 Most features are in design phase - good planning!")
	}

	return insights
}

// GetPhaseDistribution 获取各阶段的分布
func (s *Summary) GetPhaseDistribution() map[string]int {
	return map[string]int{
		"Not Started":    s.NotStartedCount,
		"In Progress":    s.InProgressCount,
		"Completed":      s.CompletedCount,
		"Blocked":        len(s.BlockedFeatures),
	}
}

// GetMostCommonStatus 获取最常见的状态
func (s *Summary) GetMostCommonStatus() FeatureStatus {
	maxCount := 0
	var mostCommon FeatureStatus

	for status, count := range s.StatusCounts {
		if count > maxCount && status != StatusUnknown {
			maxCount = count
			mostCommon = status
		}
	}

	return mostCommon
}

// SortFeaturesByStatus 按状态对 features 排序
func SortFeaturesByStatus(features []Feature) []Feature {
	sorted := make([]Feature, len(features))
	copy(sorted, features)

	statusOrder := make(map[FeatureStatus]int)
	for i, status := range AllStatuses {
		statusOrder[status] = i
	}

	sort.Slice(sorted, func(i, j int) bool {
		orderI := statusOrder[sorted[i].Status]
		orderJ := statusOrder[sorted[j].Status]
		if orderI != orderJ {
			return orderI < orderJ
		}
		return sorted[i].Name < sorted[j].Name
	})

	return sorted
}
