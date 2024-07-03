package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host         = "localhost"
	port         = 5432
	databaseName = "postgres"
	username     = "postgres"
	password     = "mypassword"
)

type Task struct {
	Id         int    `json:"id"`
	User_Id    int    `json:"user_id"`
	Name       string `json:"task_name"`
	CreateTime string `json:"created_at"`
	UpdateTime string `json:"updated_at"`
	Status     bool   `json:"status"`
}

// for create data
type CreateTask struct {
	User_Id int    `json:"user_id"`
	Name    string `json:"task_name"`
}

// for update data
type UpdateTask struct {
	Name   string `json:"task_name"`
	Status bool   `json:"status"`
}

var db *sql.DB

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, databaseName)

	sdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	// เชื่อมเสร็จโดยไม่มี error ก็ให้ค่าที่เชื่อมผ่านไปเป็นค่า global
	defer db.Close()

	db = sdb

	// db.Ping() คือคำสั่งเช็คว่าต่อได้จริงไหม
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	// Conntection Database Successful รันอันนี้ได้แปลว่าไม่เจอ error
	print("Conntection Database Successful\n")

	// Test Server
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})
	// <Create> Create data
	app.Post("/createTask", createTaskHandler)
	// <Read> Get data params
	app.Get("/getTask/:id", getTaskHandler)
	// <Update> Update data
	app.Put("/updateTask/:id", updateTaskHandler)
	// <Delete> Delete Data
	app.Delete("/deleteTask/:id", deleteTaskHandler)
	// <Read> Get many data params
	app.Get("/getTasks/:id", getTasksHandler)

	app.Listen(":8080")

	// tasks, err := getTasks()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(tasks)

}

func getTaskHandler(c *fiber.Ctx) error {
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	task, err := getTask(taskId)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return c.JSON(task)
}

func createTaskHandler(c *fiber.Ctx) error {
	taskData := new(CreateTask)
	if err := c.BodyParser(taskData); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	if err := createTask(taskData); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return c.JSON(taskData)
}

func updateTaskHandler(c *fiber.Ctx) error {
	updateTaskData := new(UpdateTask)
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	if err := c.BodyParser(updateTaskData); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	err = updateTask(taskId, updateTaskData)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return c.JSON(updateTaskData)
}

func deleteTaskHandler(c *fiber.Ctx) error {
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	err = deleteTask(taskId)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)

}

func getTasksHandler(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	task, err := getTasks(userId)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}
	return c.JSON(task)
}
