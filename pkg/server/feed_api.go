package server

import "C"
import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type FeedItem struct {
	models.BaseModel

	Alias         string `json:"alias"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	ModelType     string `json:"model_type"`
	CommentCount  int64  `json:"comment_count"`
	ReactionCount int64  `json:"reaction_count"`
	AuthorID      uint   `json:"author_id"`
	RealmID       *uint  `json:"realm_id"`

	Author       models.Account   `json:"author" gorm:"embedded"`
	ReactionList map[string]int64 `json:"reaction_list"`
}

const (
	queryArticle = "id, created_at, updated_at, alias, title, NULL as content, description, realm_id, author_id, 'article' as model_type"
	queryMoment  = "id, created_at, updated_at, alias, NULL as title, content, NULL as description, realm_id, author_id, 'moment' as model_type"
)

func listFeed(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	if take > 20 {
		take = 20
	}

	var whereCondition string

	if realmId > 0 {
		whereCondition += fmt.Sprintf("feed.realm_id = %d", realmId)
	} else {
		whereCondition += "feed.realm_id IS NULL"
	}

	var author models.Account
	if len(c.Query("authorId")) > 0 {
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		} else {
			whereCondition += fmt.Sprintf("AND feed.author_id = %d", author.ID)
		}
	}

	var result []*FeedItem

	userTable := viper.GetString("database.prefix") + "accounts"
	commentTable := viper.GetString("database.prefix") + "comments"
	reactionTable := viper.GetString("database.prefix") + "reactions"

	database.C.Raw(fmt.Sprintf(`SELECT feed.*, author.*, 
		COALESCE(comment_count, 0) as comment_count, 
		COALESCE(reaction_count, 0) as reaction_count
		FROM (? UNION ALL ?) as feed
		INNER JOIN %s as author ON author_id = author.id
		LEFT JOIN (SELECT article_id, moment_id, COUNT(*) as comment_count
            FROM %s
            GROUP BY article_id, moment_id) as comments
            ON (feed.model_type = 'article' AND feed.id = comments.article_id) OR 
			   (feed.model_type = 'moment' AND feed.id = comments.moment_id)
        LEFT JOIN (SELECT article_id, moment_id, COUNT(*) as reaction_count
        	FROM %s
            GROUP BY article_id, moment_id) as reactions
            ON (feed.model_type = 'article' AND feed.id = reactions.article_id) OR 
			   (feed.model_type = 'moment' AND feed.id = reactions.moment_id)
		WHERE %s LIMIT ? OFFSET ?`, userTable, commentTable, reactionTable, whereCondition),
		database.C.Select(queryArticle).Model(&models.Article{}),
		database.C.Select(queryMoment).Model(&models.Moment{}),
		take,
		offset,
	).Scan(&result)

	if !c.QueryBool("noReact", false) {
		var reactions []struct {
			PostID uint
			Symbol string
			Count  int64
		}

		revertReaction := func(dataset string) error {
			itemMap := lo.SliceToMap(lo.FilterMap(result, func(item *FeedItem, index int) (*FeedItem, bool) {
				return item, item.ModelType == dataset
			}), func(item *FeedItem) (uint, *FeedItem) {
				return item.ID, item
			})

			idx := lo.Map(lo.Filter(result, func(item *FeedItem, index int) bool {
				return item.ModelType == dataset
			}), func(item *FeedItem, index int) uint {
				return item.ID
			})

			if err := database.C.Model(&models.Reaction{}).
				Select(dataset+"_id as post_id, symbol, COUNT(id) as count").
				Where(dataset+"_id IN (?)", idx).
				Group("post_id, symbol").
				Scan(&reactions).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

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

			return nil
		}

		if err := revertReaction("article"); err != nil {
			return err
		}
		if err := revertReaction("moment"); err != nil {
			return err
		}
	}

	var count int64
	database.C.Raw(`SELECT COUNT(*) FROM (? UNION ALL ?) as feed`,
		database.C.Select(queryArticle).Model(&models.Article{}),
		database.C.Select(queryMoment).Model(&models.Moment{}),
	).Scan(&count)

	return c.JSON(fiber.Map{
		"count": count,
		"data":  result,
	})
}
