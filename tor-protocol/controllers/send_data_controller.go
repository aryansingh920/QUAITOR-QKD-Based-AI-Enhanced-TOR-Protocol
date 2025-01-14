package controllers

import (
	"log"
	"strings"
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

	// log.Println("Rendering Home Page")
	// log.Println("Request URL:", ctx.OriginalURL())
	// log.Println("Request Method:", ctx.Method())
	// log.Println("Request Path:", ctx.Path())
	// log.Println("Request Query:", ctx.Query())
	// log.Println("Request Body:", ctx.Body())

	//path should contain .onion domain else it will return 404
	// if !strings.Contains(ctx.Path(), ".onion") {
	// 	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"Error": "404 Not Found"})
	// }

	//remove / from the path
	// ctx.Path()
	// path := strings.TrimPrefix(ctx.Path(), "/")
	msgParam := ctx.Query("msg","")
	entryParam := ctx.Query("entry","")



	data := map[string]interface{}{
		"Title": strings.TrimPrefix(ctx.Path(), "/"),
		"Name":  "Aryan",
		"Msg": msgParam,
		"Entry": entryParam,
	}
	return renderTemplateToFiber(ctx, "layout", data)
}
