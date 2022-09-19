package assets

import (
	"embed"
)

//go:embed resources
var Resources embed.FS

//go:embed build.txt
var BuildInfo string

//go:embed buildver.txt
var BuildVer string
