package beego

import (
	"github.com/beego/beego/v2/server/web"
	"testing"
)

func TestUserController(t *testing.T)  {
	/*
	go func() {
		s := web.NewHttpSever()
		s.Run(":8082")
	}()
	 */


	web.BConfig.CopyRequestBody = true
	c := &UserController{}
	// get:GetUser --> get: request method, GetUser: function
	web.Router("/user", c, "get:GetUser")
	// listen 8081 port
	web.Run(":8081")
}
