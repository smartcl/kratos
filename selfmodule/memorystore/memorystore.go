package memorystore

import (
	"sync"
	"time"
)

// Entry 定义值结构体
type Entry struct {
	Value             string    // 字符串类型的值
	ExpiryTime        time.Time // 过期时间
	AllowOverrideTime time.Time // 允许覆盖的时间
}

// SafeMap 线程安全Map结构体
type SafeMap struct {
	mu             sync.RWMutex
	ttl            time.Duration
	overrideWindow time.Duration
	items          map[string]*Entry
}

// NewSafeMap 创建新的SafeMap实例
func NewSafeMap(ttl time.Duration, overrideWindow time.Duration) *SafeMap {
	safeMap := &SafeMap{
		items:          make(map[string]*Entry),
		ttl:            ttl,
		overrideWindow: overrideWindow,
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			safeMap.CleanupExpired()
		}
	}()
	return safeMap
}

// Get 获取值（自动处理过期条目）
func (m *SafeMap) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.items[key]
	if !exists {
		return "", false
	}
	if time.Now().After(entry.ExpiryTime) {
		return "", false
	}
	if time.Now().After(entry.AllowOverrideTime) {
		return entry.Value, false
	}
	return entry.Value, true
}

// Set 设置值（带覆盖时间检查）
func (m *SafeMap) Set(key string, value string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	newEntry := &Entry{
		Value:             value,
		ExpiryTime:        now.Add(m.ttl),
		AllowOverrideTime: now.Add(m.overrideWindow),
	}

	// 检查是否允许覆盖
	if existing, exists := m.items[key]; exists {
		if now.After(existing.AllowOverrideTime) {
			m.items[key] = newEntry
			return true // 超过允许覆盖时间，允许覆盖
		} else {
			return false // 未超过允许覆盖时间，不允许覆盖
		}
	}

	m.items[key] = newEntry
	return true
}

// Delete 删除指定键
func (m *SafeMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

// CleanupExpired 清理过期条目
func (m *SafeMap) CleanupExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, entry := range m.items {
		if now.After(entry.ExpiryTime) {
			delete(m.items, key)
		}
	}
}
