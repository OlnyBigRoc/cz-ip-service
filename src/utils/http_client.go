package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
}
var HttpClient = &http.Client{
	Timeout:   15 * time.Second,
	Transport: transport,
}

// GetForBody get请求调用
func GetForBody(url string) ([]byte, error) {
	var resp *http.Response
	var err error
	var bodyText []byte
	if resp, err = HttpClient.Get(url); err != nil {
		return bodyText, err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()
	if bodyText, err = io.ReadAll(resp.Body); err != nil {
		return bodyText, err
	}
	return bodyText, nil
}

// PostJsonForBody Post请求调用(会将参数放到请求体里)
func PostJsonForBody(url string, data map[string]interface{}) ([]byte, error) {
	var err error
	var resp *http.Response
	var bodyText []byte
	var queryText []byte
	if queryText, err = json.Marshal(data); err != nil {
		return bodyText, err
	}
	if resp, err = HttpClient.Post(url, "application/json", bytes.NewReader(queryText)); err != nil {
		return bodyText, err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()
	if bodyText, err = io.ReadAll(resp.Body); err != nil {
		return bodyText, err
	}
	return bodyText, err
}

// PostForBody Post请求调用(会将参数拼接在url上)
func PostForBody(url string, data map[string]string) ([]byte, error) {
	var err error
	var resp *http.Response
	if resp, err = HttpClient.PostForm(url, Encode2MapArray(data)); err != nil {
		return nil, err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

// RecursionGetForBody 递归发送get请求(会将参数拼接在url上)
// 递归的作用是有net error等问题可以多尝试几次
func RecursionGetForBody(url string, count uint) ([]byte, error) {
	var resp *http.Response
	var err error
	if resp, err = HttpClient.Get(url); err != nil {
		if count < 3 {
			return RecursionGetForBody(url, count+1)
		}
		return nil, err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

// RecursionGetForCookies 递归发送get请求获取Cookies
// 递归的作用是有net error等问题可以多尝试几次
func RecursionGetForCookies(url string, count uint) (map[string]string, error) {

	var resp *http.Response
	var err error
	if resp, err = HttpClient.Get(url); err != nil {
		if count < 3 {
			return RecursionGetForCookies(url, count+1)
		}
		return nil, err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()
	cookies := resp.Cookies()
	cookieMap := make(map[string]string, len(cookies))
	for i := 0; i < len(cookies); i++ {
		cookieMap[cookies[i].Name] = cookies[i].Value
	}
	return cookieMap, nil
}

func EncodeUrl(address string, args map[string]string) string {
	values := url.Values{}
	for k, v := range args {
		values.Add(k, v)
	}
	return address + values.Encode()
}

// RecursionPostForBody 递归发送post请求(会将参数拼接在url上)
// 递归的作用是有net error等问题可以多尝试几次
func RecursionPostForBody(url string, data map[string]string, count uint) ([]byte, error) {
	var err error
	var resp *http.Response
	var body []byte
	if resp, err = HttpClient.PostForm(url, Encode2MapArray(data)); err != nil {
		if count < 3 {
			return RecursionPostForBody(url, data, count+1)
		}
		return nil, err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

func Encode2MapArray(data map[string]string) map[string][]string {
	res := make(map[string][]string)
	for key := range data {
		res[key] = []string{data[key]}
	}
	return res
}
