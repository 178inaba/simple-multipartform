package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":1323"))
}

func upload(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, f := range files {
		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(f.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files with fields name=%s and email=%s.</p>", len(files), name, email))
}
