package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	// 请求计数
	totalRequests   uint64
	successRequests uint64
	failedRequests  uint64

	// 响应时间统计
	totalResponseTime uint64
	minResponseTime   uint64 = ^uint64(0)
	maxResponseTime   uint64

	// IP查询统计
	ipQueryCount     = make(map[string]uint64)
	ipQueryCountLock sync.RWMutex
)

// RecordRequest 记录请求
func RecordRequest(success bool) {
	atomic.AddUint64(&totalRequests, 1)
	if success {
		atomic.AddUint64(&successRequests, 1)
	} else {
		atomic.AddUint64(&failedRequests, 1)
	}
}

// RecordResponseTime 记录响应时间
func RecordResponseTime(duration time.Duration) {
	durationMs := uint64(duration.Milliseconds())
	atomic.AddUint64(&totalResponseTime, durationMs)

	// 更新最小响应时间
	for {
		current := atomic.LoadUint64(&minResponseTime)
		if durationMs >= current {
			break
		}
		if atomic.CompareAndSwapUint64(&minResponseTime, current, durationMs) {
			break
		}
	}

	// 更新最大响应时间
	for {
		current := atomic.LoadUint64(&maxResponseTime)
		if durationMs <= current {
			break
		}
		if atomic.CompareAndSwapUint64(&maxResponseTime, current, durationMs) {
			break
		}
	}
}

// RecordIPQuery 记录IP查询
func RecordIPQuery(ip string) {
	ipQueryCountLock.Lock()
	defer ipQueryCountLock.Unlock()
	ipQueryCount[ip]++
}

// GetMetrics 获取所有指标
func GetMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// 请求统计
	metrics["total_requests"] = atomic.LoadUint64(&totalRequests)
	metrics["success_requests"] = atomic.LoadUint64(&successRequests)
	metrics["failed_requests"] = atomic.LoadUint64(&failedRequests)

	// 响应时间统计
	total := atomic.LoadUint64(&totalResponseTime)
	count := atomic.LoadUint64(&successRequests)
	if count > 0 {
		metrics["avg_response_time_ms"] = total / count
	}
	metrics["min_response_time_ms"] = atomic.LoadUint64(&minResponseTime)
	metrics["max_response_time_ms"] = atomic.LoadUint64(&maxResponseTime)

	// IP查询统计
	ipQueryCountLock.RLock()
	topIPs := make(map[string]uint64)
	for ip, count := range ipQueryCount {
		topIPs[ip] = count
	}
	ipQueryCountLock.RUnlock()
	metrics["top_queried_ips"] = topIPs

	return metrics
}
