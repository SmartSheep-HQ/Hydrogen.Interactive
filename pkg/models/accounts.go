package models

import "time"

// Account profiles basically fetched from Hydrogen.Passport
// But cache at here for better usage
// At the same time this model can make relations between local models
type Account struct {
	BaseModel

	Name          string        `json:"name"`
	Avatar        string        `json:"avatar"`
	Description   string        `json:"description"`
	EmailAddress  string        `json:"email_address"`
	PowerLevel    int           `json:"power_level"`
	Posts         []Post        `json:"posts" gorm:"foreignKey:AuthorID"`
	LikedPosts    []PostLike    `json:"liked_posts"`
	DislikedPosts []PostDislike `json:"disliked_posts"`
	Realms        []Realm       `json:"realms"`
	ExternalID    uint          `json:"external_id"`
}

type AccountMembership struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FollowerID  uint
	FollowingID uint
}
