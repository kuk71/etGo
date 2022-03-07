/*
Отрабатывает сценарий подтверждения E-Mail
*/
package emailComfirm

import (
	"bytes"
	dbUuserEmailConfirm "et/src/db/userEmailConfirm"
	"et/src/di"
	"et/src/www"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type paramsTmpl struct {
	IsErr bool
	Err   string

	id          string // id пользователя чей Email подтверждается
	confirmCode string // присланный код подтверждения
}

func EmailСonfirm(w http.ResponseWriter, r *http.Request) {
	var params paramsTmpl
	var confirmEmail dbUuserEmailConfirm.ConfirmEmail

	// проверка переданны ли данные
	if confirmParse(r, &params, &confirmEmail) != nil {
		emailRender(w, &params)
		return
	}

	// проверка кода подтверждения
	if err := confirmEmail.IsNotValidConfirmCode(); err != nil {
		params.IsErr = true
		params.Err = err.Error()
	}

	emailRender(w, &params)
}

/*
Проверяет переданы ли данные для проверки
*/
func confirmParse(r *http.Request, params *paramsTmpl, confirmEmail *dbUuserEmailConfirm.ConfirmEmail) (err error) {
	err = r.ParseForm()

	if err != nil {
		return
	}

	// зачитывает id и приводит его к int64
	confirmEmail.UserId, err = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	confirmEmail.ConfirmCodeSent = r.URL.Query().Get("confirm")

	if err != nil || confirmEmail.ConfirmCodeSent == "" {
		err = fmt.Errorf("")

		params.IsErr = true
		params.Err = "Вы перешли по неправильной ссылке. Попробуйте скопировать ссылку и вставить ее в адресную строку. Или запросите ссылку на подтверждение повторно."
	}

	return
}

/*
Отрисовывает страницу подтверждения E-Mail
*/
func emailRender(w http.ResponseWriter, params *paramsTmpl) {
	//переменная в которую будет положен сформированный шаблон содержания страницы
	//для последующей вставки в главный шаблон
	var content bytes.Buffer

	tmpl := template.Must(template.ParseFiles("template/emailConfirm.tmpl"))
	tmpl.Execute(&content, *params)

	str := www.MainParams{template.HTML(content.String())}

	// отрисовка главного шаблона
	di.DI.TemplateMain.Execute(w, str)
}
