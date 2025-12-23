package model

import "time"

type Url struct {
	ID        string    `gorm:"column:id; not null; primary_key;" json:"id"`
	URL       string    `gorm:"column:url; not null;" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at; not null; default:now();" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; not null; default:now();" json:"updated_at"`
}


func (u *Url) GetTableName() string {
	return "short_url"
}

func (u *Url) GetFieldValue(name string) any {
	switch name {
	case "ID":
		return u.ID
	case "URL":
		return u.URL
	case "CreatedAt":
		return u.CreatedAt
	case "UpdatedAt":
		return u.UpdatedAt
	default:
		return nil
	}
}
