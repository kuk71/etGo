package mdlUsers

import (
	userEC "et/src/db/userEmailConfirm"
	dbUsers "et/src/db/users"
	"et/src/di"
	"et/src/helpers"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
)

type UserModel struct {
	Email           string
	Password        string
	PasswordReplace string

	Error string

	EmailErr           string
	PasswordErr        string
	PasswordReplaceErr string

	UserDB dbUsers.UsersDB // модель данных user
}

/*
Регистрирует нового пользователя
*/
func (user *UserModel) Registration() bool {
	// проверить переданные данные
	if !user.isValidFormSent() {
		return false
	}

	// зарегистрировать нового пользователя в БД
	if err := user.addUserInDB(); err != nil {
		return false
	}

	// Запросить подтверждение Email
	err := user.createMailConfirm()
	if err != nil {
		log.Println(err) // ошибка не критичная // поэтому false не возвращается
	}

	return true
}

/*
Тестирует на валидность данные переданные из формы
*/
func (user *UserModel) isValidFormSent() bool {
	return user.isMailValid() && user.isPassValid() && user.isPassRplValid()
}

/*
Проверяет пароль на соответствие правилам
*/
func (user *UserModel) isPassValid() bool {
	errCode := 0
	user.Password, errCode = helpers.TestString(user.Password, 8, 0, false)

	if errCode == 1 {
		user.PasswordErr = "Поле не должно быть короче 8-х символов"
		return false
	}

	return true
}

/*
Проверяет Пароль и Повторный пароль на совпадение
*/
func (user *UserModel) isPassRplValid() bool {
	if user.Password != user.PasswordReplace {
		user.PasswordReplaceErr = "В полях Пароль и Повторите пароль должны быть одинаковые значения"
		return false
	}

	return true
}

/*
Проверяет на валидность новый адрес электронной почты
*/
func (user *UserModel) isMailValid() bool {
	errCode := 0

	user.Email, errCode = helpers.TestMail(user.Email)

	switch errCode {
	case 1:
		user.EmailErr = "Поле должно быть заполнено"
	case 2, 3:
		user.EmailErr = "В поле должен быть введен адрес электронной почты"
	}

	if user.EmailErr != "" {
		return false
	}

	// проверить email на существование в БД
	if user.isMailNotNew() {
		return false
	}

	return true
}

/*
Проверяет наличие Email в БД
*/
func (user *UserModel) isMailNotNew() bool {
	user.UserDB.Email = user.Email

	isNotFound, err := user.UserDB.IsEmailAlreadyThere()

	if err != nil {
		user.EmailErr = "Ошибка БД. Повторите попытку чуть позже."
		return true
	}

	if !isNotFound {
		user.EmailErr = "Такой E-mail уже зарегистрирован в системе."
		return true
	}

	return false
}

func (user *UserModel) addUserInDB() (err error) {
	// захешировать пароль
	if err = user.HashPassword(); err != nil {
		user.Error = "Не предвиденная ошибка. Повторите попытку чуть позже. Если ошибка не исчезнет, пожалуйста, сообщите о ней администратору сайта."
		return
	}

	user.UserDB.Email = user.Email
	user.UserDB.Role = 1 // роль пользователя = пользователь

	if err = user.UserDB.AddNewUser(); err != nil {
		user.Error = "Ошибка БД. Повторите попытку чуть позже."
	}

	return
}

/*
Шифрование пароля
*/
func (user *UserModel) HashPassword() (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	user.UserDB.Password = string(bytes)

	return
}

/*
* отправляет письмо для подтверждения Email
 */
func (user *UserModel) createMailConfirm() (err error) {
	//создать запись БД с кодом подтверждения Email
	confirmCode, err := user.addConfirmCode()
	if err != nil {
		return
	}

	// отправить письмо
	err = user.sendMailConfirm(confirmCode)

	return
}

func (user *UserModel) addConfirmCode() (confirmCode string, err error) {
	var confirmEmail userEC.ConfirmEmail

	confirmEmail.UserId = user.UserDB.Id

	err = confirmEmail.AddConfirmCodeInDB()
	if err == nil {
		confirmCode = confirmEmail.ConfirmCode
	}

	return
}

func (user *UserModel) sendMailConfirm(confirmCode string) (err error) {
	body := "Вы зарегистрировались на сайте" + di.DI.ServerProtocol + "://" + di.DI.ServerDomen + ".\n" +
		"Для подтверждения адреса электронной почты, пожалуйста перейдите по ссылке:\n" +
		di.DI.ServerProtocol + "://" + di.DI.ServerDomen + "/email_confirm.html?id=" +
		strconv.FormatInt(user.UserDB.Id, 10) + "&confirm=" + confirmCode

	subject := "Запрос на подтверждение E-Mail"

	err = helpers.SendMail(user.Email, subject, body,
		di.DI.ConfirmEmail, di.DI.ServerDomen, di.DI.ConfirmEmail,
		di.DI.ConfirmEmailPassword, di.DI.ConfirmEmailHost)

	return
}
