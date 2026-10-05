package openvino

// Linux：標頭檔與函式庫要在編譯器預設搜尋路徑；裝在別的地方（例如 /opt/intel/openvino）就設
// CGO_CFLAGS="-I<include>" CGO_LDFLAGS="-L<lib> -Wl,-rpath,<lib>"。

/*
#cgo LDFLAGS: -lopenvino_c
*/
import "C"
