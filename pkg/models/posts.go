package models

import (
	"time"
)

type PostReactInfo struct {
	PostID       uint  `json:"post_id"`
	LikeCount    int64 `json:"like_count"`
	DislikeCount int64 `json:"dislike_count"`
	ReplyCount   int64 `json:"reply_count"`
	RepostCount  int64 `json:"repost_count"`
}

type PostBase struct {
	BaseModel

	Alias       string       `json:"alias" gorm:"uniqueIndex"`
	Attachments []Attachment `json:"attachments"`
	PublishedAt *time.Time   `json:"published_at"`

	AuthorID uint    `json:"author_id"`
	Author   Account `json:"author"`

	// TODO Give the reactions & replies & reposts info back
}

func (p PostBase) GetID() uint {
	return p.ID
}

func (p PostBase) GetReplyTo() PostInterface {
	return nil
}

func (p PostBase) GetRepostTo() PostInterface {
	return nil
}

func (p PostBase) GetAuthor() Account {
	return p.Author
}

func (p PostBase) GetRealm() *Realm {
	return nil
}

type PostInterface interface {
	GetID() uint
	GetHashtags() []Tag
	GetCategories() []Category
	GetReplyTo() PostInterface
	GetRepostTo() PostInterface
	GetAuthor() Account
	GetRealm() *Realm

	SetHashtags([]Tag)
	SetCategories([]Category)
}
