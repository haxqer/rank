package rank

import (
	"errors"
	"sync"
	"time"
)

// LeaderboardConfig 配置
type LeaderboardConfig struct {
	// ID 排行榜唯一标识
	ID string
	// Name 排行榜名称
	Name string
	// ScoreOrder 分数排序方式，为true时高分在前，为false时低分在前
	ScoreOrder bool
	// UpdatePolicy 更新策略，决定如何处理分数更新
	UpdatePolicy UpdatePolicy
}

// UpdatePolicy 分数更新策略
type UpdatePolicy int

const (
	// UpdateIfHigher 仅当新分数高于旧分数时更新
	UpdateIfHigher UpdatePolicy = iota
	// UpdateIfLower 仅当新分数低于旧分数时更新
	UpdateIfLower
	// UpdateAlways 始终更新分数
	UpdateAlways
)

// MemberData 排行榜成员数据
type MemberData struct {
	// Member 成员标识
	Member string
	// Score 分数
	Score int64
	// Data 额外数据
	Data interface{}
	// UpdatedAt 更新时间
	UpdatedAt time.Time
}

// RankData 排名数据
type RankData struct {
	// Rank 排名
	Rank int64
	// Member 成员数据
	MemberData
}

// Leaderboard 排行榜实现
type Leaderboard struct {
	// config 配置信息
	config LeaderboardConfig
	// skipList 底层跳表存储
	skipList *SkipList
	// mutex 互斥锁，保证并发安全
	mutex sync.RWMutex
}

// NewLeaderboard 创建新排行榜
func NewLeaderboard(config LeaderboardConfig) *Leaderboard {
	return &Leaderboard{
		config:   config,
		skipList: NewSkipList(),
		mutex:    sync.RWMutex{},
	}
}

// Add 添加或更新成员分数
func (lb *Leaderboard) Add(member string, score int64, data interface{}) (*RankData, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	
	// 检查成员是否已存在
	existing := lb.skipList.GetElementByMember(member)
	
	// 根据更新策略决定是否需要更新
	if existing != nil {
		var existingScore int64
		if md, ok := existing.Data.(MemberData); ok {
			existingScore = md.Score // 这是真实的分数，而非跳表内部存储的可能取反的分数
		} else {
			existingScore = existing.Score
		}
		
		switch lb.config.UpdatePolicy {
		case UpdateIfHigher:
			// 高分优先：新分数必须更高
			// 低分优先：新分数必须更低（更小的分数视为"更高"）
			if lb.config.ScoreOrder && score <= existingScore {
				return nil, errors.New("新分数不高于现有分数")
			}
			if !lb.config.ScoreOrder && score >= existingScore {
				return nil, errors.New("新分数不低于现有分数")
			}
		case UpdateIfLower:
			// 高分优先：新分数必须更低(更小)
			// 低分优先：新分数必须更高(更大)
			if lb.config.ScoreOrder && score >= existingScore {
				return nil, errors.New("新分数不低于现有分数")
			}
			if !lb.config.ScoreOrder && score <= existingScore {
				return nil, errors.New("新分数不高于现有分数")
			}
		}
	}
	
	// 适配分数顺序：跳表内部始终是高分靠前，所以对于低分优先的排行榜，需要将分数取反
	skipListScore := score
	if !lb.config.ScoreOrder {
		skipListScore = -score
	}
	
	// 更新元素
	memberData := MemberData{
		Member:    member,
		Score:     score, // 保存原始分数
		Data:      data,
		UpdatedAt: time.Now(),
	}
	
	lb.skipList.Insert(member, skipListScore, memberData)
	
	// 获取排名
	rank := lb.skipList.GetRank(member, skipListScore)
	
	return &RankData{
		Rank:       rank,
		MemberData: memberData,
	}, nil
}

// Remove 移除成员
func (lb *Leaderboard) Remove(member string) bool {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	
	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return false
	}
	
	return lb.skipList.Delete(member, element.Score)
}

// GetRank 获取成员排名
func (lb *Leaderboard) GetRank(member string) (int64, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	
	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return 0, errors.New("成员不存在")
	}
	
	rank := lb.skipList.GetRank(member, element.Score)
	return rank, nil
}

// GetMember 获取成员数据
func (lb *Leaderboard) GetMember(member string) (*MemberData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	
	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return nil, errors.New("成员不存在")
	}
	
	if data, ok := element.Data.(MemberData); ok {
		return &data, nil
	}
	
	return nil, errors.New("数据类型错误")
}

// GetMemberAndRank 获取成员数据和排名
func (lb *Leaderboard) GetMemberAndRank(member string) (*RankData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	
	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return nil, errors.New("成员不存在")
	}
	
	rank := lb.skipList.GetRank(member, element.Score)
	
	if data, ok := element.Data.(MemberData); ok {
		return &RankData{
			Rank:       rank,
			MemberData: data,
		}, nil
	}
	
	return nil, errors.New("数据类型错误")
}

// GetRankList 获取排行榜列表
func (lb *Leaderboard) GetRankList(start, end int64) ([]*RankData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	
	elements := lb.skipList.GetRankRange(start, end)
	result := make([]*RankData, 0, len(elements))
	
	for i, element := range elements {
		rank := start + int64(i)
		if data, ok := element.Data.(MemberData); ok {
			result = append(result, &RankData{
				Rank:       rank,
				MemberData: data,
			})
		}
	}
	
	return result, nil
}

// GetAroundMember 获取指定成员附近的排名列表
func (lb *Leaderboard) GetAroundMember(member string, count int64) ([]*RankData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	
	// 获取成员排名
	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return nil, errors.New("成员不存在")
	}
	
	rank := lb.skipList.GetRank(member, element.Score)
	
	// 计算范围
	start := rank - count
	if start < 1 {
		start = 1
	}
	
	end := rank + count
	if end > int64(lb.skipList.Len()) {
		end = int64(lb.skipList.Len())
	}
	
	// 获取排名列表
	return lb.GetRankList(start, end)
}

// GetTotal 获取排行榜总成员数
func (lb *Leaderboard) GetTotal() uint64 {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()
	
	return lb.skipList.Len()
}

// Reset 重置排行榜
func (lb *Leaderboard) Reset() {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	
	lb.skipList = NewSkipList()
} 