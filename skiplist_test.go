package rank

import (
	"testing"
)

func TestSkipListBasic(t *testing.T) {
	sl := NewSkipList()
	
	// 测试插入
	sl.Insert("key1", 100, "value1")
	sl.Insert("key2", 200, "value2")
	sl.Insert("key3", 50, "value3")
	
	// 测试长度
	if sl.Len() != 3 {
		t.Errorf("Expected length 3, got %d", sl.Len())
	}
	
	// 测试获取元素
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
	
	// 测试排名 - 分数从高到低排列，200 > 100 > 50
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
	
	// 测试根据排名获取元素
	element = sl.GetByRank(1)
	if element == nil {
		t.Fatal("Failed to get element for rank 1")
	}
	
	if element.Member != "key2" {
		t.Errorf("Expected member 'key2', got %s", element.Member)
	}
	
	// 测试更新分数
	sl.UpdateScore("key3", 300)
	
	// 检查更新后的排名，300 > 200 > 100
	rank = sl.GetRank("key3", 300)
	if rank != 1 {
		t.Errorf("Expected rank 1 after update, got %d", rank)
	}
	
	// 测试删除
	success := sl.Delete("key1", 100)
	if !success {
		t.Error("Failed to delete key1")
	}
	
	// 检查删除后的长度
	if sl.Len() != 2 {
		t.Errorf("Expected length 2 after deletion, got %d", sl.Len())
	}
	
	// 确认元素已删除
	element = sl.GetElementByMember("key1")
	if element != nil {
		t.Error("Element for key1 should be nil after deletion")
	}
}

func TestSkipListRankRange(t *testing.T) {
	sl := NewSkipList()
	
	// 插入一些数据，注意分数是按降序排列的
	for i := 1; i <= 10; i++ {
		sl.Insert("key"+string(rune('0'+i)), int64((11-i)*100), i)
	}
	
	// 测试获取排名范围
	elements := sl.GetRankRange(2, 5)
	
	if len(elements) != 4 {
		t.Errorf("Expected 4 elements in range, got %d", len(elements))
	}
	
	// 检查第一个元素是否是排名第二的元素 (分数为900)
	if elements[0].Score != 900 {
		t.Errorf("Expected score 900 for rank 2, got %d", elements[0].Score)
	}
	
	// 测试边界条件
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
	
	// 插入一些数据
	for i := 1; i <= 10; i++ {
		sl.Insert("key"+string(rune('0'+i)), int64(i*100), i)
	}
	
	// 测试获取分数范围
	elements := sl.GetScoreRange(300, 700)
	
	if len(elements) != 5 {
		t.Errorf("Expected 5 elements in range, got %d", len(elements))
	}
	
	// 测试边界条件
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
	
	// 大规模插入测试
	for i := 0; i < 1000; i++ {
		key := "key" + string(rune('0'+i%10))
		sl.Insert(key, int64(i), i)
	}
	
	// 检查长度
	if sl.Len() != 10 {
		t.Errorf("Expected length 10 after insertions with duplicate keys, got %d", sl.Len())
	}
	
	// 检查分数是否更新为最后一个插入的值
	for i := 0; i < 10; i++ {
		key := "key" + string(rune('0'+i))
		element := sl.GetElementByMember(key)
		expectedScore := 990 + i
		if element.Score != int64(expectedScore) {
			t.Errorf("Expected score %d for %s, got %d", expectedScore, key, element.Score)
		}
	}
} 