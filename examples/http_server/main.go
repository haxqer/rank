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

// Global leaderboard instance
var gameLeaderboard *rank.Leaderboard

// PlayerScore represents player score request
type PlayerScore struct {
	Member string      `json:"member"`
	Score  int64       `json:"score"`
	Data   interface{} `json:"data,omitempty"`
}

// RankResponse represents ranking response
type RankResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func init() {
	// Initialize leaderboard
	config := rank.LeaderboardConfig{
		ID:           "game_scores",
		Name:         "Game Leaderboard",
		ScoreOrder:   true, // High score first
		UpdatePolicy: rank.UpdateAlways,
	}

	gameLeaderboard = rank.NewLeaderboard(config)

	// Add some initial data
	gameLeaderboard.Add("player1", 1000, map[string]interface{}{
		"nickname": "Player One",
		"level":    10,
	})

	gameLeaderboard.Add("player2", 1500, map[string]interface{}{
		"nickname": "Player Two",
		"level":    15,
	})

	gameLeaderboard.Add("player3", 800, map[string]interface{}{
		"nickname": "Player Three",
		"level":    8,
	})
}

// handleAddScore handles requests to add/update scores
func handleAddScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	var playerScore PlayerScore
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playerScore); err != nil {
		sendResponse(w, false, "Invalid request data", nil)
		return
	}

	// Validate data
	if playerScore.Member == "" {
		sendResponse(w, false, "Member ID is required", nil)
		return
	}

	// Add to leaderboard
	rankData, err := gameLeaderboard.Add(playerScore.Member, playerScore.Score, playerScore.Data)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("Failed to add score: %v", err), nil)
		return
	}

	sendResponse(w, true, "Score added successfully", map[string]interface{}{
		"rank":       rankData.Rank,
		"member":     rankData.Member,
		"score":      rankData.Score,
		"updated_at": rankData.UpdatedAt.Format(time.RFC3339),
	})
}

// handleGetRank handles requests to get a rank
func handleGetRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	member := r.URL.Query().Get("member")
	if member == "" {
		sendResponse(w, false, "Member ID is required", nil)
		return
	}

	// Get rank and data
	rankData, err := gameLeaderboard.GetMemberAndRank(member)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("Failed to get rank: %v", err), nil)
		return
	}

	sendResponse(w, true, "Rank retrieved successfully", map[string]interface{}{
		"rank":       rankData.Rank,
		"member":     rankData.Member,
		"score":      rankData.Score,
		"data":       rankData.Data,
		"updated_at": rankData.UpdatedAt.Format(time.RFC3339),
	})
}

// handleGetTopList handles requests to get the leaderboard
func handleGetTopList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Get parameters
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

	// Get leaderboard
	rankList, err := gameLeaderboard.GetRankList(start, end)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("Failed to get leaderboard: %v", err), nil)
		return
	}

	// Build response
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

	sendResponse(w, true, "Leaderboard retrieved successfully", result)
}

// handleGetAroundMember handles requests to get ranks around a member
func handleGetAroundMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Get parameters
	member := r.URL.Query().Get("member")
	if member == "" {
		sendResponse(w, false, "Member ID is required", nil)
		return
	}

	countStr := r.URL.Query().Get("count")
	count := int64(5)

	if countStr != "" {
		if c, err := strconv.ParseInt(countStr, 10, 64); err == nil && c > 0 {
			count = c
		}
	}

	// Get ranks around member
	rankList, err := gameLeaderboard.GetAroundMember(member, count)
	if err != nil {
		sendResponse(w, false, fmt.Sprintf("Failed to get ranks around member: %v", err), nil)
		return
	}

	// Build response
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

	sendResponse(w, true, "Ranks around member retrieved successfully", result)
}

// handleRemoveMember handles requests to remove a member
func handleRemoveMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method is supported", http.StatusMethodNotAllowed)
		return
	}

	var playerScore PlayerScore
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&playerScore); err != nil {
		sendResponse(w, false, "Invalid request data", nil)
		return
	}

	if playerScore.Member == "" {
		sendResponse(w, false, "Member ID is required", nil)
		return
	}

	// Remove member
	removed := gameLeaderboard.Remove(playerScore.Member)
	if !removed {
		sendResponse(w, false, "Failed to remove member, may not exist", nil)
		return
	}

	sendResponse(w, true, "Member removed successfully", nil)
}

// handleGetTotal gets the total number of members in the leaderboard
func handleGetTotal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	total := gameLeaderboard.GetTotal()
	sendResponse(w, true, "Total retrieved successfully", map[string]interface{}{
		"total": total,
	})
}

// sendResponse sends a JSON response
func sendResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response := RankResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up routes
	http.HandleFunc("/api/score/add", handleAddScore)
	http.HandleFunc("/api/rank/get", handleGetRank)
	http.HandleFunc("/api/rank/top", handleGetTopList)
	http.HandleFunc("/api/rank/around", handleGetAroundMember)
	http.HandleFunc("/api/member/remove", handleRemoveMember)
	http.HandleFunc("/api/total", handleGetTotal)

	// Start server
	fmt.Println("Starting server on :8080")
	fmt.Println("API Endpoints:")
	fmt.Println("- POST /api/score/add - Add or update a score")
	fmt.Println("- GET /api/rank/get?member=xxx - Get a member's rank")
	fmt.Println("- GET /api/rank/top?start=1&end=10 - Get top N from the leaderboard")
	fmt.Println("- GET /api/rank/around?member=xxx&count=5 - Get ranks around a specific member")
	fmt.Println("- DELETE /api/member/remove - Remove a member")
	fmt.Println("- GET /api/total - Get the total number of members in the leaderboard")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
