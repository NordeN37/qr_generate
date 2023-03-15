package main

import (
	"encoding/base64"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Создаем структуру для передачи данных в шаблон
type Data struct {
	Png string
	Msg string
}

// Создаем переменную для хранения шаблона
var tmpl *template.Template

// Инициализируем шаблон при запуске программы
func Init() {
	tmpl = template.Must(template.ParseFiles("./html/index.html"))
}

func main() {
	Init()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var qrData string
		var resData Data
		switch r.Method {
		case http.MethodPost:
			if qrData = r.FormValue("qr_data"); qrData == "" {
				var dataErr = "input no data"
				bImg, _ := os.ReadFile("./html/image_not_found.png")
				str := base64.StdEncoding.EncodeToString(bImg)
				resData = Data{Png: str, Msg: dataErr}
			} else {
				str := generateQrBase64(qrData)
				resData = Data{Png: str, Msg: qrData}
			}
		case http.MethodGet:
			qrData = "https://example.org"
			str := generateQrBase64(qrData)
			resData = Data{Png: str, Msg: qrData}
		}
		// Применяем шаблон к данным и отправляем результат в ответ
		if err := tmpl.Execute(w, resData); err != nil {
			log.Println(err)
		}
	})
	http.ListenAndServe("0.0.0.0:80", nil)
}

func generateQrBase64(qrData string) string {
	var data []byte
	data, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		log.Println(err)
	}
	str := base64.StdEncoding.EncodeToString(data)
	return str
}
