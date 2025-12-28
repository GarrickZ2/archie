package status

import (
	"sort"
	"strings"
)

// FeatureMatch 表示一个匹配的 feature 及其相似度得分
type FeatureMatch struct {
	Name     string
	Distance int
}

// LevenshteinDistance 计算两个字符串之间的编辑距离
func LevenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if s1 == s2 {
		return 0
	}

	// 创建一个二维数组来存储编辑距离
	len1, len2 := len(s1), len(s2)
	dp := make([][]int, len1+1)
	for i := range dp {
		dp[i] = make([]int, len2+1)
	}

	// 初始化第一行和第一列
	for i := 0; i <= len1; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		dp[0][j] = j
	}

	// 动态规划计算编辑距离
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			dp[i][j] = min(
				dp[i-1][j]+1,      // 删除
				dp[i][j-1]+1,      // 插入
				dp[i-1][j-1]+cost, // 替换
			)
		}
	}

	return dp[len1][len2]
}

// FindSimilarFeatures 找到与输入相似的 features
// maxDistance 是允许的最大编辑距离
func FindSimilarFeatures(input string, features []Feature, maxDistance int) []FeatureMatch {
	var matches []FeatureMatch

	for _, feature := range features {
		distance := LevenshteinDistance(input, feature.Name)
		if distance <= maxDistance {
			matches = append(matches, FeatureMatch{
				Name:     feature.Name,
				Distance: distance,
			})
		}
	}

	// 按照编辑距离排序，距离越小越靠前
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Distance < matches[j].Distance
	})

	return matches
}

// min 返回三个整数中的最小值
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
