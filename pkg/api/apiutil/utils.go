package apiutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	err := c.ShouldBindWith(obj, b)
	return err
}
