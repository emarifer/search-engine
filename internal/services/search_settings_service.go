package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SearchSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SearchOn  bool      `json:"searchOn"`
	AddNew    bool      `json:"addNew"`
	Amount    uint      `json:"amount"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SearchSettingsServices struct {
	SearchSettings      SearchSettings
	SearchSettingsStore *gorm.DB
}

func NewSearchSettingsServices(
	ss SearchSettings, sstore *gorm.DB,
) SearchSettingsServices {

	return SearchSettingsServices{
		SearchSettings:      ss,
		SearchSettingsStore: sstore,
	}
}

func (sss *SearchSettingsServices) Get() (SearchSettings, error) {

	if err := sss.SearchSettingsStore.
		Where("id = 1").
		First(&sss.SearchSettings).
		Error; err != nil {
		return SearchSettings{}, err
	}

	return sss.SearchSettings, nil
}

func (sss *SearchSettingsServices) Upadate(
	amount uint, searchOn, addNew bool,
) error {
	sss.SearchSettings.Amount = amount
	sss.SearchSettings.SearchOn = searchOn
	sss.SearchSettings.AddNew = addNew

	tx := sss.SearchSettingsStore.
		Select("search_on", "add_new", "amount", "updated_at").
		Where("id = 1").
		Updates(&sss.SearchSettings)
	if tx.Error != nil {
		return fmt.Errorf("search not updated: %s", tx.Error)
	}

	return nil
}
