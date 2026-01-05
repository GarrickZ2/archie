package status

import (
	"fmt"
	"sort"
	"strings"
)

// DependencyGraph 表示 feature 之间的依赖关系图
type DependencyGraph struct {
	Features         []Feature
	FeaturesByKey    map[string]*Feature
	DependsOn        map[string][]string // feature-key -> list of dependencies
	DependedBy       map[string][]string // feature-key -> list of dependents (reverse)
	NoDependencies   []string            // features with no dependencies
	HasDependencies  []string            // features with dependencies
	CircularDeps     [][]string          // circular dependency chains (if any)
}

// BuildDependencyGraph 从 features 构建依赖图
func BuildDependencyGraph(features []Feature) *DependencyGraph {
	graph := &DependencyGraph{
		Features:      features,
		FeaturesByKey: make(map[string]*Feature),
		DependsOn:     make(map[string][]string),
		DependedBy:    make(map[string][]string),
	}

	// 构建索引
	for i := range features {
		feature := &features[i]
		graph.FeaturesByKey[feature.Name] = feature
	}

	// 构建依赖关系
	for _, feature := range features {
		deps := []string{}
		for depKey := range feature.Dependencies {
			deps = append(deps, depKey)
		}
		sort.Strings(deps)

		if len(deps) > 0 {
			graph.DependsOn[feature.Name] = deps
			graph.HasDependencies = append(graph.HasDependencies, feature.Name)

			// 构建反向依赖关系
			for _, depKey := range deps {
				graph.DependedBy[depKey] = append(graph.DependedBy[depKey], feature.Name)
			}
		} else {
			graph.NoDependencies = append(graph.NoDependencies, feature.Name)
		}
	}

	// 排序
	sort.Strings(graph.NoDependencies)
	sort.Strings(graph.HasDependencies)

	for key := range graph.DependedBy {
		sort.Strings(graph.DependedBy[key])
	}

	// 检测循环依赖 (简单实现，可以优化)
	graph.CircularDeps = detectCircularDependencies(graph)

	return graph
}

