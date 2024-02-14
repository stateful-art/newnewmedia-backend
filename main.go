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
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, ngrok-skip-browser-warning",
		AllowCredentials: true,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	music := app.Group("/music")
	places := app.Group("/places")

	placesroute.PlaceRoutes(places)
	musicroute.MusicRoutes(music)

	log.Fatal(app.Listen(":3000"))
}
