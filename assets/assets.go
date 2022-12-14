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

//go:embed menu-admin.json
var AdminMenudata []byte

//go:embed menu-admin.json
var OprMenudata []byte

//go:embed funcs
var LuaFuncs embed.FS

//go:embed funcs/HandlerModbusData.lua
var HandlerModbusData string

//go:embed funcs/HandlerDataStream.lua
var HandlerDataStream string
