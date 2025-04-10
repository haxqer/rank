package main

import (
	"fmt"
	"time"

	"github.com/haxqer/rank"
)

func main() {
	// Create a high-score-first leaderboard
	config := rank.LeaderboardConfig{
		ID:           "game_score",
		Name:         "Game Score Leaderboard",
		ScoreOrder:   true, // high scores first
		UpdatePolicy: rank.UpdateAlways,
	}

	leaderboard := rank.NewLeaderboard(config)

	// Add some player data
	fmt.Println("Adding player data...")

	// You can add any type of additional data
	type PlayerInfo struct {
		Nickname string
		Level    int
		Avatar   string
	}

	// Add player 1
	player1 := PlayerInfo{
		Nickname: "Super Player",
		Level:    10,
		Avatar:   "avatar1.png",
	}
	rankData1, _ := leaderboard.Add("player1", 1000, player1)
	fmt.Printf("Player 1 added successfully, Rank: %d, Score: %d\n", rankData1.Rank, rankData1.Score)

	// Add player 2
	player2 := PlayerInfo{
		Nickname: "Pro Player",
		Level:    20,
		Avatar:   "avatar2.png",
	}
	rankData2, _ := leaderboard.Add("player2", 1500, player2)
	fmt.Printf("Player 2 added successfully, Rank: %d, Score: %d\n", rankData2.Rank, rankData2.Score)

	// Add player 3
	player3 := PlayerInfo{
		Nickname: "Beginner",
		Level:    5,
		Avatar:   "avatar3.png",
	}
	rankData3, _ := leaderboard.Add("player3", 500, player3)
	fmt.Printf("Player 3 added successfully, Rank: %d, Score: %d\n", rankData3.Rank, rankData3.Score)

	// Add player 4
	player4 := PlayerInfo{
		Nickname: "Average Player",
		Level:    8,
		Avatar:   "avatar4.png",
	}
	rankData4, _ := leaderboard.Add("player4", 800, player4)
	fmt.Printf("Player 4 added successfully, Rank: %d, Score: %d\n", rankData4.Rank, rankData4.Score)

	// Add player 5
	player5 := PlayerInfo{
		Nickname: "Veteran Player",
		Level:    15,
		Avatar:   "avatar5.png",
	}
	rankData5, _ := leaderboard.Add("player5", 1200, player5)
	fmt.Printf("Player 5 added successfully, Rank: %d, Score: %d\n", rankData5.Rank, rankData5.Score)

	fmt.Println("\nTotal players in leaderboard:", leaderboard.GetTotal())

	// Get leaderboard
	fmt.Println("\nGetting top 3 players:")
	topThree, _ := leaderboard.GetRankList(1, 3)
	for _, item := range topThree {
		playerInfo := item.Data.(PlayerInfo)
		fmt.Printf("Rank: %d, Member: %s, Score: %d, Nickname: %s, Level: %d\n",
			item.Rank, item.Member, item.Score, playerInfo.Nickname, playerInfo.Level)
	}

	// Update player score
	fmt.Println("\nUpdating player 3's score...")
	updatedData, _ := leaderboard.Add("player3", 1800, player3)
	fmt.Printf("Player 3 updated successfully, New Rank: %d, New Score: %d\n", updatedData.Rank, updatedData.Score)

	// Get leaderboard again
	fmt.Println("\nUpdated top 3 players:")
	topThree, _ = leaderboard.GetRankList(1, 3)
	for _, item := range topThree {
		playerInfo := item.Data.(PlayerInfo)
		fmt.Printf("Rank: %d, Member: %s, Score: %d, Nickname: %s, Level: %d\n",
			item.Rank, item.Member, item.Score, playerInfo.Nickname, playerInfo.Level)
	}

	// Get specific player rank
	fmt.Println("\nGetting player 4's rank and data:")
	player4Data, _ := leaderboard.GetMemberAndRank("player4")
	playerInfo4 := player4Data.Data.(PlayerInfo)
	fmt.Printf("Player 4, Rank: %d, Score: %d, Nickname: %s, Level: %d\n",
		player4Data.Rank, player4Data.Score, playerInfo4.Nickname, playerInfo4.Level)

	// Get rankings around player
	fmt.Println("\nGetting rankings around player 4 (1 above and 1 below):")
	aroundPlayer4, _ := leaderboard.GetAroundMember("player4", 1)
	for _, item := range aroundPlayer4 {
		playerInfo := item.Data.(PlayerInfo)
		fmt.Printf("Rank: %d, Member: %s, Score: %d, Nickname: %s, Level: %d\n",
			item.Rank, item.Member, item.Score, playerInfo.Nickname, playerInfo.Level)
	}

	// Delete player
	fmt.Println("\nDeleting player 5...")
	leaderboard.Remove("player5")
	fmt.Println("Total players after deletion:", leaderboard.GetTotal())

	// Show timestamp
	fmt.Println("\nGetting player data with timestamp:")
	player2Data, _ := leaderboard.GetMember("player2")
	fmt.Printf("Player 2, Score: %d, Updated at: %s\n",
		player2Data.Score, player2Data.UpdatedAt.Format(time.RFC3339))

	// Demo of a low-score-first racing leaderboard
	fmt.Println("\n===============================")
	fmt.Println("Low Score First Leaderboard Example - Racing Time")
	fmt.Println("===============================")

	// Create a low-score-first leaderboard
	raceConfig := rank.LeaderboardConfig{
		ID:           "race_time",
		Name:         "Racing Time Leaderboard",
		ScoreOrder:   false,              // Low score first
		UpdatePolicy: rank.UpdateIfLower, // Only update when score is lower (for low-score-first, this actually means accepting only higher scores)
	}

	raceLeaderboard := rank.NewLeaderboard(raceConfig)

	// Add initial race results
	fmt.Println("\nAdding initial race results:")

	r1, _ := raceLeaderboard.Add("racer1", 120, "Racer 1 completion time: 120 seconds")
	fmt.Printf("Racer 1 added successfully, Rank: %d, Time: %d seconds\n", r1.Rank, r1.Score)

	r2, _ := raceLeaderboard.Add("racer2", 105, "Racer 2 completion time: 105 seconds")
	fmt.Printf("Racer 2 added successfully, Rank: %d, Time: %d seconds\n", r2.Rank, r2.Score)

	r3, _ := raceLeaderboard.Add("racer3", 130, "Racer 3 completion time: 130 seconds")
	fmt.Printf("Racer 3 added successfully, Rank: %d, Time: %d seconds\n", r3.Rank, r3.Score)

	// Show initial leaderboard
	fmt.Println("\nInitial Racing Leaderboard:")
	initialRanks, _ := raceLeaderboard.GetRankList(1, 10)
	for _, item := range initialRanks {
		fmt.Printf("Rank: %d, Racer: %s, Time: %d seconds, Note: %s\n",
			item.Rank, item.Member, item.Score, item.Data)
	}

	// Try to add a lower score, should fail
	fmt.Println("\nTrying to add a lower time (should fail):")
	_, err := raceLeaderboard.Add("racer2", 100, "Racer 2 new record: 100 seconds")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	} else {
		fmt.Println("Racer 2 updated successfully, but this is not expected!")
	}

	// Try to add a higher score, should succeed
	fmt.Println("\nTrying to add a higher time (should succeed):")
	updated, err := raceLeaderboard.Add("racer2", 110, "Racer 2 new record: 110 seconds")
	if err != nil {
		fmt.Printf("Error updating racer 2's time: %v\n", err)
	} else {
		fmt.Printf("Racer 2 updated successfully, New Rank: %d, New Time: %d seconds\n", updated.Rank, updated.Score)
	}

	// Show final leaderboard
	fmt.Println("\nFinal Racing Leaderboard:")
	finalRanks, _ := raceLeaderboard.GetRankList(1, 10)
	for _, item := range finalRanks {
		fmt.Printf("Rank: %d, Racer: %s, Time: %d seconds, Note: %s\n",
			item.Rank, item.Member, item.Score, item.Data)
	}
}
