package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/zzacong/eleven-go/go-fiber-crm/database"
	"github.com/zzacong/eleven-go/go-fiber-crm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead", models.GetLeads)
	app.Get("/api/v1/lead/:id", models.GetLead)
	app.Post("/api/v1/lead", models.NewLead)
	app.Delete("/api/v1/lead/:id", models.DeleteLead)
}

func InitDatabase() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("leads.db"))
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection opened to database")
	database.DB.AutoMigrate((&models.Lead{}))
	fmt.Println("Database migrated")
}

func main() {
	app := fiber.New()
	InitDatabase()
	SetupRoutes(app)
	log.Fatal(app.Listen("127.0.0.1:8000"))
	// defer database.DB.Close()
}
