package model

type DockerHub struct {
	ID     int64  `gorm:"primaryKey;autoIncrement"`
	UserID string `gorm:"type:varchar(64);uniqueIndex:uni_user_hub_name"`

	HubName   string `gorm:"uniqueIndex:uni_user_hub_name"` // 名称，在上传插件时需要指定名称，用户ID和名称唯一确定一个Docker repo
	Address   string // docker repo的地址
	UserName  string // docker repo的用户名
	Password  string // docker repo的登录凭证（password or token）token注意登录expire time
	IsExpired bool   // 标记是否过期
	Note      string
}
