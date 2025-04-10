package main

import (
	"fmt"
	"time"

	"github.com/haxqer/rank"
)

func main() {
	// 创建一个高分优先的排行榜
	config := rank.LeaderboardConfig{
		ID:           "game_score",
		Name:         "游戏分数排行榜",
		ScoreOrder:   true, // 高分在前
		UpdatePolicy: rank.UpdateAlways,
	}
	
	leaderboard := rank.NewLeaderboard(config)
	
	// 添加一些玩家数据
	fmt.Println("添加玩家数据...")
	
	// 可以添加任意类型的额外数据
	type PlayerInfo struct {
		Nickname string
		Level    int
		Avatar   string
	}
	
	// 添加玩家1
	player1 := PlayerInfo{
		Nickname: "超级玩家",
		Level:    10,
		Avatar:   "avatar1.png",
	}
	rankData1, _ := leaderboard.Add("player1", 1000, player1)
	fmt.Printf("玩家1添加成功, 排名: %d, 分数: %d\n", rankData1.Rank, rankData1.Score)
	
	// 添加玩家2
	player2 := PlayerInfo{
		Nickname: "大神玩家",
		Level:    20,
		Avatar:   "avatar2.png",
	}
	rankData2, _ := leaderboard.Add("player2", 1500, player2)
	fmt.Printf("玩家2添加成功, 排名: %d, 分数: %d\n", rankData2.Rank, rankData2.Score)
	
	// 添加玩家3
	player3 := PlayerInfo{
		Nickname: "新手玩家",
		Level:    5,
		Avatar:   "avatar3.png",
	}
	rankData3, _ := leaderboard.Add("player3", 500, player3)
	fmt.Printf("玩家3添加成功, 排名: %d, 分数: %d\n", rankData3.Rank, rankData3.Score)
	
	// 添加玩家4
	player4 := PlayerInfo{
		Nickname: "普通玩家",
		Level:    8,
		Avatar:   "avatar4.png",
	}
	rankData4, _ := leaderboard.Add("player4", 800, player4)
	fmt.Printf("玩家4添加成功, 排名: %d, 分数: %d\n", rankData4.Rank, rankData4.Score)
	
	// 添加玩家5
	player5 := PlayerInfo{
		Nickname: "资深玩家",
		Level:    15,
		Avatar:   "avatar5.png",
	}
	rankData5, _ := leaderboard.Add("player5", 1200, player5)
	fmt.Printf("玩家5添加成功, 排名: %d, 分数: %d\n", rankData5.Rank, rankData5.Score)
	
	fmt.Println("\n排行榜总人数:", leaderboard.GetTotal())
	
	// 获取排行榜
	fmt.Println("\n获取排行榜前3名:")
	topThree, _ := leaderboard.GetRankList(1, 3)
	for _, item := range topThree {
		playerInfo := item.Data.(PlayerInfo)
		fmt.Printf("排名: %d, 成员: %s, 分数: %d, 昵称: %s, 等级: %d\n", 
			item.Rank, item.Member, item.Score, playerInfo.Nickname, playerInfo.Level)
	}
	
	// 更新玩家分数
	fmt.Println("\n更新玩家3的分数...")
	updatedData, _ := leaderboard.Add("player3", 1800, player3)
	fmt.Printf("玩家3更新成功, 新排名: %d, 新分数: %d\n", updatedData.Rank, updatedData.Score)
	
	// 再次获取排行榜
	fmt.Println("\n更新后的排行榜前3名:")
	topThree, _ = leaderboard.GetRankList(1, 3)
	for _, item := range topThree {
		playerInfo := item.Data.(PlayerInfo)
		fmt.Printf("排名: %d, 成员: %s, 分数: %d, 昵称: %s, 等级: %d\n", 
			item.Rank, item.Member, item.Score, playerInfo.Nickname, playerInfo.Level)
	}
	
	// 获取指定玩家的排名
	fmt.Println("\n获取玩家4的排名和数据:")
	player4Data, _ := leaderboard.GetMemberAndRank("player4")
	playerInfo4 := player4Data.Data.(PlayerInfo)
	fmt.Printf("玩家4, 排名: %d, 分数: %d, 昵称: %s, 等级: %d\n", 
		player4Data.Rank, player4Data.Score, playerInfo4.Nickname, playerInfo4.Level)
	
	// 获取玩家周围的排名
	fmt.Println("\n获取玩家4周围的排名(上下各1名):")
	aroundPlayer4, _ := leaderboard.GetAroundMember("player4", 1)
	for _, item := range aroundPlayer4 {
		playerInfo := item.Data.(PlayerInfo)
		fmt.Printf("排名: %d, 成员: %s, 分数: %d, 昵称: %s, 等级: %d\n", 
			item.Rank, item.Member, item.Score, playerInfo.Nickname, playerInfo.Level)
	}
	
	// 删除玩家
	fmt.Println("\n删除玩家5...")
	leaderboard.Remove("player5")
	fmt.Println("删除后的排行榜总人数:", leaderboard.GetTotal())
	
	// 展示时间戳
	fmt.Println("\n获取带有时间戳的玩家数据:")
	player2Data, _ := leaderboard.GetMember("player2")
	fmt.Printf("玩家2, 分数: %d, 更新时间: %s\n", 
		player2Data.Score, player2Data.UpdatedAt.Format(time.RFC3339))
		
	// 创建低分优先的排行榜
	fmt.Println("\n创建一个低分优先的排行榜(如赛车游戏的用时排行)...")
	raceConfig := rank.LeaderboardConfig{
		ID:           "race_time",
		Name:         "赛车用时排行榜",
		ScoreOrder:   false, // 低分在前
		UpdatePolicy: rank.UpdateIfLower, // 只有更低的分数才会更新
	}
	
	raceLeaderboard := rank.NewLeaderboard(raceConfig)
	
	// 添加赛车成绩
	raceLeaderboard.Add("racer1", 120, "玩家1完成用时120秒")
	raceLeaderboard.Add("racer2", 105, "玩家2完成用时105秒")
	raceLeaderboard.Add("racer3", 130, "玩家3完成用时130秒")
	
	// 尝试更新一个更高的分数，应该会失败
	_, err := raceLeaderboard.Add("racer2", 110, "玩家2新的用时110秒")
	if err != nil {
		fmt.Println("预期的错误:", err)
	}
	
	// 尝试更新一个更低的分数，应该会成功
	updated, _ := raceLeaderboard.Add("racer2", 100, "玩家2新的用时100秒")
	fmt.Printf("玩家2更新成功, 新排名: %d, 新用时: %d秒\n", updated.Rank, updated.Score)
	
	// 显示赛车排行榜
	fmt.Println("\n赛车用时排行榜:")
	raceRanks, _ := raceLeaderboard.GetRankList(1, 10)
	for _, item := range raceRanks {
		fmt.Printf("排名: %d, 成员: %s, 用时: %d秒, 信息: %s\n", 
			item.Rank, item.Member, item.Score, item.Data)
	}
} 