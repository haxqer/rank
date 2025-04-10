package rank

import (
	"testing"
)

func TestLeaderboardBasic(t *testing.T) {
	// 创建排行榜
	config := LeaderboardConfig{
		ID:           "test",
		Name:         "测试排行榜",
		ScoreOrder:   true,
		UpdatePolicy: UpdateAlways,
	}
	
	lb := NewLeaderboard(config)
	
	// 测试添加
	rankData, err := lb.Add("member1", 100, "data1")
	if err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}
	
	if rankData.Rank != 1 {
		t.Errorf("Expected rank 1, got %d", rankData.Rank)
	}
	
	// 添加更多成员
	lb.Add("member2", 200, "data2")
	lb.Add("member3", 50, "data3")
	
	// 测试总数
	if lb.GetTotal() != 3 {
		t.Errorf("Expected total 3, got %d", lb.GetTotal())
	}
	
	// 测试获取排名
	rank, err := lb.GetRank("member2")
	if err != nil {
		t.Fatalf("Failed to get rank: %v", err)
	}
	
	if rank != 1 {
		t.Errorf("Expected rank 1, got %d", rank)
	}
	
	// 测试获取成员数据
	memberData, err := lb.GetMember("member1")
	if err != nil {
		t.Fatalf("Failed to get member: %v", err)
	}
	
	if memberData.Score != 100 {
		t.Errorf("Expected score 100, got %d", memberData.Score)
	}
	
	if memberData.Data != "data1" {
		t.Errorf("Expected data 'data1', got %v", memberData.Data)
	}
	
	// 测试获取成员和排名
	memberRank, err := lb.GetMemberAndRank("member3")
	if err != nil {
		t.Fatalf("Failed to get member and rank: %v", err)
	}
	
	if memberRank.Rank != 3 {
		t.Errorf("Expected rank 3, got %d", memberRank.Rank)
	}
	
	// 测试获取排行榜列表
	rankList, err := lb.GetRankList(1, 3)
	if err != nil {
		t.Fatalf("Failed to get rank list: %v", err)
	}
	
	if len(rankList) != 3 {
		t.Errorf("Expected 3 items in rank list, got %d", len(rankList))
	}
	
	if rankList[0].Rank != 1 || rankList[0].Member != "member2" {
		t.Errorf("Expected rank 1 to be member2, got %s", rankList[0].Member)
	}
	
	// 测试获取周围成员
	aroundList, err := lb.GetAroundMember("member1", 1)
	if err != nil {
		t.Fatalf("Failed to get around list: %v", err)
	}
	
	if len(aroundList) != 3 {
		t.Errorf("Expected 3 items in around list, got %d", len(aroundList))
	}
	
	// 测试移除成员
	removed := lb.Remove("member3")
	if !removed {
		t.Error("Failed to remove member3")
	}
	
	if lb.GetTotal() != 2 {
		t.Errorf("Expected total 2 after removal, got %d", lb.GetTotal())
	}
	
	// 测试重置
	lb.Reset()
	if lb.GetTotal() != 0 {
		t.Errorf("Expected total 0 after reset, got %d", lb.GetTotal())
	}
}