// detectCircularDependencies 检测循环依赖
func detectCircularDependencies(graph *DependencyGraph) [][]string {
	// 简单的 DFS 检测
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	cycles := [][]string{}

	var dfs func(node string, path []string) bool
	dfs = func(node string, path []string) bool {
		visited[node] = true
		recStack[node] = true
		path = append(path, node)

		for _, dep := range graph.DependsOn[node] {
			if !visited[dep] {
				if dfs(dep, path) {
					return true
				}
			} else if recStack[dep] {
				// 发现循环
				cycleStart := 0
				for i, n := range path {
					if n == dep {
						cycleStart = i
						break
					}
				}
				cycle := append([]string{}, path[cycleStart:]...)
				cycles = append(cycles, cycle)
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for _, feature := range graph.Features {
		if !visited[feature.Name] {
			dfs(feature.Name, []string{})
		}
	}

	return cycles
}

// GetTopologicalOrder 获取拓扑排序（设计顺序建议）
func (g *DependencyGraph) GetTopologicalOrder() []string {
	// Kahn 算法
	inDegree := make(map[string]int)
	for _, feature := range g.Features {
		inDegree[feature.Name] = len(g.DependsOn[feature.Name])
	}

	queue := []string{}
	for _, feature := range g.Features {
		if inDegree[feature.Name] == 0 {
			queue = append(queue, feature.Name)
		}
	}

	result := []string{}
	for len(queue) > 0 {
		// 按字母顺序排序以保证确定性输出
		sort.Strings(queue)

		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		for _, dependent := range g.DependedBy[current] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	return result
}

// DependencyGraphDisplay 负责展示依赖图
type DependencyGraphDisplay struct {
	graph *DependencyGraph
}

// NewDependencyGraphDisplay 创建依赖图展示器
func NewDependencyGraphDisplay(graph *DependencyGraph) *DependencyGraphDisplay {
	return &DependencyGraphDisplay{graph: graph}
}

// Show 展示依赖图
func (d *DependencyGraphDisplay) Show() {
	fmt.Println(ColorBold + "  🔗 Feature Dependencies Graph" + ColorReset)
	fmt.Println()

	if len(d.graph.Features) == 0 {
		fmt.Println(ColorDim + "  No features found" + ColorReset)
		return
	}

	// 检查是否有循环依赖
	if len(d.graph.CircularDeps) > 0 {
		d.showCircularDependencies()
		fmt.Println()
	}

	// 显示树状依赖图
	d.showDependencyTree()
	fmt.Println()

	// 显示推荐的设计顺序
	d.showDesignOrder()
}

// showCircularDependencies 显示循环依赖警告
func (d *DependencyGraphDisplay) showCircularDependencies() {
	fmt.Println(ColorRed + ColorBold + "  ⚠️  Circular Dependencies Detected!" + ColorReset)
	fmt.Println()

	for i, cycle := range d.graph.CircularDeps {
		fmt.Printf("  %sCircle %d:%s ", ColorRed, i+1, ColorReset)
		cycleStr := strings.Join(cycle, " → ")
		fmt.Printf("%s → %s\n", cycleStr, cycle[0])
	}
}

// showDependencyTree 显示树状依赖图
func (d *DependencyGraphDisplay) showDependencyTree() {
	// 找出所有根节点（独立的 features）
	roots := []string{}
	for _, featureKey := range d.graph.NoDependencies {
		// 只显示被其他 feature 依赖的根节点
		if len(d.graph.DependedBy[featureKey]) > 0 {
			roots = append(roots, featureKey)
		}
	}

	// 如果没有根节点或者所有 features 都独立，显示简化信息
	if len(roots) == 0 {
		if len(d.graph.HasDependencies) == 0 {
			fmt.Println(ColorDim + "  All features are independent (no dependencies)" + ColorReset)
			return
		}
		// 有依赖但没有内部依赖（所有依赖都指向不存在的 features）
		fmt.Println(ColorDim + "  Dependency tree (→ indicates dependency direction):" + ColorReset)
		fmt.Println()
		for _, featureKey := range d.graph.HasDependencies {
			d.showFeatureNode(featureKey, "", true, make(map[string]bool))
		}
		return
	}

	fmt.Println(ColorDim + "  Dependency tree (→ indicates enables):" + ColorReset)
	fmt.Println()

	// 为每个根节点构建树
	visited := make(map[string]bool)
	for i, root := range roots {
		isLast := i == len(roots)-1
		d.drawTree(root, "", isLast, visited)
	}
}

// drawTree 递归绘制依赖树
func (d *DependencyGraphDisplay) drawTree(featureKey string, prefix string, isLast bool, visited map[string]bool) {
	feature := d.graph.FeaturesByKey[featureKey]
	if feature == nil {
		return
	}

	// 标记已访问（但允许重复显示以展示完整的依赖路径）
	alreadyVisited := visited[featureKey]

	// 绘制当前节点
	connector := "├─"
	if isLast {
		connector = "└─"
	}

	statusColor := GetStatusColor(feature.Status)
	statusName := strings.ReplaceAll(string(feature.Status), "_", " ")

	if prefix == "" {
		// 根节点
		fmt.Printf("  %s%s%s %s[%s]%s",
			ColorBold, featureKey, ColorReset,
			statusColor, statusName, ColorReset)
	} else {
		fmt.Printf("  %s%s %s%s%s %s[%s]%s",
			prefix, connector,
			ColorBold, featureKey, ColorReset,
			statusColor, statusName, ColorReset)
	}

	// 如果这个节点有多个依赖，显示所有依赖
	deps := d.graph.DependsOn[featureKey]
	if len(deps) > 1 {
		// 显示所有依赖
		fmt.Printf(" %s(also depends on: %s)%s",
			ColorDim,
			strings.Join(deps, ", "),
			ColorReset)
	}

	fmt.Println()

	// 如果已经访问过，不再递归展开（避免无限循环）
	if alreadyVisited {
		return
	}
	visited[featureKey] = true

	// 获取依赖当前 feature 的所有 features
	dependents := d.graph.DependedBy[featureKey]
	if len(dependents) == 0 {
		return
	}

	// 准备下一层的前缀
	var newPrefix string
	if prefix == "" {
		newPrefix = "  "
	} else if isLast {
		newPrefix = prefix + "   "
	} else {
		newPrefix = prefix + "│  "
	}

	// 递归绘制子节点
	for i, dep := range dependents {
		isLastChild := i == len(dependents)-1
		d.drawTree(dep, newPrefix, isLastChild, visited)
	}
}

// showFeatureNode 显示单个 feature 节点（用于没有根节点的情况）
func (d *DependencyGraphDisplay) showFeatureNode(featureKey string, prefix string, isLast bool, visited map[string]bool) {
	feature := d.graph.FeaturesByKey[featureKey]
	if feature == nil {
		return
	}

	if visited[featureKey] {
		return
	}
	visited[featureKey] = true

	statusColor := GetStatusColor(feature.Status)
	statusName := strings.ReplaceAll(string(feature.Status), "_", " ")

	// 显示当前节点及其依赖
	fmt.Printf("  %s%s%s %s[%s]%s\n",
		ColorBold, featureKey, ColorReset,
		statusColor, statusName, ColorReset)

	deps := d.graph.DependsOn[featureKey]
	for i, depKey := range deps {
		isLastDep := i == len(deps)-1
		connector := "├─"
		if isLastDep {
			connector = "└─"
		}

		depFeature := d.graph.FeaturesByKey[depKey]
		if depFeature != nil {
			depStatusColor := GetStatusColor(depFeature.Status)
			depStatusName := strings.ReplaceAll(string(depFeature.Status), "_", " ")
			fmt.Printf("    %s→ %s%s%s %s[%s]%s\n",
				connector,
				ColorBold, depKey, ColorReset,
				depStatusColor, depStatusName, ColorReset)
		} else {
			fmt.Printf("    %s→ %s%s%s %s[NOT FOUND]%s\n",
				connector,
				ColorBold, depKey, ColorReset,
				ColorRed, ColorReset)
		}
	}
	fmt.Println()
}

// showDesignOrder 显示推荐的设计顺序
func (d *DependencyGraphDisplay) showDesignOrder() {
	order := d.graph.GetTopologicalOrder()

	if len(order) == 0 {
		return
	}

	fmt.Println()
	fmt.Println(ColorBold + "  📋 Recommended Design Order" + ColorReset)
	fmt.Println(ColorDim + "  (Features ordered by dependencies)" + ColorReset)
	fmt.Println()

	for i, featureKey := range order {
		feature := d.graph.FeaturesByKey[featureKey]
		if feature == nil {
			continue
		}

		statusColor := GetStatusColor(feature.Status)
		statusName := strings.ReplaceAll(string(feature.Status), "_", " ")

		prefix := fmt.Sprintf("%s%2d.%s", ColorDim, i+1, ColorReset)
		fmt.Printf("  %s %-30s %s[%s]%s\n",
			prefix,
			feature.Name,
			statusColor, statusName, ColorReset)
	}
}


// ShowCompact 显示紧凑版本的依赖图
func (d *DependencyGraphDisplay) ShowCompact() {
	if len(d.graph.Features) == 0 {
		return
	}

	fmt.Printf("%s🔗 Dependencies:%s %d with deps | %d independent",
		ColorBold, ColorReset,
		len(d.graph.HasDependencies),
		len(d.graph.NoDependencies))

	if len(d.graph.CircularDeps) > 0 {
		fmt.Printf(" | %s%d circular%s", ColorRed, len(d.graph.CircularDeps), ColorReset)
	}

	fmt.Println()
}
