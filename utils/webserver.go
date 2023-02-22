package utils

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type server struct {
	engine *gin.Engine
}

func RunServer(addr string) {
	var msgFinalMap sync.Map
	s := &server{
		engine: gin.Default(),
	}
	s.engine.Use()
	s.engine.GET("", get(&msgFinalMap))
	s.engine.POST("", set(&msgFinalMap))

	err := s.engine.Run(addr)
	if err != nil {
		panic(err)
	}
}

func get(msgFinalMap *sync.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		group := c.DefaultQuery("group", "")
		if group == "" {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "group is empty",
			})
			return
		}
		groupInt, err := strconv.Atoi(group)
		if err != nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "group is not int",
			})
			return
		}
		fmt.Println(groupInt)
		if msg, ok := msgFinalMap.Load(int64(groupInt)); ok {
			c.String(200, msg.(string))
			return
		} else {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "group is not exist",
			})
			return
		}
	}
}

func set(msgFinalMap *sync.Map) gin.HandlerFunc {
	return func(c *gin.Context) {
		type request struct {
			Group int64  `json:"group"`
			Msg   string `json:"msg"`
		}

		req := request{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "request body is not valid",
			})
			return
		}

		if req.Group == 0 {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "group is empty",
			})
			return
		}
		if req.Msg == "" {
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "msg is empty",
			})
			return
		}
		msgFinalMap.Store(req.Group, req.Msg)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
		})
	}
}
