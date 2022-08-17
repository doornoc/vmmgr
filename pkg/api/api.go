package api

import (
	"github.com/doornoc/vmmgr/pkg/api/core/tool/config"
	vm "github.com/doornoc/vmmgr/pkg/api/core/vm/v1"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func RestAPI() error {
	router := gin.Default()
	router.Use(cors)

	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/vm", vm.GetWebSocketAdmin)
			// noVNC
			//v1.GET("/vnc/:access_token/:node/:vm_uuid", wsVNC.GetByAdmin)
		}
	}

	go vm.HandleMessages()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Port), router))
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
