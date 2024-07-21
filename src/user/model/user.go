package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID       int64  `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	Name     string `gorm:"column:name;type:varchar(255)" json:"name"`
	Email    string `gorm:"column:email;type:varchar(255);index:users_email_idx,unique" json:"email"`
	Password string `gorm:"column:password;type:varchar(255)" json:"password"`

	DeletedAt sql.NullTime `gorm:"column:deleted_at;type:datetime;default:null" json:"deleted_at"`
	CreatedAt time.Time    `gorm:"column:created_at;type:datetime;autoCreateTime;not null" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at;type:datetime;autoUpdateTime;not null" json:"updated_at"`
}

type UserJson struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*User) TableName() string {
	return "users"
}

func (raw *User) AsJson() UserJson {
	return UserJson{
		ID:    raw.ID,
		Name:  raw.Name,
		Email: raw.Email,

		CreatedAt: raw.CreatedAt,
		UpdatedAt: raw.UpdatedAt,
	}
}
