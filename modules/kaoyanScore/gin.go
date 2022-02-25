package kaoyanScore

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type server struct {
	engine *gin.Engine
}

func runGin() {
	s := &server{
		engine: gin.Default(),
	}
	s.engine.Use()
	s.engine.GET("/score", func(c *gin.Context) {
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
	})

	err := s.engine.Run(":8090")
	if err != nil {
		panic(err)
	}
}
