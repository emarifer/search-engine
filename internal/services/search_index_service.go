package services

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type SearchIndex struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Value     string
	Urls      []CrawledUrl   `gorm:"many2many:token_urls;"`
	CreatedAt time.Time      `gorm:"datetime:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"datetime:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SearchIndexServices struct {
	Index      SearchIndex
	IndexStore *gorm.DB
}

func NewSearchIndexServices(
	i SearchIndex, iStore *gorm.DB,
) SearchIndexServices {

	return SearchIndexServices{
		Index:      i,
		IndexStore: iStore,
	}
}

// Gorm override table name (↓ SEE NOTE BELOW ↓)
func (si *SearchIndex) TableName() string {

	return "search_index"
}

func (sis *SearchIndexServices) Save(
	i map[string][]string, crUrls []CrawledUrl,
) error {
	for value, ids := range i {
		newIndex := &SearchIndex{Value: value}
		if err := sis.IndexStore.
			Where(SearchIndex{Value: value}).
			FirstOrCreate(newIndex).Error; err != nil {
			return err
		}

		var urlsToAppend []CrawledUrl
		for _, id := range ids {
			for _, url := range crUrls {
				if url.ID == id {
					urlsToAppend = append(urlsToAppend, url)

					break
				}
			}
		}

		if err := sis.IndexStore.
			Model(&newIndex).
			Association("Urls").
			Append(&urlsToAppend); err != nil {

			return err
		}
	}

	return nil
}

func (sis *SearchIndexServices) SearchFullText(v string) ([]CrawledUrl, error) {
	terms := strings.Fields(v)
	var urls []CrawledUrl

	for _, term := range terms {
		var searchIndexes []SearchIndex
		if err := sis.IndexStore.
			Preload("Urls").
			Where("value LIKE ?", "%"+term+"%").
			Find(&searchIndexes).
			Error; err != nil {
			return nil, err
		}

		for _, searchIndex := range searchIndexes {
			urls = append(urls, searchIndex.Urls...)
		}
	}

	return urls, nil
}

/* OVERRIDE TABLE NAME IN GORM:
https://gorm.io/docs/conventions.html#TableName
https://stackoverflow.com/questions/44589060/how-to-set-singular-name-for-a-table-in-gorm
*/
