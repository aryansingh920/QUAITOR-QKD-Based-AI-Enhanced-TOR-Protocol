package controllers

import (
	"log"
	"text/template"

	"github.com/gofiber/fiber/v2"
)

var tmpl *template.Template

func init() {
	// Parse all templates into a single root template
	tmpl = template.Must(template.New("").ParseGlob("templates/**/*.html"))

	// Debugging: Log all loaded templates
	for _, t := range tmpl.Templates() {
		log.Println("Loaded template during init:", t.Name())
	}
}

func renderTemplateToFiber(ctx *fiber.Ctx, tmplName string, data interface{}) error {
	// Execute the specified template
	err := tmpl.ExecuteTemplate(ctx.Context().Response.BodyWriter(), tmplName, data)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Template rendering error: " + err.Error())
	}
	return nil
}

func HomeHandler(ctx *fiber.Ctx) error {
	// Debugging: Log all loaded templates
	for _, t := range tmpl.Templates() {
		log.Println("Loaded template:", t.Name())
	}

	data := map[string]interface{}{
		"Title": "Home Page",
		"Name":  "Aryan",
	}
	return renderTemplateToFiber(ctx, "layout", data)
}
