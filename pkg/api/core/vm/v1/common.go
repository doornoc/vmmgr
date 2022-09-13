package v1

import (
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"time"
)

func (b *Base) Error(error string) {
	vm.ClientBroadcast <- vm.WebSocketResult{
		UUID:      b.UUID,
		Type:      b.Type,
		Err:       error,
		CreatedAt: time.Now(),
	}

}
