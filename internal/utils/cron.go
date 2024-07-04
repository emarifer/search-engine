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
	c.AddFunc("0 * * * *", func() { search.RunEngine(sss, us) }) // @every 120s
	// "05 * * * *" Run every hour at 5 minutes past
	c.AddFunc("05 * * * *", func() { search.RunIndex(us, sis) }) // @every 30s

	// c.AddFunc("29 * * * *", func() { search.RunEngine(sss, us) })
	// c.AddFunc("30 * * * *", func() { search.RunIndex(us, sis) })
	// c.AddFunc("31 * * * *", func() { search.RunEngine(sss, us) })
	// c.AddFunc("32 * * * *", func() { search.RunIndex(us, sis) })
	// c.AddFunc("33 * * * *", func() { search.RunEngine(sss, us) })
	// c.AddFunc("34 * * * *", func() { search.RunIndex(us, sis) })

	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", cronCount)
}
