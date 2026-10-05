// Package openvino 用 OpenVINO 的 C API 直接讀 ultralytics 匯出的 OpenVINO 模型（.xml 加 .bin）。
// 所有平台都會編進來，編譯機器要先裝好 OpenVINO C API（macOS: brew install openvino）。
package openvino

// cgo 是怎麼作用的：
//
//  1. 編譯時，Go 掃描有 import "C" 的檔案。import "C" 正上方緊貼的註解就是一小段 C 程式（下面那段 #include），
//     cgo 讀它就知道有哪些 C 的型別與函式可以用，例如 C.ov_core_t、C.ov_core_create。
//  2. cgo 會自動產生轉接層：Go 呼叫 C.ov_xxx，先進轉接函式，再由它呼叫真正的 C 函式。
//  3. #cgo 指令告訴編譯器去哪找東西，寫在 cgo_darwin.go、cgo_linux.go、cgo_windows.go（各系統不同）：
//     CFLAGS  -I<目錄>   編譯 C 程式時，去這個目錄找 openvino.h
//     LDFLAGS -lopenvino_c  連結時，把 libopenvino_c 連進來
//  4. cgo 把產生的 C 程式交給系統的 C 編譯器（clang / gcc）編譯，再和 Go 程式連結成一個執行檔。
//     所以編譯的機器要有 C 編譯器，也要有 OpenVINO 的標頭檔與函式庫，缺一樣都會連結失敗。
//  5. 執行時，作業系統載入 libopenvino_c（.dylib / .so / .dll），Go 呼叫 C.ov_xxx 就進到 OpenVINO 的
//     C 函式，由它讀 .xml 與 .bin、用 CPU 推論，再把結果交回 Go。
//
//     你的 Go 程式 → cgo 轉接層 → libopenvino_c（C/C++）→ 推論結果 → 回到 Go
//
// 各作業系統的分岔：
//
//	                         go build
//	                            │
//	            Go 看 GOOS（目標系統），只挑一個 cgo_<系統>.go
//	                            │
//	        ┌───────────────────┼───────────────────┐
//	        ▼                   ▼                   ▼
//	  GOOS=darwin         GOOS=linux          GOOS=windows
//	  cgo_darwin.go       cgo_linux.go        cgo_windows.go
//	        │                   │                   │
//	  -I/opt/homebrew     只有 -lopenvino_c   只有 -lopenvino_c
//	  /include            （標頭檔在系統      （路徑靠 CGO_CFLAGS /
//	  -L/opt/homebrew      預設位置，或靠      CGO_LDFLAGS 指定）
//	  /lib -rpath          CGO_* 指定）             │
//	        │                   │                   │
//	  clang               gcc                 MinGW gcc
//	        │                   │                   │
//	        └───────────────────┼───────────────────┘
//	                            ▼
//	              model.go（三個系統都一樣）
//	          C.ov_xxx 呼叫 OpenVINO，邏輯完全相同
//	                            │
//	                            ▼
//	              執行檔 → 執行時載入 libopenvino_c
//	              （.dylib / .so / .dll，各系統不同）
//
//	三個系統唯一的差別是：標頭檔在哪、函式庫在哪、用哪個 C 編譯器，
//	也就是 cgo_<系統>.go 裡那幾行 #cgo。
//
// 注意：這段說明和下面的 C 程式之間要空一行；註解緊貼 import "C" 的話，會被當成 C 程式。

/*
#include <stdlib.h>
#include <openvino/c/openvino.h>
*/
import "C"

import (
	"fmt"
	"slices"
	"sync"
	"unsafe"
)

// Model 是一個編譯好的 OpenVINO 模型，只有一個輸入、一個輸出，輸入是固定的 [1, 3, H, W]、float32。
// 一個 Model 只有一個 infer request，Run 會互斥，要平行推論就多載入幾個 Model。
type Model struct {
	mutex sync.Mutex

	core     *C.ov_core_t
	model    *C.ov_model_t
	compiled *C.ov_compiled_model_t
	request  *C.ov_infer_request_t

	// input 的記憶體是 OpenVINO 配的，每次 Run 把新資料複製進去，不把 Go 的記憶體交給 C 保留。
	input     *C.ov_tensor_t
	inputData unsafe.Pointer
	inputSize int

	height int
	width  int
}

// Load 讀 sXmlPath（.bin 預設在同目錄、同檔名），編譯到 sDevice（例如 "CPU"、"GPU"、"NPU"）。
func Load(sXmlPath string, sDevice string) (*Model, error) {
	oModel := &Model{}
	if err := oModel.load(sXmlPath, sDevice); err != nil {
		_ = oModel.Close()
		return nil, fmt.Errorf("openvino %s: %w", sXmlPath, err)
	}
	return oModel, nil
}

