# Rank - 基于跳表的高性能排行榜

Rank是一个Go语言实现的通用排行榜库，使用跳表算法实现高效的排名计算和查询。跳表是一种随机化数据结构，在插入、删除和查找操作上都有O(log n)的期望时间复杂度，特别适合实现排行榜功能。

## 特性

- 高性能：跳表实现，各项操作时间复杂度均为O(log n)
- 通用性：支持任意类型的成员和额外数据
- 易用性：简单直观的API设计
- 多排序支持：支持高分优先和低分优先两种排序方式
- 多更新策略：支持多种分数更新策略（总是更新/仅更高分更新/仅更低分更新）
- 线程安全：内置锁机制，支持并发访问
- 灵活查询：支持多种查询方式（单个排名/排名范围/分数范围/周围排名）

## 安装

```shell
go get github.com/haxqer/rank
```

## 使用方法

### 基本用法

```go
package main

import (
	"fmt"
	"github.com/haxqer/rank"
)

func main() {
	// 创建一个高分优先的排行榜
	config := rank.LeaderboardConfig{
		ID:           "game_score",
		Name:         "游戏分数排行榜",
		ScoreOrder:   true, // true表示高分在前，false表示低分在前
		UpdatePolicy: rank.UpdateAlways,
	}
	
	leaderboard := rank.NewLeaderboard(config)
	
	// 添加一些玩家数据
	leaderboard.Add("player1", 1000, "玩家1数据")
	leaderboard.Add("player2", 1500, "玩家2数据")
	leaderboard.Add("player3", 500, "玩家3数据")
	
	// 获取排行榜前3名
	topThree, _ := leaderboard.GetRankList(1, 3)
	for _, item := range topThree {
		fmt.Printf("排名: %d, 成员: %s, 分数: %d\n", item.Rank, item.Member, item.Score)
	}
	
	// 获取指定玩家的排名
	rank, _ := leaderboard.GetRank("player1")
	fmt.Printf("player1的排名: %d\n", rank)
	
	// 获取玩家周围的排名
	around, _ := leaderboard.GetAroundMember("player1", 1) // 获取player1上下各1名的玩家
	for _, item := range around {
		fmt.Printf("排名: %d, 成员: %s, 分数: %d\n", item.Rank, item.Member, item.Score)
	}
}
```

### 低分优先的排行榜（如竞速游戏）

```go
// 创建一个低分优先的排行榜，如竞速游戏的完成时间
config := rank.LeaderboardConfig{
    ID:           "race_time",
    Name:         "竞速时间排行榜",
    ScoreOrder:   false, // 低分在前
    UpdatePolicy: rank.UpdateIfLower, // 只在分数更低时更新
}
```

### 使用自定义数据类型

```go
// 定义玩家信息类型
type PlayerInfo struct {
    Nickname string
    Level    int
    Avatar   string
}

// 添加到排行榜
player := PlayerInfo{
    Nickname: "超级玩家",
    Level:    10,
    Avatar:   "avatar.png",
}
leaderboard.Add("player1", 1000, player)

// 获取数据时需要类型断言
rankData, _ := leaderboard.GetMemberAndRank("player1")
playerInfo := rankData.Data.(PlayerInfo)
fmt.Printf("玩家昵称: %s, 等级: %d\n", playerInfo.Nickname, playerInfo.Level)
```

## API参考

### 创建排行榜

```go
// 排行榜配置
type LeaderboardConfig struct {
    ID           string      // 排行榜唯一标识
    Name         string      // 排行榜名称
    ScoreOrder   bool        // 分数排序方式，为true时高分在前，为false时低分在前
    UpdatePolicy UpdatePolicy // 更新策略
}

// 更新策略
const (
    UpdateIfHigher UpdatePolicy = iota // 仅当新分数高于旧分数时更新
    UpdateIfLower                      // 仅当新分数低于旧分数时更新
    UpdateAlways                       // 始终更新分数
)

// 创建新排行榜
func NewLeaderboard(config LeaderboardConfig) *Leaderboard
```

### 添加或更新成员

```go
// 添加或更新成员分数
func (lb *Leaderboard) Add(member string, score int64, data interface{}) (*RankData, error)
```

### 获取成员排名

```go
// 获取成员排名
func (lb *Leaderboard) GetRank(member string) (int64, error)

// 获取成员数据
func (lb *Leaderboard) GetMember(member string) (*MemberData, error)

// 获取成员数据和排名
func (lb *Leaderboard) GetMemberAndRank(member string) (*RankData, error)
```

### 获取排行榜

```go
// 获取指定排名范围的排行榜
func (lb *Leaderboard) GetRankList(start, end int64) ([]*RankData, error)

// 获取指定成员周围的排名列表
func (lb *Leaderboard) GetAroundMember(member string, count int64) ([]*RankData, error)
```

### 其他操作

```go
// 移除成员
func (lb *Leaderboard) Remove(member string) bool

// 获取排行榜总成员数
func (lb *Leaderboard) GetTotal() uint64

// 重置排行榜
func (lb *Leaderboard) Reset()
```

## 示例

项目包含多个示例：

