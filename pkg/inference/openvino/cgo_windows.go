package openvino

// Windows：要用 MinGW 的 gcc 編譯；標頭檔與 openvino_c.dll 的位置用
// CGO_CFLAGS="-I<openvino>/runtime/include" CGO_LDFLAGS="-L<openvino>/runtime/bin/intel64/Release" 指定，
// 執行時 openvino_c.dll 所在目錄要在 PATH。

/*
#cgo LDFLAGS: -lopenvino_c
*/
import "C"
