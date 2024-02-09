package main

import (
	"bufio"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"strings"
)

type entry struct {
	key   string
	entry string
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("/")
		return c.SendString("Hello, World!")
	})
	app.Post("/api/upload", func(c *fiber.Ctx) error {
		log.Info("/api/upload")
		var envContainer []entry
		file, err := c.FormFile("envFile")
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
		}
		dat, err := file.Open()
		if err != nil {
			log.Error(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": err.Error()})
		}
		defer dat.Close()
		scanner := bufio.NewScanner(dat)
		for scanner.Scan() {
			fmt.Printf("env var: %s\n", scanner.Text())
			slicedEnvEntryLine := strings.Split(scanner.Text(), "=")
			if len(slicedEnvEntryLine) > 1 && string(slicedEnvEntryLine[0][0]) != "#" {
				envContainer = append(envContainer, entry{key: slicedEnvEntryLine[0], entry: slicedEnvEntryLine[1]})
			}
		}
		log.Info(envContainer)
		if len(envContainer) != 0 {
			return c.JSON(fiber.Map{"error": false, "msg": nil, "info": fmt.Sprintf("added %d variables", len(envContainer))})
		} else {
			return c.JSON(fiber.Map{"error": true, "msg": "no environment variables added"})
		}
	})
	app.Listen(":3002")
}
