package background

import (
	"time"

	"ramdeuter.org/solscraper/db"
	"ramdeuter.org/solscraper/scraper"
)

const scrapeInterval = 1

/**
fetch jobs every 2 minutes
run each job
*/

func StartScraping(d db.DB) {
	for {
		select {
		case <-time.After(scrapeInterval * time.Minute):
			startJobs(d)
		}
	}

}

func startJobs(d db.DB) {
	metadata, err := d.GetMetadata()
	if err != nil {
		return
	}
	for _, value := range metadata {
		q := value.Query
		if err != nil {
			continue
		}
		data := scraper.ScrapeData(q)
		d.SaveData(value.Name, data)
	}
}
