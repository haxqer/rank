# Rank - High Performance Leaderboard Based on Skip List

Rank is a general-purpose leaderboard library implemented in Go, using skip list algorithm for efficient ranking calculation and querying. Skip list is a randomized data structure with expected O(log n) time complexity for insertion, deletion, and search operations, making it particularly suitable for leaderboard functionality.

## Features

- High Performance: Skip list implementation with O(log n) time complexity for all operations
- Versatility: Support for any type of member and additional data
- Ease of Use: Simple and intuitive API design
- Multiple Sorting Options: Support for both high-score-first and low-score-first sorting
- Multiple Update Strategies: Support for various score update strategies (always update/update only if higher/update only if lower)
- Thread Safety: Built-in locking mechanism for concurrent access
- Flexible Querying: Support for various query methods (single rank/rank range/score range/nearby ranks)

## Installation

```shell
go get github.com/haxqer/rank
```

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"github.com/haxqer/rank"
)

func main() {
	// Create a high-score-first leaderboard
	config := rank.LeaderboardConfig{
		ID:           "game_score",
		Name:         "Game Score Leaderboard",
		ScoreOrder:   true, // true means high score first, false means low score first
		UpdatePolicy: rank.UpdateAlways,
	}
	
	leaderboard := rank.NewLeaderboard(config)
	
	// Add some player data
	leaderboard.Add("player1", 1000, "Player 1 Data")
	leaderboard.Add("player2", 1500, "Player 2 Data")
	leaderboard.Add("player3", 500, "Player 3 Data")
	
	// Get top 3 from the leaderboard
	topThree, _ := leaderboard.GetRankList(1, 3)
	for _, item := range topThree {
		fmt.Printf("Rank: %d, Member: %s, Score: %d\n", item.Rank, item.Member, item.Score)
	}
	
	// Get a specific player's rank
	rank, _ := leaderboard.GetRank("player1")
	fmt.Printf("player1's rank: %d\n", rank)
	
	// Get ranks around a player
	around, _ := leaderboard.GetAroundMember("player1", 1) // Get 1 player above and below player1
	for _, item := range around {
		fmt.Printf("Rank: %d, Member: %s, Score: %d\n", item.Rank, item.Member, item.Score)
	}
}
```

### Low-Score-First Leaderboard (e.g., Racing Games)

```go
// Create a low-score-first leaderboard, such as for race completion time
config := rank.LeaderboardConfig{
    ID:           "race_time",
    Name:         "Race Time Leaderboard",
    ScoreOrder:   false, // Low score first
    UpdatePolicy: rank.UpdateIfLower, // Only update when score is lower
}
```

### Using Custom Data Types

```go
// Define a player info type
type PlayerInfo struct {
    Nickname string
    Level    int
    Avatar   string
}

// Add to leaderboard
player := PlayerInfo{
    Nickname: "Super Player",
    Level:    10,
    Avatar:   "avatar.png",
}
leaderboard.Add("player1", 1000, player)

// When getting data, type assertion is needed
rankData, _ := leaderboard.GetMemberAndRank("player1")
playerInfo := rankData.Data.(PlayerInfo)
fmt.Printf("Player Nickname: %s, Level: %d\n", playerInfo.Nickname, playerInfo.Level)
```

## API Reference

### Creating a Leaderboard

```go
// Leaderboard configuration
type LeaderboardConfig struct {
    ID           string      // Unique identifier for the leaderboard
    Name         string      // Display name of the leaderboard
    ScoreOrder   bool        // Score ordering method, true for high scores first, false for low scores first
    UpdatePolicy UpdatePolicy // Update policy
}

// Update policy
const (
    UpdateIfHigher UpdatePolicy = iota // Only update when new score is higher than old score
    UpdateIfLower                      // Only update when new score is lower than old score
    UpdateAlways                       // Always update the score
)

// Create a new leaderboard
func NewLeaderboard(config LeaderboardConfig) *Leaderboard
```

### Adding or Updating a Member

```go
// Add or update a member's score
func (lb *Leaderboard) Add(member string, score int64, data interface{}) (*RankData, error)
```

### Getting Member Rank

```go
// Get a member's rank
func (lb *Leaderboard) GetRank(member string) (int64, error)

// Get a member's data
func (lb *Leaderboard) GetMember(member string) (*MemberData, error)

// Get a member's data and rank
func (lb *Leaderboard) GetMemberAndRank(member string) (*RankData, error)
```

### Getting Leaderboard

```go
// Get a specific rank range from the leaderboard
func (lb *Leaderboard) GetRankList(start, end int64) ([]*RankData, error)

