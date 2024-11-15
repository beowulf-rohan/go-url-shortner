package cronjob

import (
	"log"
	"time"

	"github.com/beowulf-rohan/go-url-shortner/config"
	"github.com/beowulf-rohan/go-url-shortner/elasticsearch"
	"github.com/go-co-op/gocron"
)

const (
	DEFAULT_INDEX_CLEAR_INTERVAL = 30
)

func RunElasticClearUp() {
	config := config.GlobalConfig

	interval := config.IndexClearInterval
	if interval == 0 {
		interval = DEFAULT_INDEX_CLEAR_INTERVAL
	}

	cron := gocron.NewScheduler(time.UTC)
	cron.Every(interval).Hours().Do(func() {
		ElasticClient, err := elasticsearch.GetElasticClient(config.UrlMetadataIndex)
		if err != nil {
			return
		}

		err = ElasticClient.ClearExpiredDocuments()
		if err != nil {
			log.Println("error in monthly clearup", err)
			return
		}
	})
}
