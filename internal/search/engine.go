package search

import (
	"fmt"
	"log"
	"time"

	"github.com/emarifer/search-engine/internal/services"
)

func loggerEndEngine(s time.Time) {
	endEngine := time.Now()
	log.Printf("🏁 Search engine crawl has finished at %v\n", endEngine.Sub(s))
}

func RunEngine(
	sss *services.SearchSettingsServices, us *services.UrlServices,
) {
	startEngine := time.Now()
	log.Println("🚀 Started search engine crawl…")

	// defer log.Println("Search engine crawl has finished")
	defer loggerEndEngine(startEngine)

	// Get crawl settings from DB
	_, err := sss.Get()
	if err != nil {
		fmt.Printf("something went wrong getting the settings: %s\n", err)

		return
	}

	// Check if search is turned on by checking settings
	if !sss.SearchSettings.SearchOn {
		fmt.Println("search is turned off")

		return
	}

	// Get next X urls to be tested
	nextUrls, err := us.GetNextCrawlUrls(sss.SearchSettings.Amount)
	if err != nil {
		fmt.Printf("something went wrong getting next urls: %s\n", err)

		return
	}

	newUrls := []services.CrawledUrl{}
	testedTime := time.Now()

	// Loop over the slice and run crawl on each url
	for _, next := range nextUrls {
		result := runCrawl(next.Url)

		// Check if the crawl was not successul
		if !result.Success {
			// Update row in database with the failed crawl
			err := us.UpdateUrl(services.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlBody.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlBody.PageTitle,
				PageDescription: result.CrawlBody.PageDescription,
				Headings:        result.CrawlBody.Headings,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Printf(
					"something went wrong updating a failed url: %s\n", err,
				)
			}

			continue
		}

		// Update a successful row in database
		err := us.UpdateUrl(services.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlBody.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlBody.PageTitle,
			PageDescription: result.CrawlBody.PageDescription,
			Headings:        result.CrawlBody.Headings,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Printf(
				"something went wrong updating %s\n", next.Url,
			)
		}

		// Push the newly found external urls to a slice
		for _, newUrl := range result.CrawlBody.Links.External {
			newUrls = append(newUrls, services.CrawledUrl{Url: newUrl})
		}
	} // End of loop

	// Check if we should add the newly found urls to the database
	if !sss.SearchSettings.AddNew {
		fmt.Println("Adding new urls to database is disabled")

		return
	}

	countNotAdded := 0
	// Insert newly found urls into database
	/* if err := us.SaveBatch(&newUrls); err != nil {
		fmt.Printf("something went wrong adding new urls to DB: %s\n", err)
	} */
	for _, newUrl := range newUrls {
		err := us.Save(&newUrl)
		if err != nil {
			countNotAdded++
			fmt.Printf(
				"something went wrong adding new url to database: %s\n",
				newUrl.Url,
			)
		}
	}

	fmt.Printf(
		"\nAdded %d new urls to database\n\n", len(newUrls)-countNotAdded,
	)
}

func RunIndex(us *services.UrlServices, sis *services.SearchIndexServices) {
	log.Println("🚀 Started search indexing…")

	defer log.Println("🏁 Search indexing has finished")

	// Get index settings from DB - Get all urls that are not indexed
	notIndexed, err := us.GetNotIndexed()
	if err != nil {
		fmt.Println("something went wrong getting the not indexed urls:", err)

		return
	}
	fmt.Println("not indexed urls:", len(notIndexed))

	// Create a new index map
	idx := make(Index)

	// Add the not indexed urls to the index
	idx.Add(notIndexed)

	// Save the index to DB
	err = sis.Save(idx, notIndexed)
	if err != nil {
		fmt.Println("something went wrong saving the index:", err)

		return
	}

	// Update the urls to be indexed=true
	err = us.SetIndexedTrue(notIndexed)
	if err != nil {
		fmt.Println("something went wrong updating the indexed urls:", err)

		return
	}
}
