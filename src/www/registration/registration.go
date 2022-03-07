package registration

import (
	"bytes"
	"et/src/di"
	mdlUsers "et/src/models/users"
	"et/src/www"
	"html/template"
	"net/http"
)

type paramsTmpl struct {
	IsTotalErr bool // наличие общей ошибки выполнения

	IsPassErr    bool // наличие ошибки в поле Пароль
	IsPassRplErr bool // наличие ошибки в поле Повторите пароль
	IsEmailErr   bool // наличие ошибки в поле Email

	User mdlUsers.UserModel // структура с методами управления Моделью user
}

func Registration(w http.ResponseWriter, r *http.Request) {
	// параметры для передачи в шаблон рисования страницы
	var params paramsTmpl

	if err := r.ParseForm(); err != nil {
		return
	}

	isUserRegistr := false

	// проверить переданные данные
	if r.PostForm.Get("send") != "" {
		// данные переданы - провести продседуру регистрации

		params.User.Email = r.PostForm.Get("mail")
		params.User.Password = r.PostForm.Get("pass")
		params.User.PasswordReplace = r.PostForm.Get("passRpl")

		isUserRegistr = params.User.Registration()
	}

	if !isUserRegistr {
		// аутентификация пользователя

		// переадресация на главную старницу аккаунта
	}

	params.setParamsErr()
	registrationRender(w, &params)
}

// сформировать ответ
func registrationRender(w http.ResponseWriter, params *paramsTmpl) {
	/* переменная в которую будет положен сформированный шаблон содержания страницы
	для последующей вставки в главный шаблон
	*/
	var content bytes.Buffer

	tmpl := template.Must(template.ParseFiles("template/registration.tmpl"))
	tmpl.Execute(&content, *params)

	str := www.MainParams{template.HTML(content.String())}

	di.DI.TemplateMain.Execute(w, str)
}

/*
Определяет какие ошибки выводить в шаблоне
*/
func (params *paramsTmpl) setParamsErr() {
	if params.User.Error != "" {
		params.IsTotalErr = true
	}

	if params.User.EmailErr != "" {
		params.IsEmailErr = true
	}

	if params.User.PasswordErr != "" {
		params.IsPassErr = true
	}

	if params.User.PasswordReplaceErr != "" {
		params.IsPassRplErr = true
	}

}
