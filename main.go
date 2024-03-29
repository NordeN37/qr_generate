package main

import (
	"encoding/base64"
	"github.com/skip2/go-qrcode"
	"html/template"
	"log"
	"net/http"
	"os"
	"qr_code_generate/utils/logger"
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
	log := logger.New("debug")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var qrData string
		var resData Data
		switch r.Method {
		case http.MethodPost:
			if qrData = r.FormValue("qr_data"); qrData == "" {
				var dataErr = "input no data"
				bImg, err := os.ReadFile("./html/image_not_found.png")
				if err != nil {
					log.Err(err).Send()
					return
				}
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
			log.Err(err).Send()
			return
		}
		log.Info().Str("qr_data", qrData).Send()
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("server listen error")
	}
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