// Get ranks around a specific member
func (lb *Leaderboard) GetAroundMember(member string, count int64) ([]*RankData, error)
```

### Other Operations

```go
// Remove a member
func (lb *Leaderboard) Remove(member string) bool

// Get total number of members in the leaderboard
func (lb *Leaderboard) GetTotal() uint64

// Reset the leaderboard
func (lb *Leaderboard) Reset()
```

## Examples

The project includes multiple examples:

1. [Basic Usage](examples/basic/main.go) - Demonstrates basic leaderboard operations
2. [HTTP Server](examples/http_server/main.go) - Demonstrates how to build an HTTP leaderboard service

### Running the HTTP Server Example

```shell
cd examples/http_server
go run main.go
```

Then you can use the following APIs:

- `POST /api/score/add`: Add or update a score
- `GET /api/rank/get?member=xxx`: Get a specific member's rank
- `GET /api/rank/top?start=1&end=10`: Get top N from the leaderboard
- `GET /api/rank/around?member=xxx&count=5`: Get ranks around a specific member
- `DELETE /api/member/remove`: Remove a member
- `GET /api/total`: Get the total number of members in the leaderboard

## Performance

Based on the skip list implementation, this leaderboard has good performance for all operations:

- Adding/updating members: O(log n)
- Finding a member's rank: O(log n)
- Getting a member at a specific rank: O(log n)
- Getting a rank range: O(log n) + O(m), where m is the range size

### Benchmark Results

The following are benchmark results on an Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz:

#### Skip List Operations Performance

| Operation | Data Size | Execution Time (ns/op) | Memory Allocation (B/op) | Memory Allocation Count (allocs/op) |
|-----------|-----------|------------------------|--------------------------|-------------------------------------|
| Insert    | 100       | 588.4                  | 96                       | 3                                   |
| Insert    | 1,000     | 656.9                  | 96                       | 3                                   |
| Insert    | 10,000    | 1,229                  | 97                       | 3                                   |
| Insert    | 100,000   | 3,082                  | 121                      | 3                                   |
| Lookup    | 100       | 18.43                  | 0                        | 0                                   |
| Lookup    | 1,000     | 33.63                  | 0                        | 0                                   |
| Lookup    | 10,000    | 48.78                  | 0                        | 0                                   |
| Lookup    | 100,000   | 59.12                  | 0                        | 0                                   |
| Get Rank  | 100       | 44.46                  | 0                        | 0                                   |
| Get Rank  | 1,000     | 148.1                  | 0                        | 0                                   |
| Get Rank  | 10,000    | 430.0                  | 0                        | 0                                   |
| Get Rank  | 100,000   | 980.9                  | 0                        | 0                                   |
| Get By Rank | 100     | 61.67                  | 0                        | 0                                   |
| Get By Rank | 1,000   | 125.4                  | 0                        | 0                                   |
| Get By Rank | 10,000  | 372.7                  | 0                        | 0                                   |
| Get By Rank | 100,000 | 1,237                  | 0                        | 0                                   |
| Get Rank Range(10) | 100    | 731.0            | 248                      | 5                                   |
| Get Rank Range(10) | 1,000  | 1,254            | 248                      | 5                                   |
| Get Rank Range(10) | 10,000 | 2,559            | 248                      | 5                                   |
| Get Rank Range(10) | 100,000 | 5,418           | 248                      | 5                                   |
| Get Rank Range(100) | 100,000 | 41,406         | 2,168                    | 8                                   |

#### Leaderboard Operations Performance

| Operation | Data Size | Execution Time (ns/op) | Memory Allocation (B/op) | Memory Allocation Count (allocs/op) |
|-----------|-----------|------------------------|--------------------------|-------------------------------------|
| Add       | 100       | 846.8                  | 240                      | 5                                   |
| Add       | 1,000     | 1,111                  | 240                      | 5                                   |
| Add       | 10,000    | 1,693                  | 241                      | 5                                   |
| Add       | 100,000   | 3,851                  | 261                      | 5                                   |
| Get Rank  | 100       | 64.96                  | 0                        | 0                                   |
| Get Rank  | 1,000     | 207.2                  | 0                        | 0                                   |
| Get Rank  | 10,000    | 512.1                  | 0                        | 0                                   |
| Get Rank  | 100,000   | 1,349                  | 0                        | 0                                   |
| Get Rank List(10) | 100    | 1,352             | 1,128                    | 16                                  |
| Get Rank List(10) | 1,000  | 1,829             | 1,128                    | 16                                  |
| Get Rank List(10) | 10,000 | 3,179             | 1,128                    | 16                                  |
| Get Rank List(10) | 100,000 | 9,527            | 1,128                    | 16                                  |
| Get Rank List(100) | 100,000 | 61,814          | 11,064                   | 109                                 |

Performance Analysis Conclusion: 