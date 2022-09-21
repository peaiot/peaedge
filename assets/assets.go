package assets

import (
	"embed"
)

//go:embed testdata
var TestData embed.FS

//go:embed build.txt
var BuildInfo string

//go:embed buildver.txt
var BuildVer string