func TestLeaderboardUpdatePolicy(t *testing.T) {
	// 测试高分优先 + UpdateIfHigher策略
	higherConfig := LeaderboardConfig{
		ID:           "higher_test",
		Name:         "高分优先测试",
		ScoreOrder:   true, // 高分优先
		UpdatePolicy: UpdateIfHigher,
	}
	
	higherLB := NewLeaderboard(higherConfig)
	
	// 添加初始分数
	higherLB.Add("player1", 100, nil)
	
	// 尝试添加更低的分数，应该失败
	_, err := higherLB.Add("player1", 50, nil)
	if err == nil {
		t.Error("Expected error when adding lower score with UpdateIfHigher policy")
	}
	
	// 尝试添加更高的分数，应该成功
	rankData, err := higherLB.Add("player1", 150, nil)
	if err != nil {
		t.Errorf("Failed to add higher score with UpdateIfHigher policy: %v", err)
	}
	
	if rankData.Score != 150 {
		t.Errorf("Expected score 150, got %d", rankData.Score)
	}
	
	// 测试低分优先 + UpdateIfLower策略
	lowerConfig := LeaderboardConfig{
		ID:           "lower_test",
		Name:         "低分优先测试",
		ScoreOrder:   false, // 低分优先
		UpdatePolicy: UpdateIfLower,
	}
	
	lowerLB := NewLeaderboard(lowerConfig)
	
	// 添加初始分数
	lowerLB.Add("player1", 100, nil)
	
	// 尝试添加更高的分数（对于低分优先，这相当于更低的排名），应该失败
	_, err = lowerLB.Add("player1", 150, nil)
	if err == nil {
		t.Error("Expected error when adding higher score with UpdateIfLower policy in a low-score-first leaderboard")
	}
	
	// 尝试添加更低的分数（对于低分优先，这相当于更高的排名），应该成功
	rankData, err = lowerLB.Add("player1", 50, nil)
	if err != nil {
		t.Errorf("Failed to add lower score with UpdateIfLower policy in a low-score-first leaderboard: %v", err)
	} else if rankData == nil {
		t.Error("Expected non-nil rankData")
	} else if rankData.Score != 50 {
		t.Errorf("Expected score 50, got %d", rankData.Score)
	}
	
	// 测试高分优先 + UpdateIfLower策略
	higherLowerConfig := LeaderboardConfig{
		ID:           "higher_lower_test",
		Name:         "高分优先低分更新测试",
		ScoreOrder:   true, // 高分优先
		UpdatePolicy: UpdateIfLower, // 只接受更低的分数
	}
	
	higherLowerLB := NewLeaderboard(higherLowerConfig)
	
	// 添加初始分数
	higherLowerLB.Add("player1", 100, nil)
	
	// 尝试添加更高的分数，应该失败
	_, err = higherLowerLB.Add("player1", 150, nil)
	if err == nil {
		t.Error("Expected error when adding higher score with UpdateIfLower policy in a high-score-first leaderboard")
	}
	
	// 尝试添加更低的分数，应该成功
	rankData, err = higherLowerLB.Add("player1", 50, nil)
	if err != nil {
		t.Errorf("Failed to add lower score with UpdateIfLower policy in a high-score-first leaderboard: %v", err)
	} else if rankData == nil {
		t.Error("Expected non-nil rankData")
	} else if rankData.Score != 50 {
		t.Errorf("Expected score 50, got %d", rankData.Score)
	}
}

func TestLeaderboardScoreOrder(t *testing.T) {
	// 测试高分优先
	highConfig := LeaderboardConfig{
		ID:           "high_first",
		Name:         "高分优先",
		ScoreOrder:   true,
		UpdatePolicy: UpdateAlways,
	}
	
	highLB := NewLeaderboard(highConfig)
	
	highLB.Add("player1", 100, nil)
	highLB.Add("player2", 200, nil)
	highLB.Add("player3", 150, nil)
	
	// 检查排名
	rank, _ := highLB.GetRank("player2")
	if rank != 1 {
		t.Errorf("Expected player2 to be rank 1, got %d", rank)
	}
	
	rank, _ = highLB.GetRank("player3")
	if rank != 2 {
		t.Errorf("Expected player3 to be rank 2, got %d", rank)
	}
	
	rank, _ = highLB.GetRank("player1")
	if rank != 3 {
		t.Errorf("Expected player1 to be rank 3, got %d", rank)
	}
	
	// 测试低分优先
	lowConfig := LeaderboardConfig{
		ID:           "low_first",
		Name:         "低分优先",
		ScoreOrder:   false,
		UpdatePolicy: UpdateAlways,
	}
	
	lowLB := NewLeaderboard(lowConfig)
	
	lowLB.Add("player1", 100, nil)
	lowLB.Add("player2", 200, nil)
	lowLB.Add("player3", 150, nil)
	
	// 检查排名
	rank, _ = lowLB.GetRank("player1")
	if rank != 1 {
		t.Errorf("Expected player1 to be rank 1, got %d", rank)
	}
	
	rank, _ = lowLB.GetRank("player3")
	if rank != 2 {
		t.Errorf("Expected player3 to be rank 2, got %d", rank)
	}
	
	rank, _ = lowLB.GetRank("player2")
	if rank != 3 {
		t.Errorf("Expected player2 to be rank 3, got %d", rank)
	}
} 