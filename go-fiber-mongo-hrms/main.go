package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// autoload .env
	_ "github.com/joho/godotenv/autoload"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "go-fiber-mongo-hrms"

type Employee struct {
	Id     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func Connect() error {
	mongoURI := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {
		query := bson.D{}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(employees)
	})

	app.Post("employee", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		employee.Id = "" // reset id
		res, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		filter := bson.D{{Key: "_id", Value: res.InsertedID}}
		var createdEmployee Employee
		err = collection.FindOne(c.Context(), filter).Decode(&createdEmployee)
		if err == mongo.ErrNoDocuments {
			return fiber.ErrNotFound
		} else if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(201).JSON(employee)
	})

	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		employeeId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		update := bson.D{{
			Key: "$set", Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "age", Value: employee.Age},
				{Key: "salary", Value: employee.Salary},
			}}}

		err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err == mongo.ErrNoDocuments {
			return fiber.ErrNotFound
		} else if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(employee)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {

		employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		res, err := mg.Db.Collection("employees").DeleteOne(c.Context(), &query)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if res.DeletedCount < 1 {
			return fiber.ErrNotFound
		}

		return c.JSON(true)
	})

	fmt.Println("Starting server at http://localhost:8000")
	log.Fatal(app.Listen("127.0.0.1:8000"))
}
