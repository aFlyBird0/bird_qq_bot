package kaoyanScore

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type server struct {
	engine *gin.Engine
}

// RunServer 运行一个local webserver，展示考研信息
func RunServer(addr string, msgFinalMap *sync.Map) {
	s := &server{
		engine: gin.Default(),
	}
	s.engine.Use()
	s.engine.GET("", get(msgFinalMap))

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
