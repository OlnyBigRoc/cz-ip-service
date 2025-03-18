package vo

import (
	"errors"
)

type Version struct {
	Code    int    `json:"code"`    // 状态码
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 提示信息
	Data    string `json:"data"`    // 版本号
	Time    string `json:"time"`    // 时间
}

type CheckVersion struct {
	IsUpdate   bool   `json:"isUpdate"`
	Msg        string `json:"msg"`
	OldVersion string `json:"oldVersion"`
	NewVersion string `json:"newVersion"`
}

func NewCheckVersion(isUpdate bool, msg string, oldVersion string, newVersion string) *CheckVersion {
	return &CheckVersion{
		IsUpdate:   isUpdate,
		Msg:        msg,
		OldVersion: oldVersion,
		NewVersion: newVersion,
	}
}
func (receiver *Version) CheckError() error {
	if receiver.Code != 200 {
		return errors.New(receiver.Message)
	}
	return nil
}
