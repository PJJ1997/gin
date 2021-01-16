package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	ginCookie()
}

// 第一个简单实例
func ginFirstSimple() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		// 返回JSON
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}

// 不同的请求模式
func ginRequestPattern() {
	router := gin.Default()
	router.GET("/get", func(c *gin.Context) {
		// getInfo
	})
	router.DELETE("/delete", func(c *gin.Context) {
		// deleteInfo
	})
	router.POST("/post", func(c *gin.Context) {
		// postInfo
	})
	router.PUT("/put", func(c *gin.Context) {
		// putInfo
	})
	router.PATCH("/patch", func(c *gin.Context) {
		// patchInfo
	})
	router.HEAD("/head", func(c *gin.Context) {
		// headInfo
	})
	router.OPTIONS("/options", func(c *gin.Context) {
		// optionsInfo
	})
	router.Run(":3000") // 可以指定端口号
}

// 获取get请求restful参数
func ginGetRestfulParam() {
	router := gin.Default()
	// 命中 user/pengjj -> name = pengjj
	router.GET("user/:name", func(c *gin.Context) {
		name := c.Param("name")
		// 返回String
		c.String(http.StatusOK, "Hello %s", name)
	})

	// 命中 user/pengjj/eat -> name = pengjj && action = /eat
	router.GET("user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	router.Run()
}

// 获取get请求参数
func ginGetNormalParam() {
	router := gin.Default()

	router.GET("/welcome", func(c *gin.Context) {
		// 如果参数为空，赋予默认值
		firstName := c.DefaultQuery("firstName", "Guest")
		lastName := c.Query("lastName")
		c.JSON(http.StatusOK, gin.H{
			"message": gin.H{
				"firstName": firstName,
				"lastName":  lastName,
			},
		})
	})

	router.Run()
}

// 获取post请求参数
func ginPostParam() {
	router := gin.Default()
	router.POST("/post", func(c *gin.Context) {
		message := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("nick") // 需要结合前端html的form标签使用

		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    name,
			"page":    page,
		})
	})
	log.Println("Server started ...")
	router.Run()
}

// 上传单个文件
func ginUploadSingleFile() {
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file") // 需要结合前端form标签中的file标签使用
		if err != nil {
			c.String(http.StatusBadRequest, "文件上传出错")
		}
		log.Println(file.Filename)

		c.SaveUploadedFile(file, "/workspace/dev/files/")
		c.String(http.StatusOK, "'%s'上传成功", file.Filename)
	})
	router.Run()
}

// 上传多个文件
func ginUploadMultiFile() {
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "文件上传出错")
		}

		files := form.File["upload[]"]
		for _, file := range files {
			log.Println(file.Filename)
			c.SaveUploadedFile(file, "workspace/dev/files/")
		}
		c.String(http.StatusOK, "'%d'个文件已上传成功", len(files))
	})
}

// 访问路径分组
func ginGroup() {
	// gin.New()：无中间件
	// gin.Default()：带日志，捕获错误中间件
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/login", func(c *gin.Context) { log.Println("login ...") })
		v1.POST("/submit", func(c *gin.Context) { log.Println("submit ...") })
		v1.POST("/read", func(c *gin.Context) { log.Println("read ...") })
	}

	v2 := router.Group("/v2")
	{
		v2.GET("/login", func(c *gin.Context) { log.Println("login ...") })
		v2.GET("/submit", func(c *gin.Context) { log.Println("login ...") })
		v2.GET("/read", func(c *gin.Context) { log.Println("login ...") })
	}
	router.Run()
}

// 自定义日志打印
func ginCustomLog() {
	router := gin.New()
	// 将日志打印到日志文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout) 控制台和日志文件都打印到
	// 使用中间件
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s %s %s %s\n",
			param.ClientIP, param.TimeStamp.Format(time.RFC1123), param.Method, param.ErrorMessage)

	}))
	router.Use(gin.Recovery())

	router.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.Run()
}

// 模型绑定
func ginModelBind() {
	type Login struct {
		User     string `form:"user" json:"user"`
		Password string `form:"password" json:"password"`
	}
	router := gin.Default()
	router.GET("/loginJSON", func(c *gin.Context) {
		var json Login
		// ShouldBind：只绑定查询字符串的参数
		// ShouldBindQuery：GET请求绑定查询字符串的参数, POST请求绑定FORM表单的参数
		if err := c.ShouldBind(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println(json.User)
		log.Println(json.Password)
		if json.User != "pengjj" && json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "登录成功"})
	})
	router.Run()
}

// 模型绑定restful参数
func ginModelBindRestfulUri() {
	type Person struct {
		ID   string `uri:"id"`
		Name string `uri:"name"`
	}
	router := gin.Default()
	router.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "id": person.ID})
	})
	router.Run()
}

// 重定向
func ginRedirect() {
	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://baidu.com")
	})

	// 但uri不会变化
	router.GET("/test1", func(c *gin.Context) {
		c.Request.URL.Path = "/test2"
		router.HandleContext(c)
	})

	router.GET("/test2", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})
	router.Run()
}

// ====================== 运行多个服务 ======================
var eg errgroup.Group

func gin01() http.Handler {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Welcome Server01"})
		return
	})
	return router
}

func gin02() http.Handler {
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Welcome Server02"})
		return
	})
	return router
}

func ginStartAllServer() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      gin01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      gin02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	eg.Go(func() error {
		return server01.ListenAndServe()
	})
	eg.Go(func() error {
		return server02.ListenAndServe()
	})
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

// 获取或设置Cookie
func ginCookie() {
	router := gin.Default()
	router.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			log.Println(err)
			cookie = "pengjj"
			c.SetCookie("gin_cookie", "test", 3600, "/", "127.0.0.1", false, true)
		}
		// 无需再调一次Cookie()来获取
		log.Println("cookie: ", cookie)
	})
	router.Run()
}

// 自定义中间件：其实就是一个函数，可以在handle前后做相应的处理
// https://docs.lvrui.io/2019/01/27/gin%E6%B7%BB%E5%8A%A0%E8%87%AA%E5%AE%9A%E4%B9%89%E4%B8%AD%E9%97%B4%E4%BB%B6/
