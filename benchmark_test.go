package rank

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())
}

// Generate random string ID
func generateID(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Benchmark: skip list insertion
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
			// Pre-generate all IDs and scores to avoid generation during timing
			ids := make([]string, bm.size)
			scores := make([]int64, bm.size)
			for i := 0; i < bm.size; i++ {
				ids[i] = generateID(8)
				scores[i] = rand.Int63n(10000000)
			}

			// Create skip list
			sl := NewSkipList()

			// Insert some elements first to better simulate real-world insertion performance
			// Different pre-insert size for different test sizes
			preInsertSize := bm.size / 2
			for i := 0; i < preInsertSize; i++ {
				sl.Insert(generateID(8), rand.Int63n(10000000), nil)
			}

			// Reset timer
			b.ResetTimer()

			// Execute b.N separate Insert operations and calculate average time
			for i := 0; i < b.N; i++ {
				// Use i as an index to loop through pre-generated IDs and scores
				idx := i % bm.size
				sl.Insert(ids[idx], scores[idx], nil)
			}
		})
	}
}

// Benchmark: skip list lookup
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

		// Randomly select IDs to look up
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

// Benchmark: skip list get rank
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

		// Randomly select IDs and scores to look up
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

// Benchmark: skip list get element by rank
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
				// Randomly get a rank
				rank := int64(rand.Intn(int(sl.Len())) + 1)
				_ = sl.GetByRank(rank)
			}
		})
	}
}

// Benchmark: skip list get rank range
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
				// Randomly get a range
				start := int64(rand.Intn(int(sl.Len()-uint64(bm.rangeSize))) + 1)
				end := start + int64(bm.rangeSize) - 1
				_ = sl.GetRankRange(start, end)
			}
		})
	}
}

// Benchmark: leaderboard add
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
			// Pre-generate all IDs and scores
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

			// Create leaderboard
			lb := NewLeaderboard(config)

			// Pre-insert some elements to simulate real environment
			preInsertSize := bm.size / 2
			for i := 0; i < preInsertSize; i++ {
				lb.Add(generateID(8), rand.Int63n(1000000), nil)
			}

			b.ResetTimer()
			// Test single Add operation performance
			for i := 0; i < b.N; i++ {
				idx := i % bm.size
				_, _ = lb.Add(ids[idx], scores[idx], nil)
			}
		})
	}
}

// Benchmark: leaderboard get rank
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
		// Pre-generate all IDs and scores
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

		// Randomly select IDs to look up
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

// Benchmark: leaderboard get rank list
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
		// Pre-generate all IDs and scores
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
				// Randomly get a range
				start := int64(rand.Intn(int(lb.GetTotal()-uint64(bm.rangeSize))) + 1)
				end := start + int64(bm.rangeSize) - 1
				_, _ = lb.GetRankList(start, end)
			}
		})
	}
}

// Run performance test and generate report
func TestBenchmarkAndReport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip benchmark test")
	}

	// This function is only used for generating reports and will not run in normal tests
	// Use go test -bench=. command to run benchmark test
}
