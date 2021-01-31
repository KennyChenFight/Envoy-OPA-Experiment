package main

import (
	"github.com/gofiber/fiber/v2"
	admissionv1 "k8s.io/api/admission/v1"
)


func main() {
	app := fiber.New()
	app.Post("/mutate", func(c *fiber.Ctx) error {
		resp := admissionv1.AdmissionResponse{
			Allowed:          true,
			// todo
			Patch:            nil,
			PatchType:  func() *admissionv1.PatchType {
				pt := admissionv1.PatchTypeJSONPatch
				return &pt
			}(),
		}
		return c.JSON(resp)
	})
}