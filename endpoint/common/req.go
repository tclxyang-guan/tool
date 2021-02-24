/**
* @Auther:gy
* @Date:2021/2/2 17:32
 */

package common

import (
	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, req interface{}) error {
	err := c.BindJSON(req)
	c.Set("req", req)
	return err
}
