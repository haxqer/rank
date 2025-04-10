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

	// 展示低分优先的赛车排行榜示例
	fmt.Println("\n===============================")
	fmt.Println("低分优先排行榜示例 - 赛车比赛用时")
	fmt.Println("===============================")

	// 创建低分优先排行榜
	raceConfig := rank.LeaderboardConfig{
		ID:           "race_time",
		Name:         "Racing Time Leaderboard",
		ScoreOrder:   false,              // 低分优先
		UpdatePolicy: rank.UpdateIfLower, // 只更新更低的分数（在低分优先的情况下，实际是只接受更高的分数）
	}

	raceLeaderboard := rank.NewLeaderboard(raceConfig)

	// 添加初始比赛成绩
	fmt.Println("\n添加初始比赛成绩:")

	r1, _ := raceLeaderboard.Add("racer1", 120, "选手1完成时间: 120秒")
	fmt.Printf("选手1添加成功, 排名: %d, 用时: %d秒\n", r1.Rank, r1.Score)

	r2, _ := raceLeaderboard.Add("racer2", 105, "选手2完成时间: 105秒")
	fmt.Printf("选手2添加成功, 排名: %d, 用时: %d秒\n", r2.Rank, r2.Score)

	r3, _ := raceLeaderboard.Add("racer3", 130, "选手3完成时间: 130秒")
	fmt.Printf("选手3添加成功, 排名: %d, 用时: %d秒\n", r3.Rank, r3.Score)

	// 显示初始排行榜
	fmt.Println("\n初始赛车排行榜:")
	initialRanks, _ := raceLeaderboard.GetRankList(1, 10)
	for _, item := range initialRanks {
		fmt.Printf("排名: %d, 选手: %s, 用时: %d秒, 备注: %s\n",
			item.Rank, item.Member, item.Score, item.Data)
	}

	// 尝试添加更低的分数，应该失败
	fmt.Println("\n尝试添加更低的用时 (应该失败):")
	_, err := raceLeaderboard.Add("racer2", 100, "选手2新记录: 100秒")
	if err != nil {
		fmt.Printf("预期错误: %v\n", err)
	} else {
		fmt.Println("选手2更新成功，但这不符合预期!")
	}

	// 尝试添加更高的分数，应该成功
	fmt.Println("\n尝试添加更高的用时 (应该成功):")
	updated, err := raceLeaderboard.Add("racer2", 110, "选手2新记录: 110秒")
	if err != nil {
		fmt.Printf("更新选手2的成绩出错: %v\n", err)
	} else {
		fmt.Printf("选手2更新成功, 新排名: %d, 新用时: %d秒\n", updated.Rank, updated.Score)
	}

	// 显示最终排行榜
	fmt.Println("\n最终赛车排行榜:")
	finalRanks, _ := raceLeaderboard.GetRankList(1, 10)
	for _, item := range finalRanks {
		fmt.Printf("排名: %d, 选手: %s, 用时: %d秒, 备注: %s\n",
			item.Rank, item.Member, item.Score, item.Data)
	}
}
