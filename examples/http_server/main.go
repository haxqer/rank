package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/haxqer/rank"
)

// 全局排行榜实例
var gameLeaderboard *rank.Leaderboard

// PlayerScore 玩家分数请求
type PlayerScore struct {
	Member string      `json:"member"`
	Score  int64       `json:"score"`
	Data   interface{} `json:"data,omitempty"`
}

// RankResponse 排名响应
type RankResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func init() {
	// 初始化排行榜
	config := rank.LeaderboardConfig{
		ID:           "game_scores",
		Name:         "游戏排行榜",
		ScoreOrder:   true, // 高分在前
		UpdatePolicy: rank.UpdateAlways,
	}
	
	gameLeaderboard = rank.NewLeaderboard(config)
	
	// 添加一些初始数据
	gameLeaderboard.Add("player1", 1000, map[string]interface{}{
		"nickname": "玩家一号",
		"level":    10,
	})
	
	gameLeaderboard.Add("player2", 1500, map[string]interface{}{
		"nickname": "玩家二号",
		"level":    15,
	})
	
	gameLeaderboard.Add("player3", 800, map[string]interface{}{
		"nickname": "玩家三号",
		"level":    8,
	})
}

// handleAddScore 处理添加/更新分数的请求
func handleAddScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}
	
	var playerScore PlayerScore
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playerScore); err != nil {
		sendResponse(w, false, "无效的请求数据", nil)
		return
	}
	
	// 验证数据
	if playerScore.Member == "" {
		sendResponse(w, false, "必须提供成员ID", nil)
		return
	}
	
	// 添加到排行榜
	rankData, err := gameLeaderboard.Add(playerScore.Member, playerScore.Score, playerScore.Data)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("添加分数失败: %v", err), nil)
		return
	}
	
	sendResponse(w, true, "分数添加成功", map[string]interface{}{
		"rank":       rankData.Rank,
		"member":     rankData.Member,
		"score":      rankData.Score,
		"updated_at": rankData.UpdatedAt.Format(time.RFC3339),
	})
}

// handleGetRank 处理获取排名的请求
func handleGetRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}
	
	member := r.URL.Query().Get("member")
	if member == "" {
		sendResponse(w, false, "必须提供成员ID", nil)
		return
	}
	
	// 获取排名和数据
	rankData, err := gameLeaderboard.GetMemberAndRank(member)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("获取排名失败: %v", err), nil)
		return
	}
	
	sendResponse(w, true, "获取排名成功", map[string]interface{}{
		"rank":       rankData.Rank,
		"member":     rankData.Member,
		"score":      rankData.Score,
		"data":       rankData.Data,
		"updated_at": rankData.UpdatedAt.Format(time.RFC3339),
	})
}

// handleGetTopList 处理获取排行榜的请求
func handleGetTopList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}
	
	// 获取参数
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	
	start := int64(1)
	end := int64(10)
	
	if startStr != "" {
		if s, err := strconv.ParseInt(startStr, 10, 64); err == nil && s > 0 {
			start = s
		}
	}
	
	if endStr != "" {
		if e, err := strconv.ParseInt(endStr, 10, 64); err == nil && e >= start {
			end = e
		}
	}
	
	// 获取排行榜
	rankList, err := gameLeaderboard.GetRankList(start, end)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("获取排行榜失败: %v", err), nil)
		return
	}
	
	// 构建响应
	var result []map[string]interface{}
	for _, item := range rankList {
		result = append(result, map[string]interface{}{
			"rank":       item.Rank,
			"member":     item.Member,
			"score":      item.Score,
			"data":       item.Data,
			"updated_at": item.UpdatedAt.Format(time.RFC3339),
		})
	}
	
	sendResponse(w, true, "获取排行榜成功", result)
}

// handleGetAroundMember 处理获取指定成员周围排名的请求
func handleGetAroundMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}
	
	// 获取参数
	member := r.URL.Query().Get("member")
	if member == "" {
		sendResponse(w, false, "必须提供成员ID", nil)
		return
	}
	
	countStr := r.URL.Query().Get("count")
	count := int64(5)
	
	if countStr != "" {
		if c, err := strconv.ParseInt(countStr, 10, 64); err == nil && c > 0 {
			count = c
		}
	}
	
	// 获取周围的排名
	rankList, err := gameLeaderboard.GetAroundMember(member, count)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("获取周围排名失败: %v", err), nil)
		return
	}
	
	// 构建响应
	var result []map[string]interface{}
	for _, item := range rankList {
		result = append(result, map[string]interface{}{
			"rank":       item.Rank,
			"member":     item.Member,
			"score":      item.Score,
			"data":       item.Data,
			"updated_at": item.UpdatedAt.Format(time.RFC3339),
		})
	}
	
	sendResponse(w, true, "获取周围排名成功", result)
}

// handleRemoveMember 处理删除成员的请求
func handleRemoveMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "只支持DELETE方法", http.StatusMethodNotAllowed)
		return
	}
	
	var playerScore PlayerScore
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playerScore); err != nil {
		sendResponse(w, false, "无效的请求数据", nil)
		return
	}
	
	if playerScore.Member == "" {
		sendResponse(w, false, "必须提供成员ID", nil)
		return
	}
	
	// 删除成员
	removed := gameLeaderboard.Remove(playerScore.Member)
	if !removed {
		sendResponse(w, false, "删除成员失败，可能不存在", nil)
		return
	}
	
	sendResponse(w, true, "成员删除成功", nil)
}

// handleGetTotal 获取排行榜总人数
func handleGetTotal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}
	
	total := gameLeaderboard.GetTotal()
	
	sendResponse(w, true, "获取总人数成功", map[string]interface{}{
		"total": total,
	})
}

// sendResponse 发送统一格式的响应
func sendResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
	response := RankResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
	
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		http.Error(w, "响应编码失败", http.StatusInternalServerError)
	}
}

func main() {
	// 设置路由
	http.HandleFunc("/api/score/add", handleAddScore)
	http.HandleFunc("/api/rank/get", handleGetRank)
	http.HandleFunc("/api/rank/top", handleGetTopList)
	http.HandleFunc("/api/rank/around", handleGetAroundMember)
	http.HandleFunc("/api/member/remove", handleRemoveMember)
	http.HandleFunc("/api/total", handleGetTotal)
	
	// 启动服务器
	port := ":8080"
	fmt.Printf("HTTP服务器启动，监听在 %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
} 