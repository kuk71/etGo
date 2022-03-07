/*
Пакет с зависимостями
*/
package di

import (
	"database/sql"
	"github.com/spf13/viper"
	"html/template"
)

type globalVar struct {
	DB *sql.DB // ссылка на подключение к БД

	ServerPort     string // порт который слушает GO сервер
	ServerDomen    string // доменное имя сайта
	ServerProtocol string // протокол на котором работает сайт http / https

	DbDriver   string // драйвер для подключения к БД
	DbHost     string
	DbPort     string
	DbName     string
	DbUser     string
	DbPassword string
	DbCharset  string

	ConfirmEmail          string // Email с которого осущеcтвляется отправка писем на подтверждение пользовательского Email-a
	ConfirmEmailPassword  string // Пароль для подключения к smtp хосту
	ConfirmEmailHost      string //smtp host и порт для отправки почты
	ConfirmAttempts       int    // число допустимых ошибок подтверждения
	ConfirmWaitingTimeSec int64  // через сколько секунд код подтверждения будет удален как протухший

	TemplateMain *template.Template // откомпилированный шаблон верхнего уровня для вставки в него содержимого HTML страниц
}

var DI globalVar

func Config(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	// настройки сервера и сайта
	DI.ServerPort = viper.GetString("server.port")
	DI.ServerDomen = viper.GetString("server.domen")
	DI.ServerProtocol = viper.GetString("server.protocol")

	DI.DbDriver = viper.GetString("db.driver")
	DI.DbHost = viper.GetString("db.host")
	DI.DbPort = viper.GetString("db.port")
	DI.DbName = viper.GetString("db.dbName")
	DI.DbUser = viper.GetString("db.user")
	DI.DbPassword = viper.GetString("db.password")
	DI.DbCharset = viper.GetString("db.charset")

	DI.ConfirmEmail = viper.GetString("site.confirmEmail")
	DI.ConfirmEmailPassword = viper.GetString("site.confirmEmailPassword")
	DI.ConfirmEmailHost = viper.GetString("site.confirmEmailHost")
	DI.ConfirmAttempts = viper.GetInt("site.confirmAttempts")
	DI.ConfirmWaitingTimeSec = viper.GetInt64("site.confirmWaitingTimeSec")

	// компилируется шаблон верхнего уровня
	DI.TemplateMain = template.Must(template.ParseFiles("template/main.tmpl"))
}
