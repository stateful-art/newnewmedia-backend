package service

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/nats-io/nats.go"
)

type SearchService struct {
	elasticClient *elasticsearch.Client
	natsClient    *nats.Conn
}

func NewSearchService(elasticClient *elasticsearch.Client, natsClient *nats.Conn) *SearchService {
	return &SearchService{elasticClient: elasticClient, natsClient: natsClient}
}
