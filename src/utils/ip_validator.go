package utils

import "net"

// IsValidIP 验证IP地址格式是否有效
func IsValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// ValidateIPList 验证IP列表
// maxSize: 最大IP数量限制
// 返回: (是否有效, 错误信息)
func ValidateIPList(ips []string, maxSize int) (bool, string) {
	if len(ips) == 0 {
		return false, "ip list cannot be empty"
	}
	if len(ips) > maxSize {
		return false, "too many ips, maximum is 100"
	}
	for _, ip := range ips {
		if !IsValidIP(ip) {
			return false, "invalid ip format: " + ip
		}
	}
	return true, ""
}
