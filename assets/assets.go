package assets

import (
	"embed"
)

//go:embed static
var StaticFs embed.FS

//go:embed templates
var TemplatesFs embed.FS

//go:embed testdata
var TestData embed.FS

//go:embed build.txt
var BuildInfo string

//go:embed buildver.txt
var BuildVer string

//go:embed funcs
var LuaFuncs embed.FS

//go:embed menu-admin.json
var AdminMenudata []byte

//go:embed menu-admin.json
var OprMenudata []byte
