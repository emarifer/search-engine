package utils

import (
	"fmt"

	"github.com/emarifer/search-engine/internal/search"
	"github.com/emarifer/search-engine/internal/services"
	"github.com/robfig/cron/v3"
)

func StartCronJobs(
	sss *services.SearchSettingsServices,
	us *services.UrlServices,
	sis *services.SearchIndexServices,
) {
	c := cron.New()
	// "0 * * * *" Run every hour
	c.AddFunc("35 * * * *", func() { search.RunEngine(sss, us) }) // @every 120s
	// "" Run every hour at 15 minutes past
	c.AddFunc("37 * * * *", func() { search.RunIndex(us, sis) }) // @every 30s

	// c.AddFunc("38 * * * *", func() { search.RunEngine(sss, us) })
	// c.AddFunc("40 * * * *", func() { search.RunIndex(us, sis) })

	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", cronCount)
}
