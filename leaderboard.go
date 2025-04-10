package rank

import (
	"errors"
	"sync"
	"time"
)

// LeaderboardConfig configuration
type LeaderboardConfig struct {
	// ID unique identifier for the leaderboard
	ID string
	// Name display name of the leaderboard
	Name string
	// ScoreOrder score ordering method, true for high scores first, false for low scores first
	ScoreOrder bool
	// UpdatePolicy policy for handling score updates
	UpdatePolicy UpdatePolicy
}

// UpdatePolicy score update policy
type UpdatePolicy int

const (
	// UpdateIfHigher only update when new score is higher than old score
	UpdateIfHigher UpdatePolicy = iota
	// UpdateIfLower only update when new score is lower than old score
	UpdateIfLower
	// UpdateAlways always update the score
	UpdateAlways
)

// MemberData leaderboard member data
type MemberData struct {
	// Member member identifier
	Member string
	// Score member's score
	Score int64
	// Data additional data
	Data interface{}
	// UpdatedAt last update time
	UpdatedAt time.Time
}

// RankData ranking data
type RankData struct {
	// Rank position in the leaderboard
	Rank int64
	// Member member data
	MemberData
}

// Leaderboard implementation
type Leaderboard struct {
	// config configuration information
	config LeaderboardConfig
	// skipList underlying skip list storage
	skipList *SkipList
	// mutex mutex for thread safety
	mutex sync.RWMutex
}

// NewLeaderboard creates a new leaderboard
func NewLeaderboard(config LeaderboardConfig) *Leaderboard {
	return &Leaderboard{
		config:   config,
		skipList: NewSkipList(),
		mutex:    sync.RWMutex{},
	}
}

// Add adds or updates a member's score
func (lb *Leaderboard) Add(member string, score int64, data interface{}) (*RankData, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	// Check if member already exists
	existing := lb.skipList.GetElementByMember(member)

	// Decide whether to update based on update policy
	if existing != nil {
		var existingScore int64
		if md, ok := existing.Data.(MemberData); ok {
			existingScore = md.Score // This is the actual score, not the potentially inverted score stored in the skip list
		} else {
			existingScore = existing.Score
		}

		switch lb.config.UpdatePolicy {
		case UpdateIfHigher:
			// High score priority: new score must be higher
			// Low score priority: new score must be lower (smaller scores are considered "higher")
			if lb.config.ScoreOrder && score <= existingScore {
				return nil, errors.New("new score is not higher than existing score")
			}
			if !lb.config.ScoreOrder && score >= existingScore {
				return nil, errors.New("new score is not lower than existing score")
			}
		case UpdateIfLower:
			// High score priority: new score must be lower (smaller)
			// Low score priority: new score must be higher (higher times are worse)
			if lb.config.ScoreOrder && score >= existingScore {
				return nil, errors.New("new score is not lower than existing score")
			}
			if !lb.config.ScoreOrder && score <= existingScore {
				return nil, errors.New("new score is not higher than existing score")
			}
		}
	}

	// Adapt score ordering: skip list always keeps high scores at the front,
	// so for low-score-first leaderboards, we need to invert the score
	skipListScore := score
	if !lb.config.ScoreOrder {
		skipListScore = -score
	}

	// Update element
	memberData := MemberData{
		Member:    member,
		Score:     score, // Store original score
		Data:      data,
		UpdatedAt: time.Now(),
	}

	lb.skipList.Insert(member, skipListScore, memberData)

	// Get rank
	rank := lb.skipList.GetRank(member, skipListScore)

	return &RankData{
		Rank:       rank,
		MemberData: memberData,
	}, nil
}

// Remove removes a member
func (lb *Leaderboard) Remove(member string) bool {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return false
	}

	return lb.skipList.Delete(member, element.Score)
}

// GetRank gets a member's rank
func (lb *Leaderboard) GetRank(member string) (int64, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return 0, errors.New("member does not exist")
	}

	rank := lb.skipList.GetRank(member, element.Score)
	return rank, nil
}

// GetMember gets a member's data
func (lb *Leaderboard) GetMember(member string) (*MemberData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return nil, errors.New("member does not exist")
	}

	if data, ok := element.Data.(MemberData); ok {
		return &data, nil
	}

	return nil, errors.New("data type error")
}

// GetMemberAndRank gets a member's data and rank
func (lb *Leaderboard) GetMemberAndRank(member string) (*RankData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return nil, errors.New("member does not exist")
	}

	rank := lb.skipList.GetRank(member, element.Score)

	if data, ok := element.Data.(MemberData); ok {
		return &RankData{
			Rank:       rank,
			MemberData: data,
		}, nil
	}

	return nil, errors.New("data type error")
}

// GetRankList gets a list of rankings
func (lb *Leaderboard) GetRankList(start, end int64) ([]*RankData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	elements := lb.skipList.GetRankRange(start, end)
	result := make([]*RankData, 0, len(elements))

	for _, element := range elements {
		// Calculate rank correctly
		member := element.Member
		score := element.Score
		rank := lb.skipList.GetRank(member, score)

		if data, ok := element.Data.(MemberData); ok {
			result = append(result, &RankData{
				Rank:       rank,
				MemberData: data,
			})
		}
	}

	return result, nil
}

// GetAroundMember gets a list of rankings around a specified member
func (lb *Leaderboard) GetAroundMember(member string, count int64) ([]*RankData, error) {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	// Get member's rank
	element := lb.skipList.GetElementByMember(member)
	if element == nil {
		return nil, errors.New("member does not exist")
	}

	rank := lb.skipList.GetRank(member, element.Score)

	// Calculate range
	start := rank - count
	if start < 1 {
		start = 1
	}

	end := rank + count
	if end > int64(lb.skipList.Len()) {
		end = int64(lb.skipList.Len())
	}

	// Get rank list
	return lb.GetRankList(start, end)
}

// GetTotal gets the total number of members in the leaderboard
func (lb *Leaderboard) GetTotal() uint64 {
	lb.mutex.RLock()
	defer lb.mutex.RUnlock()

	return lb.skipList.Len()
}

// Reset resets the leaderboard
func (lb *Leaderboard) Reset() {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	lb.skipList = NewSkipList()
}
