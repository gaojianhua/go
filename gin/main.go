package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 从path获取参数
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})
	// 从？后边获取参数
	r.GET("/name", func(c *gin.Context) {
		query := c.Query("name")
		//等价于：c.Request.URL.Query().Get("name")
		c.JSON(http.StatusOK, "Hello "+query)
	})

	// 获取请求体
	//r.POST("/stu", func(c *gin.Context) {
	//	body := c.Request.Body
	//	all, err := io.ReadAll(body)
	//	stu := Student{}
	//	json.Unmarshal(all, &stu)
	//	if err != nil {
	//		panic(err)
	//	}
	//	c.JSON(http.StatusOK, stu)
	//})

	r.POST("/stu", func(c *gin.Context) {
		stu := Student{}
		// c.ShouldBindJSON 使用了 c.Request.Body，不可重用。
		c.ShouldBindJSON(&stu)
		str := fmt.Sprintf("%+v", stu)
		fmt.Println(str)
		c.JSON(http.StatusOK, stu)
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST http://localhost:8080/admin -H 'authorization: Basic Zm9vOmJhcg==' -H 'content-type: application/json' -d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
