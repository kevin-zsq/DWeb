package main

import (
	"fmt"
	"net/http"
	"time"

	controller "github.com/kevin-zsq/DWeb/web"
)

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


func main() {
	engine := controller.NewEngine()
	engine.GET("/", func(c *controller.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	engine.GET("/hello/:name", func(c *controller.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
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
		v1.GET("/", func(c *controller.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *controller.Context) {
			// expect /hello?name=geektutu
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

	fmt.Println("Server start, listening port 9999")
	engine.Run(":9999")
}
