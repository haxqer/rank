package rank

import (
	"math/rand"
	"time"
)

const (
	// MaxLevel 跳表的最大层级
	MaxLevel = 42
	// Probability 跳表层级提升概率，值为1/4
	Probability = 0.25
)

var defaultRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Element 是跳表中存储的元素
type Element struct {
	// Member 成员ID或名称
	Member string
	// Score 分数，用于排序
	Score int64
	// Data 可以存储的额外数据
	Data interface{}
}

// node 内部节点结构
type node struct {
	element Element
	// level[i] 表示第i层的下一个节点和跨度
	level []*levelNode
}

// levelNode 表示跳表中某一层级的节点
type levelNode struct {
	forward *node  // 指向这一层的下一个节点
	span    uint64 // 到下一个节点的跨度
}

// SkipList 跳表实现
type SkipList struct {
	head       *node            // 头节点，不包含实际数据
	tail       *node            // 尾节点
	length     uint64           // 元素数量
	level      int              // 当前的最大层级
	elementMap map[string]*node // 成员到节点的映射，用于快速查找
}

// NewSkipList 创建一个新的跳表
func NewSkipList() *SkipList {
	head := &node{
		level: make([]*levelNode, MaxLevel),
	}

	for i := 0; i < MaxLevel; i++ {
		head.level[i] = &levelNode{
			forward: nil,
			span:    0,
		}
	}

	return &SkipList{
		head:       head,
		level:      1,
		elementMap: make(map[string]*node),
	}
}

// randomLevel 随机生成层数
func randomLevel() int {
	level := 1
	for level < MaxLevel && defaultRand.Float64() < Probability {
		level++
	}
	return level
}

// Insert 插入元素，如果已存在则更新
func (sl *SkipList) Insert(member string, score int64, data interface{}) *Element {
	// 如果已存在，先删除旧的
	if oldNode, ok := sl.elementMap[member]; ok {
		sl.Delete(member, oldNode.element.Score)
	}

	// 创建新节点
	level := randomLevel()
	if level > sl.level {
		sl.level = level
	}

	newNode := &node{
		element: Element{
			Member: member,
			Score:  score,
			Data:   data,
		},
		level: make([]*levelNode, level),
	}

	for i := 0; i < level; i++ {
		newNode.level[i] = &levelNode{
			forward: nil,
			span:    0,
		}
	}

	// 获取插入位置
	var update [MaxLevel]*node
	var rank [MaxLevel]uint64

	// 查找位置
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		// 注意这里的比较逻辑: 分数高的排在前面，分数相同时按照成员ID字典序排序
		for x.level[i].forward != nil &&
			(x.level[i].forward.element.Score > score ||
				(x.level[i].forward.element.Score == score &&
					x.level[i].forward.element.Member < member)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	// 插入节点
	for i := 0; i < level; i++ {
		newNode.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = newNode

		// 更新跨度
		newNode.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// 更新跨度
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}

	// 如果是尾节点，则更新尾指针
	if newNode.level[0].forward == nil {
		sl.tail = newNode
	}

	// 保存到映射表
	sl.elementMap[member] = newNode
	sl.length++

	return &newNode.element
}

// Delete 删除元素
func (sl *SkipList) Delete(member string, score int64) bool {
	// 查找要删除的节点
	var update [MaxLevel]*node

	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		// 注意这里的比较逻辑: 分数高的排在前面，分数相同时按照成员ID字典序排序
		for x.level[i].forward != nil &&
			(x.level[i].forward.element.Score > score ||
				(x.level[i].forward.element.Score == score &&
					x.level[i].forward.element.Member < member)) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	// 找到待删除节点
	x = x.level[0].forward
	if x != nil && x.element.Score == score && x.element.Member == member {
		// 从各层中删除
		for i := 0; i < sl.level; i++ {
			if update[i].level[i].forward == x {
				update[i].level[i].span += x.level[i].span - 1
				update[i].level[i].forward = x.level[i].forward
			} else {
				update[i].level[i].span--
			}
		}

		// 如果删除的是尾节点
		if x.level[0].forward == nil {
			sl.tail = update[0]
		}

		// 更新最高层级
		for sl.level > 1 && sl.head.level[sl.level-1].forward == nil {
			sl.level--
		}

		// 从映射表中删除
		delete(sl.elementMap, member)
		sl.length--

		return true
	}

	return false
}

// GetRank 获取指定成员的排名，从1开始计算（排名第一的分数最高）
func (sl *SkipList) GetRank(member string, score int64) int64 {
	var rank uint64 = 0
	x := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		// 注意这里的比较逻辑: 分数高的排在前面，分数相同时按照成员ID字典序排序
		for x.level[i].forward != nil &&
			(x.level[i].forward.element.Score > score ||
				(x.level[i].forward.element.Score == score &&
					x.level[i].forward.element.Member < member)) {
			rank += x.level[i].span
			x = x.level[i].forward
		}
	}

	x = x.level[0].forward
	if x != nil && x.element.Member == member {
		return int64(rank + 1)
	}

	return 0
}

// GetByRank 根据排名获取元素，排名从1开始
func (sl *SkipList) GetByRank(rank int64) *Element {
	if rank <= 0 || rank > int64(sl.length) {
		return nil
	}

	var traversed uint64 = 0
	x := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span < uint64(rank) {
			traversed += x.level[i].span
			x = x.level[i].forward
		}

		if traversed+1 == uint64(rank) {
			if x.level[i].forward != nil {
				return &x.level[i].forward.element
			}
		}
	}

	return nil
}

// GetElementByMember 通过成员名获取元素
func (sl *SkipList) GetElementByMember(member string) *Element {
	if node, ok := sl.elementMap[member]; ok {
		return &node.element
	}
	return nil
}

// UpdateScore 更新成员的分数
func (sl *SkipList) UpdateScore(member string, newScore int64) bool {
	if node, ok := sl.elementMap[member]; ok {
		oldScore := node.element.Score
		data := node.element.Data

		// 删除旧节点并添加新节点
		sl.Delete(member, oldScore)
		sl.Insert(member, newScore, data)
		return true
	}
	return false
}

// GetRankRange 获取指定排名范围的元素
func (sl *SkipList) GetRankRange(start, end int64) []*Element {
	var elements []*Element

	// 边界检查
	if start <= 0 {
		start = 1
	}

	if end > int64(sl.length) {
		end = int64(sl.length)
	}

	if start > end {
		return elements
	}

	// 获取指定范围的元素
	for i := start; i <= end; i++ {
		element := sl.GetByRank(i)
		if element != nil {
			elements = append(elements, element)
		}
	}

	return elements
}

// GetScoreRange 获取指定分数范围的元素
func (sl *SkipList) GetScoreRange(minScore, maxScore int64) []*Element {
	var elements []*Element

	// 从头节点开始遍历
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.element.Score > maxScore {
			x = x.level[i].forward
		}
	}

	// 找到第一个分数小于等于maxScore的节点
	x = x.level[0].forward

	// 收集所有分数在范围内的元素
	for x != nil && x.element.Score >= minScore {
		elements = append(elements, &x.element)
		x = x.level[0].forward
	}

	return elements
}

// Len 返回跳表中元素的数量
func (sl *SkipList) Len() uint64 {
	return sl.length
}
