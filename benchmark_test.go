package rank

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())
}

// 生成随机字符串ID
func generateID(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 基准测试：跳表插入
func BenchmarkSkipListInsert(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Size_100", 100},
		{"Size_1000", 1000},
		{"Size_10000", 10000},
		{"Size_100000", 100000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// 预生成所有ID和分数，避免在计时部分生成
			ids := make([]string, bm.size)
			scores := make([]int64, bm.size)
			for i := 0; i < bm.size; i++ {
				ids[i] = generateID(8)
				scores[i] = rand.Int63n(10000000)
			}

			// 创建跳表
			sl := NewSkipList()
			
			// 先插入一些元素，以便更准确地模拟真实场景的插入性能
			// 不同大小的测试，预插入不同数量的元素
			preInsertSize := bm.size / 2
			for i := 0; i < preInsertSize; i++ {
				sl.Insert(generateID(8), rand.Int63n(10000000), nil)
			}
			
			// 重置计时器
			b.ResetTimer()
			
			// 执行b.N次单独的Insert操作，计算平均时间
			for i := 0; i < b.N; i++ {
				// 使用循环中的i作为索引，循环读取已经预生成的IDs和分数
				idx := i % bm.size
				sl.Insert(ids[idx], scores[idx], nil)
			}
		})
	}
}

// 基准测试：跳表查找
func BenchmarkSkipListGet(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Size_100", 100},
		{"Size_1000", 1000},
		{"Size_10000", 10000},
		{"Size_100000", 100000},
	}

	for _, bm := range benchmarks {
		ids := make([]string, bm.size)
		scores := make([]int64, bm.size)
		for i := 0; i < bm.size; i++ {
			ids[i] = generateID(8)
			scores[i] = rand.Int63n(1000000)
		}

		sl := NewSkipList()
		for i := 0; i < bm.size; i++ {
			sl.Insert(ids[i], scores[i], nil)
		}

		// 随机选择要查找的ID
		lookupIndices := make([]int, 1000)
		if bm.size < 1000 {
			lookupIndices = make([]int, bm.size)
			for i := 0; i < bm.size; i++ {
				lookupIndices[i] = i
			}
		} else {
			for i := 0; i < 1000; i++ {
				lookupIndices[i] = rand.Intn(bm.size)
			}
		}

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := lookupIndices[i%len(lookupIndices)]
				_ = sl.GetElementByMember(ids[idx])
			}
		})
	}
}

// 基准测试：跳表获取排名
func BenchmarkSkipListGetRank(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Size_100", 100},
		{"Size_1000", 1000},
		{"Size_10000", 10000},
		{"Size_100000", 100000},
	}

	for _, bm := range benchmarks {
		ids := make([]string, bm.size)
		scores := make([]int64, bm.size)
		for i := 0; i < bm.size; i++ {
			ids[i] = generateID(8)
			scores[i] = rand.Int63n(1000000)
		}

		sl := NewSkipList()
		for i := 0; i < bm.size; i++ {
			sl.Insert(ids[i], scores[i], nil)
		}

		// 随机选择要查找的ID和分数
		lookupIndices := make([]int, 1000)
		if bm.size < 1000 {
			lookupIndices = make([]int, bm.size)
			for i := 0; i < bm.size; i++ {
				lookupIndices[i] = i
			}
		} else {
			for i := 0; i < 1000; i++ {
				lookupIndices[i] = rand.Intn(bm.size)
			}
		}

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := lookupIndices[i%len(lookupIndices)]
				_ = sl.GetRank(ids[idx], scores[idx])
			}
		})
	}
}

// 基准测试：跳表按排名获取元素
func BenchmarkSkipListGetByRank(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Size_100", 100},
		{"Size_1000", 1000},
		{"Size_10000", 10000},
		{"Size_100000", 100000},
	}

	for _, bm := range benchmarks {
		ids := make([]string, bm.size)
		scores := make([]int64, bm.size)
		for i := 0; i < bm.size; i++ {
			ids[i] = generateID(8)
			scores[i] = rand.Int63n(1000000)
		}

		sl := NewSkipList()
		for i := 0; i < bm.size; i++ {
			sl.Insert(ids[i], scores[i], nil)
		}

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// 随机获取一个排名
				rank := int64(rand.Intn(int(sl.Len())) + 1)
				_ = sl.GetByRank(rank)
			}
		})
	}
}

