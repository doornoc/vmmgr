package v1

import (
	"errors"
	"github.com/vmmgr/controller/pkg/api/core/tool/template"
	"github.com/vmmgr/controller/pkg/api/core/vm"
)

func CheckImageTemplate(imageTemplates []template.ImageTemplate, templateID string) (*template.ImageTemplate, error) {
	var imageTemplate *template.ImageTemplate = nil

	for _, v := range imageTemplates {
		if v.Disable {
			continue
		}
		if v.Name == templateID {
			imageTemplate = &v
			if imageTemplate != nil {
				break
			}
		}
	}
	if imageTemplate == nil {
		return imageTemplate, errors.New("template not found")
	}

	return imageTemplate, nil
}

func CheckImageTemplateSpec(imageTemplatesSpecPlans []template.SpecPlan, input vm.VMInput) error {
	for _, v := range imageTemplatesSpecPlans {
		if v.Disable {
			continue
		}
		if v.CPU == input.CPU && v.Memory == input.Memory {
			return nil
		}
	}
	return errors.New("template spec not found")
}

func CheckImageTemplateStorage(imageTemplatesStoragePlans []template.StoragePlan, name string, size uint) (*template.StoragePlan, error) {
	var storageTemplate *template.StoragePlan = nil

	for _, v := range imageTemplatesStoragePlans {
		if v.Disable {
			continue
		}
		if v.StorageID == name {
			for _, tmpSize := range v.Size {
				if tmpSize == size {
					storageTemplate = &v
					if storageTemplate != nil {
						break
					}
				}
			}
			if storageTemplate != nil {
				break
			}
		}
	}
	if storageTemplate == nil {
		return storageTemplate, errors.New("storage template not found")
	}

	return storageTemplate, nil
}
