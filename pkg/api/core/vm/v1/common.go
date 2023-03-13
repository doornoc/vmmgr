package v1

import (
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"log"
	"time"
)

func (b *Base) Error(error string) {
	delete(Store, b.UUID)
	log.Printf("[%s_%d] %s\n", b.UUID, b.Type, error)
	vm.ClientBroadcast <- vm.WebSocketResult{
		UUID:      b.UUID,
		Type:      b.Type,
		Err:       error,
		CreatedAt: time.Now(),
	}
}