// 基准测试：跳表获取排名范围
func BenchmarkSkipListGetRankRange(b *testing.B) {
	benchmarks := []struct {
		name      string
		size      int
		rangeSize int
	}{
		{"Size_100_Range10", 100, 10},
		{"Size_1000_Range10", 1000, 10},
		{"Size_10000_Range10", 10000, 10},
		{"Size_100000_Range10", 100000, 10},
		{"Size_100000_Range100", 100000, 100},
	}

	for _, bm := range benchmarks {
		ids := make([]string, bm.size)
		scores := make([]int64, bm.size)
		for i := 0; i < bm.size; i++ {
			ids[i] = generateID(8)
			scores[i] = rand.Int63n(1000000)
		}

		sl := NewSkipList()
		for i := 0; i < bm.size; i++ {
			sl.Insert(ids[i], scores[i], nil)
		}

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// 随机获取一个范围
				start := int64(rand.Intn(int(sl.Len()-uint64(bm.rangeSize))) + 1)
				end := start + int64(bm.rangeSize) - 1
				_ = sl.GetRankRange(start, end)
			}
		})
	}
}

// 基准测试：排行榜添加元素
func BenchmarkLeaderboardAdd(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Size_100", 100},
		{"Size_1000", 1000},
		{"Size_10000", 10000},
		{"Size_100000", 100000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// 预生成所有ID和分数
			ids := make([]string, bm.size)
			scores := make([]int64, bm.size)
			for i := 0; i < bm.size; i++ {
				ids[i] = generateID(8)
				scores[i] = rand.Int63n(1000000)
			}

			config := LeaderboardConfig{
				ID:           "bench_test",
				Name:         "Benchmark Test",
				ScoreOrder:   true,
				UpdatePolicy: UpdateAlways,
			}

			// 创建排行榜
			lb := NewLeaderboard(config)
			
			// 预先插入一些元素，以模拟真实环境
			preInsertSize := bm.size / 2
			for i := 0; i < preInsertSize; i++ {
				lb.Add(generateID(8), rand.Int63n(1000000), nil)
			}

			b.ResetTimer()
			// 测试单个Add操作的性能
			for i := 0; i < b.N; i++ {
				idx := i % bm.size
				_, _ = lb.Add(ids[idx], scores[idx], nil)
			}
		})
	}
}

// 基准测试：排行榜获取排名
func BenchmarkLeaderboardGetRank(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Size_100", 100},
		{"Size_1000", 1000},
		{"Size_10000", 10000},
		{"Size_100000", 100000},
	}

	for _, bm := range benchmarks {
		// 预生成所有ID和分数
		ids := make([]string, bm.size)
		scores := make([]int64, bm.size)
		for i := 0; i < bm.size; i++ {
			ids[i] = generateID(8)
			scores[i] = rand.Int63n(1000000)
		}

		config := LeaderboardConfig{
			ID:           "bench_test",
			Name:         "Benchmark Test",
			ScoreOrder:   true,
			UpdatePolicy: UpdateAlways,
		}

		lb := NewLeaderboard(config)
		for i := 0; i < bm.size; i++ {
			lb.Add(ids[i], scores[i], nil)
		}

		// 随机选择要查找的ID
		lookupIndices := make([]int, 1000)
		if bm.size < 1000 {
			lookupIndices = make([]int, bm.size)
			for i := 0; i < bm.size; i++ {
				lookupIndices[i] = i
			}
		} else {
			for i := 0; i < 1000; i++ {
				lookupIndices[i] = rand.Intn(bm.size)
			}
		}

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := lookupIndices[i%len(lookupIndices)]
				_, _ = lb.GetRank(ids[idx])
			}
		})
	}
}

// 基准测试：排行榜获取排名列表
func BenchmarkLeaderboardGetRankList(b *testing.B) {
	benchmarks := []struct {
		name      string
		size      int
		rangeSize int
	}{
		{"Size_100_Range10", 100, 10},
		{"Size_1000_Range10", 1000, 10},
		{"Size_10000_Range10", 10000, 10},
		{"Size_100000_Range10", 100000, 10},
		{"Size_100000_Range100", 100000, 100},
	}

	for _, bm := range benchmarks {
		// 预生成所有ID和分数
		ids := make([]string, bm.size)
		scores := make([]int64, bm.size)
		for i := 0; i < bm.size; i++ {
			ids[i] = generateID(8)
			scores[i] = rand.Int63n(1000000)
		}

		config := LeaderboardConfig{
			ID:           "bench_test",
			Name:         "Benchmark Test",
			ScoreOrder:   true,
			UpdatePolicy: UpdateAlways,
		}

		lb := NewLeaderboard(config)
		for i := 0; i < bm.size; i++ {
			lb.Add(ids[i], scores[i], nil)
		}

		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// 随机获取一个范围
				start := int64(rand.Intn(int(lb.GetTotal()-uint64(bm.rangeSize))) + 1)
				end := start + int64(bm.rangeSize) - 1
				_, _ = lb.GetRankList(start, end)
			}
		})
	}
}

// 运行性能测试并生成报告
func TestBenchmarkAndReport(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过基准测试")
	}

	// 此函数仅用于生成报告，实际不会在正常测试中运行
	// 使用 go test -bench=. 命令运行基准测试
}
