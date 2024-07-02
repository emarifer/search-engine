package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Url             string         `gorm:"unique;not null" json:"url"`
	Success         bool           `gorm:"default:null" json:"success"`
	CrawlDuration   time.Duration  `json:"crawlDuration"`
	ResponseCode    int            `gorm:"type:smallint" json:"responseCode"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Headings        string         `json:"headings"`
	LastTested      *time.Time     `json:"lastTested"` // Use pointer so this value can be nil
	Indexed         bool           `json:"indexed" gorm:"default:false"`
	CreatedAt       time.Time      `gorm:"datetime:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"datetime:timestamp" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type UrlServices struct {
	Url      CrawledUrl
	UrlStore *gorm.DB
}

func NewUrlServices(u CrawledUrl, uStore *gorm.DB) UrlServices {

	return UrlServices{
		Url:      u,
		UrlStore: uStore,
	}
}

func (u *UrlServices) UpdateUrl(input CrawledUrl) error {
	tx := u.UrlStore.Select(
		"url",
		"success",
		"crawl_duration",
		"response_code",
		"page_title",
		"page_description",
		"headings",
		"last_tested",
		"updated_at",
	).Omit("created_at").Save(&input)
	if tx.Error != nil {
		return fmt.Errorf("url not updated: %s", tx.Error)
	}

	return nil
}

func (u *UrlServices) GetNextCrawlUrls(limit uint) ([]CrawledUrl, error) {
	var urls []CrawledUrl

	tx := u.UrlStore.Where("last_tested IS NULL").Limit(int(limit)).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, fmt.Errorf("urls not found: %s", tx.Error)
	}

	return urls, nil
}

/* func (u *UrlServices) SaveBatch(input *[]CrawledUrl) error {
	tx := u.UrlStore.Create(input)
	if tx.Error != nil {
		return fmt.Errorf("some urls could not be saved: %s", tx.Error)
	}

	return nil
} */

func (u *UrlServices) Save(input *CrawledUrl) error {
	tx := u.UrlStore.Save(input)
	if tx.Error != nil {
		return fmt.Errorf("the url could not be saved: %s", tx.Error)
	}

	return nil
}

func (u *UrlServices) GetNotIndexed() ([]CrawledUrl, error) {
	var urls []CrawledUrl

	tx := u.UrlStore.
		Where("indexed = ? AND last_tested IS NOT NULL", false).
		Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, fmt.Errorf(
			"something went wrong when getting non-indexed urls: %s", tx.Error,
		)
	}

	return urls, nil
}

func (u *UrlServices) SetIndexedTrue(urls []CrawledUrl) error {
	for _, url := range urls {
		url.Indexed = true
		tx := u.UrlStore.Save(&url)
		if tx.Error != nil {
			return fmt.Errorf(
				"something went wrong when saving indexed urls: %s",
				tx.Error,
			)
		}
	}

	return nil
}
