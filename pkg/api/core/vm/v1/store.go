package v1

import (
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
	"time"
)

// Store
var Store = map[string]StoreStruct{}

type StoreStruct struct {
	StartTime     time.Time
	UpdateTime    time.Time
	Base          Base
	HostName      string
	VMInfo        vm.VMInput
	Template      template.Template
	ImageTemplate template.ImageTemplate
	DestBasePath  string
}
