package main

import (
	"encoding/json"
	"github.com/KennyChenFight/Envoy-OPA-Experiment/k8s-envoy-opa/jwt-example/service/model"
	"github.com/KennyChenFight/Envoy-OPA-Experiment/k8s-envoy-opa/jwt-example/service/store"
	"github.com/gofiber/fiber/v2"
)

func main() {
	memoryStore := store.NewMemoryStore()

	app := fiber.New()

	app.Get("/people", func(c *fiber.Ctx) error {
		return c.JSON(memoryStore.GetAll())
	})

	app.Post("/people", func(c *fiber.Ctx) error {
		var p model.Person
		if err := json.Unmarshal(c.Body(), &p); err != nil {
			return c.JSON(fiber.Map{"error": "parse body fail"})
		}
		id := memoryStore.Save(p)
		return c.JSON(fiber.Map{"id": id})
	})

	app.Listen(":8080")
}