1. [基本用法](examples/basic/main.go) - 展示排行榜的基本操作
2. [HTTP服务](examples/http_server/main.go) - 演示如何构建一个HTTP排行榜服务

### 运行HTTP服务示例

```shell
cd examples/http_server
go run main.go
```

然后可以使用以下API：

- `POST /api/score/add`: 添加或更新分数
- `GET /api/rank/get?member=xxx`: 获取指定成员的排名
- `GET /api/rank/top?start=1&end=10`: 获取排行榜前N名
- `GET /api/rank/around?member=xxx&count=5`: 获取指定成员周围的排名
- `DELETE /api/member/remove`: 删除成员
- `GET /api/total`: 获取排行榜总人数

## 性能

基于跳表的实现，该排行榜在各项操作上都有很好的性能表现：

- 添加/更新成员: O(log n)
- 查找成员排名: O(log n)
- 获取指定排名的成员: O(log n)
- 获取排名范围: O(log n) + O(m)，其中m是范围大小

### 性能测试结果

以下是在 Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz 环境下的性能测试结果：

#### 跳表操作性能

| 操作 | 数据规模 | 执行时间 (ns/op) | 内存分配 (B/op) | 内存分配次数 (allocs/op) |
|------|----------|------------------|----------------|-------------------------|
| 插入 | 100      | 588.4            | 96             | 3                       |
| 插入 | 1,000    | 656.9            | 96             | 3                       |
| 插入 | 10,000   | 1,229            | 97             | 3                       |
| 插入 | 100,000  | 3,082            | 121            | 3                       |
| 查找 | 100      | 18.43            | 0              | 0                       |
| 查找 | 1,000    | 33.63            | 0              | 0                       |
| 查找 | 10,000   | 48.78            | 0              | 0                       |
| 查找 | 100,000  | 59.12            | 0              | 0                       |
| 获取排名 | 100      | 44.46        | 0              | 0                       |
| 获取排名 | 1,000    | 148.1        | 0              | 0                       |
| 获取排名 | 10,000   | 430.0        | 0              | 0                       |
| 获取排名 | 100,000  | 980.9        | 0              | 0                       |
| 按排名获取 | 100      | 61.67      | 0              | 0                       |
| 按排名获取 | 1,000    | 125.4      | 0              | 0                       |
| 按排名获取 | 10,000   | 372.7      | 0              | 0                       |
| 按排名获取 | 100,000  | 1,237      | 0              | 0                       |
| 获取排名范围(10) | 100      | 731.0 | 248            | 5                       |
| 获取排名范围(10) | 1,000    | 1,254 | 248            | 5                       |
| 获取排名范围(10) | 10,000   | 2,559 | 248            | 5                       |
| 获取排名范围(10) | 100,000  | 5,418 | 248            | 5                       |
| 获取排名范围(100) | 100,000  | 41,406 | 2,168        | 8                       |

#### 排行榜操作性能

| 操作 | 数据规模 | 执行时间 (ns/op) | 内存分配 (B/op) | 内存分配次数 (allocs/op) |
|------|----------|------------------|----------------|-------------------------|
| 添加 | 100      | 846.8            | 240            | 5                       |
| 添加 | 1,000    | 1,111            | 240            | 5                       |
| 添加 | 10,000   | 1,693            | 241            | 5                       |
| 添加 | 100,000  | 3,851            | 261            | 5                       |
| 获取排名 | 100      | 64.96        | 0              | 0                       |
| 获取排名 | 1,000    | 207.2        | 0              | 0                       |
| 获取排名 | 10,000   | 512.1        | 0              | 0                       |
| 获取排名 | 100,000  | 1,349        | 0              | 0                       |
| 获取排名列表(10) | 100     | 1,352  | 1,128          | 16                      |
| 获取排名列表(10) | 1,000   | 1,829  | 1,128          | 16                      |
| 获取排名列表(10) | 10,000  | 3,179  | 1,128          | 16                      |
| 获取排名列表(10) | 100,000 | 9,527  | 1,128          | 16                      |
| 获取排名列表(100) | 100,000 | 61,814 | 11,064        | 109                     |

性能分析结论：

1. **插入操作极快**：单次插入操作在数据量达到10万级别时也只需要约3微秒，内存分配非常少
2. **排行榜添加高效**：即使在大数据量(10万)的情况下，单次添加操作也只需约3.8微秒
3. **查找操作极为高效**：即使在10万数据规模下，成员查找操作也只需要约59ns，基本无内存分配
4. **排名查询出色**：10万数据规模下的排名查询为约980ns至1.3μs，符合O(log n)的复杂度预期
5. **范围查询优异**：获取10个元素的排名范围在10万数据规模下仅需约5.4μs至9.5μs
6. **内存使用合理**：对于查询操作几乎无额外内存分配，插入操作的内存分配很少且不随数据规模增长

## 线程安全

排行榜实现使用了读写锁机制，可以安全地在多个goroutine中并发使用。对于读多写少的场景（典型的排行榜使用模式），性能表现尤为出色。

## 许可证

MIT License

## 贡献

欢迎提交issues和PR！ 