package elastic

import (
	"log"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func New() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Panic(err)
	}
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	return es
}
