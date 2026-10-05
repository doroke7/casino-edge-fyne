package openvino

// macOS：brew install openvino，標頭檔與函式庫都在 /opt/homebrew。
//
// 檔名有 _darwin 後綴，所以只有在 macOS 編譯時才會用到這個檔案，其他系統會直接忽略它。
// 在 macOS 編譯時，下面的 #cgo 指令告訴編譯器與連結器：
//
//	-I/opt/homebrew/include   編譯 C 程式時，去這個目錄找標頭檔（openvino/c/openvino.h）。
//	                          brew 裝的東西不在系統預設位置，所以要指定。
//	-L/opt/homebrew/lib       連結時，去這個目錄找函式庫檔案。
//	-lopenvino_c              連結時，把 libopenvino_c 連進來，OpenVINO 的函式實作在裡面。
//	-Wl,-rpath,/opt/homebrew/lib
//	                          -Wl, 是把後面逗號分隔的內容原樣交給連結器，也就是 -rpath /opt/homebrew/lib。
//	                          把「執行時去哪找 libopenvino_c.dylib」這個路徑寫進執行檔，
//	                          這樣執行時不用設環境變數。不寫的話編譯能過，但程式一啟動就會報
//	                          dyld: Library not loaded。
//
// 注意：說明和下面的 C 程式之間要空一行；註解緊貼 import "C" 的話，會被當成 C 程式。

/*
#cgo CFLAGS: -I/opt/homebrew/include
#cgo LDFLAGS: -L/opt/homebrew/lib -lopenvino_c -Wl,-rpath,/opt/homebrew/lib
*/
import "C"
