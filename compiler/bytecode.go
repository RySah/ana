package compiler

import "math"

type BytecodeHandler struct {
	byteStream []byte
}

func NewBytecodeHandler(byteStream []byte) *BytecodeHandler {
	return &BytecodeHandler{byteStream: byteStream}
}

// Remove selected bits in the value
func removeBits(mask byte, value byte) byte {
	return (value | mask) ^ mask
}

const maxStackSize = math.MaxUint
const getStackSizeBitMask = 0b10000000

func (handler *BytecodeHandler) getStackSizes() (instrStackSize uint, paramStackSize uint, i byte) {
	var totalSize uint = 0
	i = 1

	instrStackSize = uint(handler.byteStream[0])
	if (instrStackSize >> 7) == 1 {
		totalSize = uint(removeBits(getStackSizeBitMask, byte(instrStackSize)))
		for (handler.byteStream[i] >> 7) == 1 {
			totalSize += uint(removeBits(getStackSizeBitMask, handler.byteStream[i]))
			i++
		}
		instrStackSize = totalSize
	} else {
		instrStackSize = uint(removeBits(getStackSizeBitMask, byte(instrStackSize)))
	}

	paramStackSize = uint(handler.byteStream[i])
	if (paramStackSize >> 7) == 1 {
		totalSize = uint(removeBits(getStackSizeBitMask, byte(paramStackSize)))
		i++
		for (handler.byteStream[i] >> 7) == 1 {
			totalSize += uint(removeBits(getStackSizeBitMask, handler.byteStream[i]))
			i++
		}
		paramStackSize = totalSize
	} else {
		paramStackSize = uint(removeBits(getStackSizeBitMask, byte(paramStackSize)))
	}

	return instrStackSize, paramStackSize, i
}

func (handler *BytecodeHandler) GetStacks() (instrStack Stack, paramStack Stack) {
	instrStackSize, paramStackSize, searched := handler.getStackSizes()

	var i uint = uint(searched)
	for i < instrStackSize+uint(searched) {
		instrStack.Push(handler.byteStream[i])
		i++
	}

	for i < paramStackSize+instrStackSize+uint(searched) {
		paramStack.Push(handler.byteStream[i])
		i++
	}

	return instrStack, paramStack
}
