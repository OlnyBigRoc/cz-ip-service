package vo

import (
	"github.com/vmihailenco/msgpack/v5"
	"strings"
)

type Reqs struct {
	IPs []string `json:"ips" form:"ips"`
}

type Result[T any] struct {
	Code    int    `msgpack:"code"`              // 状态码
	Data    T      `msgpack:"data"`              // 数据
	Message string `msgpack:"message,omitempty"` // 消息
}
type IPInfo struct {
	IP       string `json:"ip" msgpack:"ip"`                // IP地址
	Country  string `json:"country" msgpack:"country"`      // 国家
	Province string `json:"province" msgpack:"province"`    // 省份
	City     string `json:"city" msgpack:"city"`            // 市区
	Isp      string `json:"isp" msgpack:"isp"`              // 运营商
	Time     int64  `json:"time,omitempty" msgpack:"time,"` // 耗时
}

func NewInfo(str string, ip string) *IPInfo {
	info := &IPInfo{
		IP: ip,
	}
	// 中国-陕西-西安\t电信
	strArray := strings.Split(str, "\t")
	if len(strArray) == 2 {
		reArray := strings.Split(strArray[0], "–")
		if len(reArray) >= 1 {
			info.Country = strings.ReplaceAll(reArray[0], "0", "")
		}
		if len(reArray) >= 2 {
			info.Province = strings.ReplaceAll(reArray[1], "0", "")
		}
		if len(reArray) >= 3 {
			info.City = strings.ReplaceAll(reArray[2], "0", "")
		}
		info.Isp = strings.ReplaceAll(strArray[1], "0", "")
	}
	return info
}

func (info *IPInfo) ToString(seq ...string) string { // 返回字符串
	if len(seq) == 0 {
		seq = []string{"-"}
	}
	var str []string
	if info.Country != "" {
		str = append(str, info.Country)
	}
	if info.Province != "" {
		str = append(str, info.Province)
	}
	if info.City != "" {
		str = append(str, info.City)
	}
	if info.Isp != "" {
		str = append(str, info.Isp)
	}
	return strings.Join(str, seq[0])
}
func (info *IPInfo) LastValidName() string { // 返回最后一个有效的名字
	if info.City != "" {
		return info.City
	}
	if info.Province != "" {
		return info.Province
	}
	if info.Country != "" {
		return info.Country
	}
	return ""

}

func (r *Result[T]) Error(err error) *Result[T] {
	r.Code = 500
	r.Message = err.Error()
	return r
}
func (r *Result[T]) Success(data T) *Result[T] {
	r.Code = 200
	r.Data = data
	return r
}

func (r *Result[T]) ErrorMsgpack(err error) []byte {
	r.Code = 500
	r.Message = err.Error()
	b, _ := msgpack.Marshal(r)
	return b
}
func (r *Result[T]) SuccessMsgpack(data T) []byte {
	r.Code = 200
	r.Data = data
	b, _ := msgpack.Marshal(r)
	return b
}
