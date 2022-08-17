package v1

import (
	"github.com/doornoc/vmmgr/pkg/api/core/vm"
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
