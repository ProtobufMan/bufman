package model

import "time"

type Plugin struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	UserID     string `gorm:"type:varchar(64);uniqueIndex:uni_user_id_name"` // 插件名，用户ID、插件名、版本、修订版本组成唯一索引
	UserName   string `gorm:"type:varchar(200);not null"`
	PluginID   string `gorm:"type:varchar(64);unique;not null"`
	PluginName string `gorm:"type:varchar(200);uniqueIndex:uni_plugin"` // 插件名，用户ID、插件名、版本、修订版本组成唯一索引
	Version    string `gorm:"uniqueIndex:uni_user_id_name"`             // 插件版本
	Reversion  uint32 `gorm:"uniqueIndex:uni_user_id_name"`             // 修订版本
	BinaryName string `gorm:"not null"`                                 // 插件二进制执行文件名称

	Description    string    // 插件描述信息
	Visibility     uint8     `gorm:"default:1"` // 可见性，1:public 2:private
	Deprecated     bool      // 是否弃用
	DeprecationMsg string    // 弃用说明
	CreatedTime    time.Time `gorm:"autoCreateTime"`
	UpdateTime     time.Time `gorm:"autoUpdateTime"`
}
