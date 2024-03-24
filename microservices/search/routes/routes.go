package searchroutes

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	controller "newnew.media/microservices/search/controller"
)

func SearchRoutes(app fiber.Router, elasticClient *elasticsearch.Client) {
	// Pass the PlaylistService instance to the PlaylistController
	controller := controller.NewElasticIndexController(elasticClient)

	// Admin Index routes
	app.Post("/admin/create-index/:indexName", controller.CreateIndex)
	app.Delete("/admin/delete-index/:indexName", controller.DeleteIndex)

	// Place routes
	app.Post("/index-place", controller.IndexPlace)
	app.Get("/place", controller.SearchPlaceByName)
	// app.Get("/place/prefix", controller.PrefixSearchPlaceByName)

}
