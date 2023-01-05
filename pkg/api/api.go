package api

import (
	"github.com/gin-gonic/gin"
	node "github.com/vmmgr/controller/pkg/api/core/controller/v1"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	vm "github.com/vmmgr/controller/pkg/api/core/vm/v1"
	"log"
	"net/http"
	"strconv"
)

func RestAPI() error {
	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/controller", node.ReceiveNode)
		}
	}

	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/vm", vm.GetWebSocketAdmin)
			// noVNC
			//v1.GET("/vnc/:access_token/:controller/:vm_uuid", wsVNC.GetByAdmin)
		}
	}

	go vm.HandleMessages()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Port), router))
	return nil
}

func cors(c *gin.Context) {

	//c.Header("Access-Control-Allow-Headers", "Accept, Content-ID, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-ID", "application/json")
	c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
