package rank

import (
	"testing"
)

func TestSkipListBasic(t *testing.T) {
	sl := NewSkipList()

	// Test insertion
	sl.Insert("key1", 100, "value1")
	sl.Insert("key2", 200, "value2")
	sl.Insert("key3", 50, "value3")

	// Test length
	if sl.Len() != 3 {
		t.Errorf("Expected length 3, got %d", sl.Len())
	}

	// Test getting element
	element := sl.GetElementByMember("key2")
	if element == nil {
		t.Fatal("Failed to get element for key2")
	}

	if element.Score != 200 {
		t.Errorf("Expected score 200, got %d", element.Score)
	}

	if element.Data != "value2" {
		t.Errorf("Expected data 'value2', got %v", element.Data)
	}

	// Test ranking - scores are sorted in descending order, 200 > 100 > 50
	rank := sl.GetRank("key2", 200)
	if rank != 1 {
		t.Errorf("Expected rank 1, got %d", rank)
	}

	rank = sl.GetRank("key1", 100)
	if rank != 2 {
		t.Errorf("Expected rank 2, got %d", rank)
	}

	rank = sl.GetRank("key3", 50)
	if rank != 3 {
		t.Errorf("Expected rank 3, got %d", rank)
	}

	// Test getting element by rank
	element = sl.GetByRank(1)
	if element == nil {
		t.Fatal("Failed to get element for rank 1")
	}

	if element.Member != "key2" {
		t.Errorf("Expected member 'key2', got %s", element.Member)
	}

	// Test updating score
	sl.UpdateScore("key3", 300)

	// Check ranking after update, 300 > 200 > 100
	rank = sl.GetRank("key3", 300)
	if rank != 1 {
		t.Errorf("Expected rank 1 after update, got %d", rank)
	}

	// Test deletion
	success := sl.Delete("key1", 100)
	if !success {
		t.Error("Failed to delete key1")
	}

	// Check length after deletion
	if sl.Len() != 2 {
		t.Errorf("Expected length 2 after deletion, got %d", sl.Len())
	}

	// Confirm element is deleted
	element = sl.GetElementByMember("key1")
	if element != nil {
		t.Error("Element for key1 should be nil after deletion")
	}
}

func TestSkipListRankRange(t *testing.T) {
	sl := NewSkipList()

	// Insert some data, note that scores are arranged in descending order
	for i := 1; i <= 10; i++ {
		sl.Insert("key"+string(rune('0'+i)), int64((11-i)*100), i)
	}

	// Test getting rank range
	elements := sl.GetRankRange(2, 5)

	if len(elements) != 4 {
		t.Errorf("Expected 4 elements in range, got %d", len(elements))
	}

	// Check if the first element is the second ranked element (score 900)
	if elements[0].Score != 900 {
		t.Errorf("Expected score 900 for rank 2, got %d", elements[0].Score)
	}

	// Test boundary conditions
	elements = sl.GetRankRange(9, 15)
	if len(elements) != 2 {
		t.Errorf("Expected 2 elements in range, got %d", len(elements))
	}

	elements = sl.GetRankRange(0, 3)
	if len(elements) != 3 {
		t.Errorf("Expected 3 elements in range, got %d", len(elements))
	}

	elements = sl.GetRankRange(11, 15)
	if len(elements) != 0 {
		t.Errorf("Expected 0 elements in range, got %d", len(elements))
	}
}

func TestSkipListScoreRange(t *testing.T) {
	sl := NewSkipList()

	// Insert some data
	for i := 1; i <= 10; i++ {
		sl.Insert("key"+string(rune('0'+i)), int64(i*100), i)
	}

	// Test getting score range
	elements := sl.GetScoreRange(300, 700)

	if len(elements) != 5 {
		t.Errorf("Expected 5 elements in range, got %d", len(elements))
	}

	// Test boundary conditions
	elements = sl.GetScoreRange(1100, 2000)
	if len(elements) != 0 {
		t.Errorf("Expected 0 elements in range, got %d", len(elements))
	}

	elements = sl.GetScoreRange(0, 100)
	if len(elements) != 1 {
		t.Errorf("Expected 1 element in range, got %d", len(elements))
	}
}

func TestSkipListLarge(t *testing.T) {
	sl := NewSkipList()

	// Large-scale insertion test
	for i := 0; i < 1000; i++ {
		key := "key" + string(rune('0'+i%10))
		sl.Insert(key, int64(i), i)
	}

	// Check length
	if sl.Len() != 10 {
		t.Errorf("Expected length 10 after insertions with duplicate keys, got %d", sl.Len())
	}

	// Check if scores are updated to the last inserted value
	for i := 0; i < 10; i++ {
		key := "key" + string(rune('0'+i))
		element := sl.GetElementByMember(key)
		expectedScore := 990 + i
		if element.Score != int64(expectedScore) {
			t.Errorf("Expected score %d for %s, got %d", expectedScore, key, element.Score)
		}
	}
}
