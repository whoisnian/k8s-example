package model

import (
	"database/sql"
	"time"
)

type File struct {
	ID     int64  `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	UserID int64  `gorm:"column:user_id;type:bigint(20);index:files_user_id_idx" json:"user_id"`
	Name   string `gorm:"column:name;type:varchar(255)" json:"name"`
	Size   int64  `gorm:"column:size;type:bigint(20)" json:"size"`
	Digest string `gorm:"column:digest;type:varchar(255);" json:"digest"`

	BucketName string `gorm:"column:bucket_name;type:varchar(255)" json:"bucket_name"`
	ObjectName string `gorm:"column:object_name;type:varchar(255)" json:"object_name"`

	DeletedAt sql.NullTime `gorm:"column:deleted_at;type:datetime;default:null" json:"deleted_at"`
	CreatedAt time.Time    `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
}

func (*File) TableName() string {
	return "files"
}
