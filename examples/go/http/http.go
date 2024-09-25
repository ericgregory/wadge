//go:generate go run github.com/bytecodealliance/wasm-tools-go/cmd/wit-bindgen-go generate -w app -o bindings ./wit

package main

import (
	"unsafe"

	"github.com/bytecodealliance/wasm-tools-go/cm"
	incominghandler "github.com/wasmCloud/wadge/examples/go/http/bindings/wasi/http/incoming-handler"
	"github.com/wasmCloud/wadge/examples/go/http/bindings/wasi/http/types"
)

func init() {
	incominghandler.Exports.Handle = func(request types.IncomingRequest, responseOut types.ResponseOutparam) {
		if err := handle(request, responseOut); err != nil {
			types.ResponseOutparamSet(responseOut, cm.Err[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](*err))
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}

func handle(req types.IncomingRequest, out types.ResponseOutparam) *types.ErrorCode {
	resp := types.NewOutgoingResponse(req.Headers())

	body := resp.Body()
	if body.IsErr() {
		return ptr(types.ErrorCodeInternalError(cm.Some("failed to get response body")))
	}
	bodyOut := body.OK()

	bodyWrite := bodyOut.Write()
	if bodyWrite.IsErr() {
		return ptr(types.ErrorCodeInternalError(cm.Some("failed to get response body stream")))
	}

	types.ResponseOutparamSet(out, cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](resp))
	stream := bodyWrite.OK()
	s := "hello world"
	writeRes := stream.BlockingWriteAndFlush(cm.NewList(unsafe.StringData(s), uint(len(s))))
	if writeRes.IsErr() {
		return nil
	}
	stream.ResourceDrop()

	finishRes := types.OutgoingBodyFinish(*bodyOut, cm.None[types.Fields]())
	if finishRes.IsErr() {
		return nil
	}
	return nil
}

func main() {}
