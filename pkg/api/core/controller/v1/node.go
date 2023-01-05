package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	vmInterface "github.com/vmmgr/controller/pkg/api/core/vm"
	vm "github.com/vmmgr/controller/pkg/api/core/vm/v1"
	"net/http"
	"strconv"
	"time"
)

func ReceiveNode(c *gin.Context) {
	var node controller.Node
	if err := c.BindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, "")
	}

	// when copy finish, send result to controller
	if node.Progress == 100 && node.Status && node.Comment == "success" {
		go vm.CreateForCloudInit(node.UUID)
	}

	if node.Progress < 100 && node.Status {
		storeVM := vm.Store[node.UUID]
		storeVM.UpdateTime = time.Now()
		vm.Store[node.UUID] = storeVM
		vmInterface.ClientBroadcast <- vmInterface.WebSocketResult{
			Type:      vmInterface.MessageTypeCreateVM,
			CreatedAt: time.Now(),
			UUID:      node.UUID,
			VMDetail: []vmInterface.VMDetail{{
				Node: storeVM.HostName,
			}},
			Data: map[string]string{"create_progress": "10", "copy_progress": strconv.Itoa(int(node.Progress))},
		}
	}

	c.JSON(http.StatusOK, "")
}
