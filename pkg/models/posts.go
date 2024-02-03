package models

import "time"

type Post struct {
	BaseModel

	Alias            string        `json:"alias" gorm:"uniqueIndex"`
	Title            string        `json:"title"`
	Content          string        `json:"content"`
	Tags             []Tag         `gorm:"many2many:post_tags"`
	Categories       []Category    `gorm:"many2many:post_categories"`
	LikedAccounts    []PostLike    `json:"liked_accounts"`
	DislikedAccounts []PostDislike `json:"disliked_accounts"`
	RepostTo         *Post         `json:"repost_to" gorm:"foreignKey:RepostID"`
	ReplyTo          *Post         `json:"reply_to" gorm:"foreignKey:ReplyID"`
	PublishedAt      time.Time     `json:"published_at"`
	RepostID         *uint         `json:"repost_id"`
	ReplyID          *uint         `json:"reply_id"`
	RealmID          *uint         `json:"realm_id"`
	AuthorID         uint          `json:"author_id"`
	Author           Account       `json:"author"`

	// Dynamic Calculating Values
	LikeCount    int64 `json:"like_count" gorm:"-"`
	DislikeCount int64 `json:"dislike_count" gorm:"-"`
}
