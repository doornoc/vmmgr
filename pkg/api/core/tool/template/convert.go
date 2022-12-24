package template

func GetArchStr(ty uint) string {
	typeStr := "x86_64"
	switch ty {
	case 1:
		typeStr = "i386"
		break
	}

	return typeStr
}

func GetExtensionStr(extension uint) string {
	// default is qcow2
	switch extension {
	case 1:
		return "raw"
	}
	return "qcow2"
}
