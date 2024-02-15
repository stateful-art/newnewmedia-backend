package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	db "newnewmedia.com/db"
	musicroute "newnewmedia.com/microservices/music/routes"
	placesroute "newnewmedia.com/microservices/place/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	app := fiber.New()
	db.ConnectDB()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == "https://www.newnewmedia.com" || origin == "http://localhost:5173"
		},
		AllowOrigins:     "http://localhost",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
	}))
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	music := app.Group("/music")
	places := app.Group("/places")

	placesroute.PlaceRoutes(places)
	musicroute.MusicRoutes(music)

	log.Fatal(app.Listen(":3000"))
}
