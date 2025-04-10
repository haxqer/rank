package rank

import (
	"testing"
)

func TestLeaderboardBasic(t *testing.T) {
	// Create leaderboard
	config := LeaderboardConfig{
		ID:           "test",
		Name:         "Test Leaderboard",
		ScoreOrder:   true,
		UpdatePolicy: UpdateAlways,
	}

	lb := NewLeaderboard(config)

	// Test adding
	rankData, err := lb.Add("member1", 100, "data1")
	if err != nil {
		t.Fatalf("Failed to add member: %v", err)
	}

	if rankData.Rank != 1 {
		t.Errorf("Expected rank 1, got %d", rankData.Rank)
	}

	// Add more members
	lb.Add("member2", 200, "data2")
	lb.Add("member3", 50, "data3")

	// Test total count
	if lb.GetTotal() != 3 {
		t.Errorf("Expected total 3, got %d", lb.GetTotal())
	}

	// Test getting rank
	rank, err := lb.GetRank("member2")
	if err != nil {
		t.Fatalf("Failed to get rank: %v", err)
	}

	if rank != 1 {
		t.Errorf("Expected rank 1, got %d", rank)
	}

	// Test getting member data
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

	// Test getting member and rank
	memberRank, err := lb.GetMemberAndRank("member3")
	if err != nil {
		t.Fatalf("Failed to get member and rank: %v", err)
	}

	if memberRank.Rank != 3 {
		t.Errorf("Expected rank 3, got %d", memberRank.Rank)
	}

	// Test getting rank list
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

	// Test getting members around
	aroundList, err := lb.GetAroundMember("member1", 1)
	if err != nil {
		t.Fatalf("Failed to get around list: %v", err)
	}

	if len(aroundList) != 3 {
		t.Errorf("Expected 3 items in around list, got %d", len(aroundList))
	}

	// Test removing member
	removed := lb.Remove("member3")
	if !removed {
		t.Error("Failed to remove member3")
	}

	if lb.GetTotal() != 2 {
		t.Errorf("Expected total 2 after removal, got %d", lb.GetTotal())
	}

	// Test reset
	lb.Reset()
	if lb.GetTotal() != 0 {
		t.Errorf("Expected total 0 after reset, got %d", lb.GetTotal())
	}
}

func TestLeaderboardUpdatePolicy(t *testing.T) {
	// Test high score priority + UpdateIfHigher policy
	higherConfig := LeaderboardConfig{
		ID:           "higher_test",
		Name:         "High Score Priority Test",
		ScoreOrder:   true, // High score priority
		UpdatePolicy: UpdateIfHigher,
	}

	higherLB := NewLeaderboard(higherConfig)

	// Add initial score
	higherLB.Add("player1", 100, nil)

	// Try to add a lower score, should fail
	_, err := higherLB.Add("player1", 50, nil)
	if err == nil {
		t.Error("Expected error when adding lower score with UpdateIfHigher policy")
	}

	// Try to add a higher score, should succeed
	rankData, err := higherLB.Add("player1", 150, nil)
	if err != nil {
		t.Errorf("Failed to add higher score with UpdateIfHigher policy: %v", err)
	}

	if rankData.Score != 150 {
		t.Errorf("Expected score 150, got %d", rankData.Score)
	}

	// Test low score priority + UpdateIfLower policy
	lowerConfig := LeaderboardConfig{
		ID:           "lower_test",
		Name:         "Low Score Priority Test",
		ScoreOrder:   false, // Low score priority
		UpdatePolicy: UpdateIfLower,
	}

	lowerLB := NewLeaderboard(lowerConfig)

	// Add initial score
	lowerLB.Add("player1", 100, nil)

	// Try to add a lower score, should fail
	_, err = lowerLB.Add("player1", 50, nil)
	if err == nil {
		t.Error("Expected error when adding lower score with UpdateIfLower policy in a low-score-first leaderboard")
	}

	// Try to add a higher score, should succeed
	rankData, err = lowerLB.Add("player1", 150, nil)
	if err != nil {
		t.Errorf("Failed to add higher score with UpdateIfLower policy in a low-score-first leaderboard: %v", err)
	} else if rankData == nil {
		t.Error("Expected non-nil rankData")
	} else if rankData.Score != 150 {
		t.Errorf("Expected score 150, got %d", rankData.Score)
	}

	// Test high score priority + UpdateIfLower policy
	higherLowerConfig := LeaderboardConfig{
		ID:           "higher_lower_test",
		Name:         "High Score Priority Low Score Update Test",
		ScoreOrder:   true,          // High score priority
		UpdatePolicy: UpdateIfLower, // Only accept lower scores
	}

	higherLowerLB := NewLeaderboard(higherLowerConfig)

	// Add initial score
	higherLowerLB.Add("player1", 100, nil)

	// Try to add a higher score, should fail
	_, err = higherLowerLB.Add("player1", 150, nil)
	if err == nil {
		t.Error("Expected error when adding higher score with UpdateIfLower policy in a high-score-first leaderboard")
	}

	// Try to add a lower score, should succeed
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
	// Test high score priority
	highConfig := LeaderboardConfig{
		ID:           "high_first",
		Name:         "High Score First",
		ScoreOrder:   true,
		UpdatePolicy: UpdateAlways,
	}

	highLB := NewLeaderboard(highConfig)

	highLB.Add("player1", 100, nil)
	highLB.Add("player2", 200, nil)
	highLB.Add("player3", 150, nil)

	// Check ranks
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

	// Test low score priority
	lowConfig := LeaderboardConfig{
		ID:           "low_first",
		Name:         "Low Score First",
		ScoreOrder:   false,
		UpdatePolicy: UpdateAlways,
	}

	lowLB := NewLeaderboard(lowConfig)

	lowLB.Add("player1", 100, nil)
	lowLB.Add("player2", 200, nil)
	lowLB.Add("player3", 150, nil)

	// Check ranks
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
