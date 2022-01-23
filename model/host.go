package model

import "time"

type Host struct {
	ID               int       // id
	Hostname         string    // 主机名
	IP               string    // 主机ip
	Port             int       // 端口
	User             string    // 用户
	Password         string    // 密码
	PrivateKey       string    // 私钥
	PrivateKeyPhrase string    // 私钥密码
	Desc             string    // 说明
	CreateTime       time.Time // 创建时间
	UpdateTime       time.Time // 更新时间
}

func (h *Host) TableName() string {
	return "host"
}
