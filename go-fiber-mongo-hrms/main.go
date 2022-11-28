package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

type Employee struct {
	Id     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float32 `json:"salary"`
	Age    int     `json:"age"`
}

const dbName = "fiber-hrms"
const mongoURL = "mongodb://localhost:27017/" + dbName

var mg MongoInstance

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err = client.Connect(ctx)
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
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), bson.D{{}})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var employees []Employee
		if err = cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(employees)
	})

	app.Post("/employee", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")

		newEmployee := new(Employee)
		if err := c.BodyParser(newEmployee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		newEmployee.Id = ""
		insertOneResult, err := collection.InsertOne(c.Context(), newEmployee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertOneResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := new(Employee)
		createdRecord.Decode(createdEmployee)

		return c.Status(201).JSON(createdEmployee)
	})

	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		employeeId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		employee := new(Employee)
		if err = c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		update := bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "age", Value: employee.Age},
					{Key: "salary", Value: employee.Salary},
				},
			},
		}

		err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(404)
			}
			return c.Status(500).SendString(err.Error())
		}

		employee.Id = idParam
		return c.Status(201).JSON(employee)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		employeeId, err := primitive.ObjectIDFromHex(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		result, err := mg.Db.Collection("employees").DeleteOne(c.Context(), query)

		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).SendString("Record Delete")
	})

	log.Fatal(app.Listen(":3000"))
}
