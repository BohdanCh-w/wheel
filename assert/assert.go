package assert

import (
	"fmt"
	"os"
	"runtime/debug"
)

var isProd = false

const (
	ModeProduction = iota
	ModeDevelopment
)

func SetMode(mode int) {
	isProd = mode == ModeProduction
}

func assert(condition bool, typ, msg string, args ...any) {
	if condition {
		return
	}

	args = append([]any{"msg", msg}, args...)

	fmt.Fprintf(os.Stderr, "ASSERT %s failed\n", typ)
	for i := 0; i < len(args)-1; i += 2 {
		fmt.Fprintf(os.Stderr, "   %s: %v\n", args[i], args[i+1])
	}

	if len(args)%2 == 1 {
		fmt.Fprintf(os.Stderr, "   %s\n", args[len(args)-1])
	}

	fmt.Fprintf(os.Stderr, "\nStack trace:\n%s\n", debug.Stack())
	os.Exit(1)
}

func Assert(truth bool, msg string, data ...any) {
	if isProd {
		return
	}

	assert(truth, "True", msg, data...)
}

func Nil(item any, msg string, data ...any) {
	if isProd {
		return
	}

	assert(item == nil, "Nil", msg, data...)
}

func NotNil(item any, msg string, data ...any) {
	if isProd {
		return
	}

	assert(item != nil, "Not Nil", msg, data...)
}

func Never(msg string, data ...any) {
	if isProd {
		return
	}

	assert(false, "Never", msg, data...)
}

func NoError(err error, msg string, data ...any) {
	if isProd {
		return
	}

	assert(err == nil, "No Error", msg, data...)
}
