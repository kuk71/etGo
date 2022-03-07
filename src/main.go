package main

import (
	"database/sql"
	"et/src/di"
	"et/src/www/emailComfirm"
	"et/src/www/registration"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

// переменная для скомпилированного главного шаблона
var tmplMain *template.Template

var db *sql.DB

func main() {
	// загрузка конфигурации
	di.Config("./")

	// Инициализация
	start()
	defer di.DI.DB.Close()

	db = di.DI.DB

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/registration.html", registration.Registration)
	adminMux.HandleFunc("/email_confirm.html", emailComfirm.EmailСonfirm)

	siteHandler := wrapper(adminMux)

	fmt.Println("starting server at :", di.DI.ServerPort)
	http.ListenAndServe(":"+di.DI.ServerPort, siteHandler)
}

func wrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("\n**************\nPanic\n**************\n")
			}

		}()

		next.ServeHTTP(w, r)
	})
}
