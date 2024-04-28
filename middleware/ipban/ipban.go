package ipban

import (
	"sync"
	"time"
)

type IPBanManager struct {
	bans     map[string]time.Time
	mu       sync.RWMutex
	duration time.Duration
}

var (
	manager *IPBanManager
	once    sync.Once
)

// GetManager returns the singleton instance of IPBanManager
func GetManager() *IPBanManager {
	once.Do(func() {
		manager = &IPBanManager{
			bans:     make(map[string]time.Time),
			duration: 24 * time.Hour, // 24小时封禁时间
		}
		// 启动清理过期封禁的goroutine
		go manager.cleanupLoop()
	})
	return manager
}

// IsBanned checks if an IP is currently banned
func (m *IPBanManager) IsBanned(ip string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if banTime, exists := m.bans[ip]; exists {
		if time.Now().Before(banTime) {
			return true
		}
		// 如果封禁时间已过，删除记录
		delete(m.bans, ip)
	}
	return false
}

// BanIP bans an IP for the configured duration
func (m *IPBanManager) BanIP(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.bans[ip] = time.Now().Add(m.duration)
}

// cleanupLoop periodically cleans up expired bans
func (m *IPBanManager) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		m.cleanup()
	}
}

// cleanup removes expired bans
func (m *IPBanManager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for ip, banTime := range m.bans {
		if now.After(banTime) {
			delete(m.bans, ip)
		}
	}
} 