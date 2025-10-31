package traceLog

import (
	"bytes"
	"context"
	"runtime"
	"strconv"
	"sync"
	"unsafe"
)

var contextMap sync.Map

var gidCache sync.Map

func getGID() uint64 {
	var x int
	ptr := unsafe.Pointer(&x)
	
	if v, ok := gidCache.Load(ptr); ok {
		return v.(uint64)
	}
	
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := bytes.Fields(buf[:n])
	if len(fields) < 2 {
		return 0
	}
	id, _ := strconv.ParseUint(string(fields[1]), 10, 64)
	
	gidCache.Store(ptr, id)
	return id
}

func SetGoroutineContext(ctx context.Context) {
	id := getGID()
	contextMap.Store(id, ctx)
}

func GetGoroutineContext() context.Context {
	id := getGID()
	if v, ok := contextMap.Load(id); ok {
		if ctx, ok := v.(context.Context); ok {
			return ctx
		}
	}
	return context.Background()
}

func DelGoroutineContext() {
	id := getGID()
	contextMap.Delete(id)
}
