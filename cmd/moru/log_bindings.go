package main

// #include <stdlib.h>
//
// typedef void (*log_write_t)(const void* msg, int len);
// void bridge_log_write(log_write_t f, const void* msg, int len);
import "C"

type nativeLogger struct {
	log C.log_write_t
}

func newNativeLogger() *nativeLogger {
	return &nativeLogger{}
}

func (n *nativeLogger) Bind(log C.log_write_t) {
	n.log = log
}

func (n *nativeLogger) Write(bs []byte) (int, error) {
	if n.log == nil {
		return len(bs), nil
	}

	if len(bs) == 0 {
		return 0, nil
	}

	// append a zero byte to the end of the buffer to make it a C string
	bs = append(bs, 0)

	p := C.CBytes(bs)
	// we cannot free the memory when this function returns, since
	// the memory is still being used asynchronously by the native code.

	C.bridge_log_write(n.log, p, C.int(len(bs)))
	return len(bs), nil
}
