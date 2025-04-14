package controller

import (
	"com.bookstore/demo/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// 获取书架
func GetShelfHandler(c *gin.Context) {
	// 获取参数信息
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	// 检查参数是否合法
	if id <= 0 {
		zap.L().Error("invalid id", zap.Int64("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "invalid id"})
		return
	}
	// 安装参数到数据库中查询书架
	data, err := logic.GetShelf(id)
	if err != nil {
		zap.L().Error("get shelf failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// 获取书架列表
func ShelvesListHandler(c *gin.Context) {}

// 创建书架
func CreateShelfHandler(c *gin.Context) {}

// 删除书架
func DelateShelfHandler(c *gin.Context) {}
