package czdb

import (
	"cz-ip-service/src/constant"
	"testing"
)

var Ipv6DbPath = "D:\\360MoveData\\Users\\meet-hive\\Desktop\\czdb\\cz88_public_v6.czdb"
var Ipv4DbPath = "D:\\360MoveData\\Users\\meet-hive\\Desktop\\czdb\\cz88_public_v4.czdb"
var KEY = "xxxx"
var IPV4 = "112.91.94.38"
var IPV6 = "240e:456:230:6dac:f862:6bff:fe95:5108"

func TestMemoryIPV4(t *testing.T) {
	performQuery(IPV4, Ipv4DbPath, constant.MEMORY, t)
}

func TestMemoryIPV6(t *testing.T) {
	performQuery(IPV6, Ipv6DbPath, constant.MEMORY, t)
}

func performQuery(ip, dbPath string, queryType constant.QueryType, t *testing.T) {
	// 执行查询操作
	searcher, err := NewDbSearcher(dbPath, queryType, KEY)
	if err != nil {
		t.Errorf("NewDbSearcher failed: %v", err)
		return
	}

	// 执行查询操作
	result, err := searcher.Search(ip)
	if err != nil {
		t.Errorf("Search failed: %v", err)
		return
	}

	// 打印查询结果
	t.Logf("Query result: %v", result)
}
