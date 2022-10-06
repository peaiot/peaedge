package goscript

import (
	"github.com/goplus/igop"
	_ "github.com/goplus/igop"
	_ "github.com/goplus/igop/gopbuild"
	_ "github.com/goplus/igop/pkg/bytes"
	_ "github.com/goplus/igop/pkg/errors"
	_ "github.com/goplus/igop/pkg/fmt"
	_ "github.com/goplus/igop/pkg/math"
	_ "github.com/goplus/igop/pkg/os"
	_ "github.com/goplus/igop/pkg/path/filepath"
	_ "github.com/goplus/igop/pkg/reflect"
	_ "github.com/goplus/igop/pkg/runtime"
	_ "github.com/goplus/igop/pkg/strings"
	_ "github.com/goplus/igop/pkg/sync"
	_ "github.com/goplus/igop/pkg/time"
)

func RunFunc(filename, source string, funcName string, args ...interface{}) (interface{}, error) {
	ctx := igop.NewContext(0)
	p, err := ctx.LoadFile(filename, source)
	if err != nil {
		return nil, err
	}

	v, err := ctx.RunFunc(p, funcName, args...)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func RunScript(filename, source string) (exitCode int, err error) {
	return igop.RunFile(filename, source, nil, 0)
}