func (oSelf *Model) load(sXmlPath string, sDevice string) error {
	if err := check(C.ov_core_create(&oSelf.core), "create core"); err != nil {
		return err
	}

	cXmlPath := C.CString(sXmlPath)
	defer C.free(unsafe.Pointer(cXmlPath))
	if err := check(C.ov_core_read_model(oSelf.core, cXmlPath, nil, &oSelf.model), "read model"); err != nil {
		return err
	}

	aInputShape, err := inputShape(oSelf.model)
	if err != nil {
		return err
	}
	if len(aInputShape) != 4 || aInputShape[0] != 1 || aInputShape[1] != 3 || aInputShape[2] <= 0 || aInputShape[3] <= 0 {
		return fmt.Errorf("input shape %v is not a fixed [1, 3, H, W]", aInputShape)
	}
	oSelf.height = int(aInputShape[2])
	oSelf.width = int(aInputShape[3])
	oSelf.inputSize = 3 * oSelf.height * oSelf.width

	cDevice := C.CString(sDevice)
	defer C.free(unsafe.Pointer(cDevice))
	if err := check(C.ov_core_compile_model_props(oSelf.core, oSelf.model, cDevice, 0, nil, &oSelf.compiled), "compile model for "+sDevice); err != nil {
		return err
	}
	if err := check(C.ov_compiled_model_create_infer_request(oSelf.compiled, &oSelf.request), "create infer request"); err != nil {
		return err
	}

	var cShape C.ov_shape_t
	if err := check(C.ov_shape_create(C.int64_t(len(aInputShape)), (*C.int64_t)(unsafe.Pointer(&aInputShape[0])), &cShape), "create input shape"); err != nil {
		return err
	}
	defer C.ov_shape_free(&cShape)
	if err := check(C.ov_tensor_create(C.F32, cShape, &oSelf.input), "create input tensor"); err != nil {
		return err
	}
	if err := check(C.ov_tensor_data(oSelf.input, &oSelf.inputData), "get input tensor data"); err != nil {
		return err
	}
	return check(C.ov_infer_request_set_input_tensor(oSelf.request, oSelf.input), "set input tensor")
}

func (oSelf *Model) InputSize() (int, int) {
	return oSelf.height, oSelf.width
}

func (oSelf *Model) Run(aInput []float32) ([]float32, []int64, error) {
	if len(aInput) != oSelf.inputSize {
		return nil, nil, fmt.Errorf("input has %d values, want %d", len(aInput), oSelf.inputSize)
	}

	oSelf.mutex.Lock()
	defer oSelf.mutex.Unlock()

	copy(unsafe.Slice((*float32)(oSelf.inputData), oSelf.inputSize), aInput)
	if err := check(C.ov_infer_request_infer(oSelf.request), "infer"); err != nil {
		return nil, nil, err
	}

	var cOutput *C.ov_tensor_t
	if err := check(C.ov_infer_request_get_output_tensor(oSelf.request, &cOutput), "get output tensor"); err != nil {
		return nil, nil, err
	}
	defer C.ov_tensor_free(cOutput)

	var cType C.ov_element_type_e
	if err := check(C.ov_tensor_get_element_type(cOutput, &cType), "get output type"); err != nil {
		return nil, nil, err
	}
	if cType != C.F32 {
		return nil, nil, fmt.Errorf("model output is not float32")
	}

	var cShape C.ov_shape_t
	if err := check(C.ov_tensor_get_shape(cOutput, &cShape), "get output shape"); err != nil {
		return nil, nil, err
	}
	aShape := slices.Clone(unsafe.Slice((*int64)(unsafe.Pointer(cShape.dims)), int(cShape.rank)))
	C.ov_shape_free(&cShape)

	var cSize C.size_t
	if err := check(C.ov_tensor_get_size(cOutput, &cSize), "get output size"); err != nil {
		return nil, nil, err
	}
	var pData unsafe.Pointer
	if err := check(C.ov_tensor_data(cOutput, &pData), "get output data"); err != nil {
		return nil, nil, err
	}

	// tensor 在 return 時就釋放了，所以要複製一份。
	return slices.Clone(unsafe.Slice((*float32)(pData), int(cSize))), aShape, nil
}

func (oSelf *Model) Close() error {
	if oSelf.input != nil {
		C.ov_tensor_free(oSelf.input)
		oSelf.input = nil
	}
	if oSelf.request != nil {
		C.ov_infer_request_free(oSelf.request)
		oSelf.request = nil
	}
	if oSelf.compiled != nil {
		C.ov_compiled_model_free(oSelf.compiled)
		oSelf.compiled = nil
	}
	if oSelf.model != nil {
		C.ov_model_free(oSelf.model)
		oSelf.model = nil
	}
	if oSelf.core != nil {
		C.ov_core_free(oSelf.core)
		oSelf.core = nil
	}
	return nil
}

// inputShape 讀模型第一個輸入的形狀。
func inputShape(cModel *C.ov_model_t) ([]int64, error) {
	var cPort *C.ov_output_const_port_t
	if err := check(C.ov_model_const_input(cModel, &cPort), "get input port"); err != nil {
		return nil, err
	}
	defer C.ov_output_const_port_free(cPort)

	var cShape C.ov_shape_t
	if err := check(C.ov_const_port_get_shape(cPort, &cShape), "get input shape"); err != nil {
		return nil, err
	}
	defer C.ov_shape_free(&cShape)

	return slices.Clone(unsafe.Slice((*int64)(unsafe.Pointer(cShape.dims)), int(cShape.rank))), nil
}

func check(cStatus C.ov_status_e, sAction string) error {
	if cStatus == 0 {
		return nil
	}
	return fmt.Errorf("%s: %s (status %d)", sAction, C.GoString(C.ov_get_last_err_msg()), int(cStatus))
}
