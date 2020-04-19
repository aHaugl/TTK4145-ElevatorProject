package Hardware

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func io_init() int {
	return int(C.io_init())
}

func io_setBit(channel int) {
	C.io_set_bit(C.int(channel))
}

func io_clearBit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func io_writeAnalog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}

func io_readBit(channel int) int {
	return int(C.io_read_bit(C.int(channel)))
}

func io_readAnalog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}
