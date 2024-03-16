package services

import (
	"fmt"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func (v *PostTypeContext) ListComment(id uint, take int, offset int, noReact ...bool) ([]*models.Feed, error) {
	if take > 20 {
		take = 20
	}

	var items []*models.Feed
	table := viper.GetString("database.prefix") + "comments"
	userTable := viper.GetString("database.prefix") + "accounts"
	if err := v.Tx.
		Table(table).
		Select("*, ? as model_type", "comment").
		Where(v.ColumnName+"_id = ?", id).
		Joins(fmt.Sprintf("INNER JOIN %s as author ON author_id = author.id", userTable)).
		Limit(take).Offset(offset).Find(&items).Error; err != nil {
		return items, err
	}

	idx := lo.Map(items, func(item *models.Feed, index int) uint {
		return item.ID
	})

	if len(noReact) <= 0 || !noReact[0] {
		var reactions []struct {
			PostID uint
			Symbol string
			Count  int64
		}

		if err := database.C.Model(&models.Reaction{}).
			Select("comment_id as post_id, symbol, COUNT(id) as count").
			Where("comment_id IN (?)", idx).
			Group("post_id, symbol").
			Scan(&reactions).Error; err != nil {
			return items, err
		}

		itemMap := lo.SliceToMap(items, func(item *models.Feed) (uint, *models.Feed) {
			return item.ID, item
		})

		list := map[uint]map[string]int64{}
		for _, info := range reactions {
			if _, ok := list[info.PostID]; !ok {
				list[info.PostID] = make(map[string]int64)
			}
			list[info.PostID][info.Symbol] = info.Count
		}

		for k, v := range list {
			if post, ok := itemMap[k]; ok {
				post.ReactionList = v
			}
		}
	}

	return items, nil
}

func (v *PostTypeContext) CountComment(id uint) (int64, error) {
	var count int64
	if err := database.C.
		Model(&models.Comment{}).
		Where(v.ColumnName+"_id = ?", id).
		Where("published_at <= ?", time.Now()).
		Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}
