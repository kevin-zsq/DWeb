package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	controller "github.com/kevin-zsq/DWeb/web"
)

type student struct {
	Name string
	Age  int8
}

func onlyForV2() controller.HandleFunc {
	return func(c *controller.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.String(500, "Internal Server Error")
		// Calculate resolution time
		fmt.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	engine := controller.NewEngine()
	engine.Use(controller.Recovery())
	engine.Use(controller.Logger())
	// html render
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHTMLGlob("templates/*")
	// server static file
	engine.Static("/assets", "/Users/kevin/go/src/DWeb/static/css")

	// url/group register
	engine.GET("/hello/:name", func(c *controller.Context) {
		// expect /hello/kevin
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	//test internal error
	engine.GET("/panic", func(c *controller.Context) {
		// expect /hello/kevin
		names := []string{"kevin"}
		c.String(http.StatusOK, "hello %s\n", names[100])
	})

	engine.POST("/login", func(c *controller.Context) {
		c.JSON(http.StatusOK, controller.H{
			"user":     c.PostForm("user"),
			"passWord": c.PostForm("passWord"),
		})
	})

	engine.POST("/assets/*filepath", func(c *controller.Context) {
		c.JSON(http.StatusOK, controller.H{"filepath": c.Param("filepath")})
	})

	v1 := engine.Group("/v1")
	{
		v1.GET("/hello", func(c *controller.Context) {
			// expect /hello?name=kevin
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := engine.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *controller.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *controller.Context) {
			c.JSON(http.StatusOK, controller.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	// template render
	engine.GET("/", func(c *controller.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	engine.GET("/students", func(c *controller.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", controller.H{
			"title":  "students info",
			"stuArr": []student{{"kevin", 18}, {"sissi", 17}},
		})
	})
	engine.GET("/data", func(c *controller.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", controller.H{
			"title": "show time",
			"now":   time.Date(1994, 12, 20, 0, 0, 0, 0, time.UTC),
		})
	})
	fmt.Println("Server start, listening port 9999")
	engine.Run(":9999")
}
