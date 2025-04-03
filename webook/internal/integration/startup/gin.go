package startup

import "github.com/gin-gonic/gin"

func InitGin() *gin.Engine {
	s := gin.New()
	return s
}
