## 這個目錄架構的三大實戰原則：
無 Fyne 污染的 internal/ (最重要)：
internal/ 資料夾裡的程式碼應該純粹只做資料運算或網路通訊。這裡面絕對不要出現 import "fyne.io/fyne/v2"。 這樣能確保你的商業邏輯與 UI 框架完全脫鉤，未來如果想把軟體轉成網頁版或換成 Wails，底層邏輯可以 100% 沿用。

讓 main.go 保持極簡：
不要把所有的視窗佈局都塞在 main.go 裡。main.go 的職責只有一個：呼叫 internal 讀取設定、把資料傳給 ui/screens/ 畫出首頁，然後執行 app.Run()。

靜態資源全部打包進程式碼：
桌面軟體最怕使用者遺失旁邊的圖片檔。Fyne 提供了極度好用的指令 fyne bundle。平常把圖片放在 assets/ 下，編譯前跑一次工具，它會把圖片轉成 bundle/bundled.go 裡的 byte 陣列。最終發布時，你依然只會產出一個乾淨無比的單一執行檔（.exe 或 macOS 的 .app）。