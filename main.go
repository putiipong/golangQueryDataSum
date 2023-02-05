package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/test")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Test!")
	})

	app.Get("/api/create", func(ctx *fiber.Ctx) error {
		_, errTable := db.Query("CREATE TABLE data (id int NOT NULL PRIMARY KEY AUTO_INCREMENT, active_power INT,power_input INT);")
		if errTable != nil {
			return ctx.JSON(fiber.Map{"error_create_table": errTable.Error()})
		}
		for i := 1; i <= 100; i++ {

			_, errData := db.Exec("INSERT INTO data (active_power,power_input) VALUES (floor(rand() * 1000 + 1), floor(rand() * 1000 + 1));")
			if errData != nil {
				return ctx.JSON(fiber.Map{"error_insert_table": errData.Error()})
			}
		}

		return ctx.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/api/dataSum", func(ctx *fiber.Ctx) error {
		param := ctx.Query("sum")
		var result int
		query := fmt.Sprintf("SELECT SUM(%v) FROM data;", param)
		rows, err := db.Query(query)

		if err != nil {
			ctx.Status(500)
			return ctx.JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&result); err != nil {
				return ctx.JSON(fiber.Map{"error_scan": err.Error()})
			}
		}
		return ctx.JSON(fiber.Map{"sum_value": result})
	})

	app.Listen(":4000")

}
