package dbUserEmailConfirm

/*
Структура осуществляет работу по подтверждению Email

Одновременно модель работает и с моделью данных. Разносить модели хранения и обработки данных
смысла не имеет. Т.е. Эти модели применяются единожды и только для подтверждения Email
*/

import (
	dbUsers "et/src/db/users"
	"et/src/di"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type ConfirmEmail struct {
	ConfirmCodeSent string // код подтверждения полученный для проверки от пользователя

	UserId         int64  // id пользователя Email которого подтверждается
	ConfirmCode    string // код для подтверждения сохраненный в БД
	IsEmailConfirm int    // состояние подтверждения - 1 - подтверждено, 0 - нет
	TimestampAdd   int64  // время когда был создан запрос на подтверждение
	Attempts       int    // число попыток подтверждения по текущему ConfirmCode

	UserDB dbUsers.UsersDB // ссылка на модель данных users
}

/*
Добавляет запись с кодом подтверждения
*/
func (cnfMl *ConfirmEmail) AddConfirmCodeInDB() (err error) {
	// создать код подтверждения
	if err = cnfMl.createConfirmCode(); err != nil {
		return
	}

	sql := "INSERT INTO " +
		"users_email_confirm (`user_id`, `confirm_code`, `timestamp_add`) VALUES (?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE `confirm_code` = ?, `timestamp_add` = ?, attempts = 0"

	_, err = di.DI.DB.Exec(sql,
		cnfMl.UserId, cnfMl.ConfirmCode, cnfMl.TimestampAdd,
		cnfMl.ConfirmCode, cnfMl.TimestampAdd)

	return
}

/*
проверяет валидность кода подтверждения E-Mail
*/
func (cnfMl *ConfirmEmail) IsNotValidConfirmCode() (err error) {
	// получить данные о состоянии подтверждения e-mail
	err = cnfMl.findConfirmCode()
	if err != nil {
		return
	}

	// проверить данные на валидность
	err = cnfMl.isValidConfirmEmail()
	if err != nil {
		return
	}

	// данные валидны // отметить E-Mail как подтвержденный
	cnfMl.UserDB.Id = cnfMl.UserId

	if cnfMl.UserDB.UpdateConfirmEmail() != nil {
		err = fmt.Errorf("Ошибка БД. Повторите попытку чуть позже.")
		return
	}

	return
}

/*
Получает запись о коде подтверждения пользователя
*/
func (cnfMl *ConfirmEmail) findConfirmCode() (err error) {
	sql := "SELECT " +
		"confirm_code, timestamp_add, attempts, email_confirm " +
		"FROM " +
		"users_email_confirm " +
		"JOIN users ON (users_email_confirm.user_id = users.id)" +
		"WHERE user_id = ?"

	rows, err := di.DI.DB.Query(sql, cnfMl.UserId)

	if err != nil {
		err = fmt.Errorf("Ошибка БД. Повторите попытку чуть позже.")
		return
	}
	defer rows.Close()

	// передать в структуру полученые из БД данные
	countRow := 0
	for rows.Next() {
		err = rows.Scan(
			&cnfMl.ConfirmCode, &cnfMl.TimestampAdd,
			&cnfMl.Attempts, &cnfMl.IsEmailConfirm)

		if err != nil {
			err = fmt.Errorf("Ошибка БД. Повторите попытку чуть позже.")
			return
		}

		countRow++
	}

	if countRow == 0 {
		err = fmt.Errorf("Вы перешли по неправильной ссылке. Попробуйте скопировать ссылку и вставить ее в адресную строку. Или запросите ссылку на подтверждение повторно.")
		return
	}

	return
}

func (cnfMl *ConfirmEmail) isValidConfirmEmail() (err error) {
	if cnfMl.IsEmailConfirm == 1 {
		err = fmt.Errorf("Ваш E-Mail успешно подтерждён.")
		return // email уже был ранее успешно подтвержден
	}

	if cnfMl.TimestampAdd < (time.Now().Unix() - di.DI.ConfirmWaitingTimeSec) {
		// код поддтверждения просрочен
		err = fmt.Errorf("Срок действия Вашей ссылки истек. Запросите ссылку на подтверждение повторно.")

		// удалить просроченную запись из базы
		// удаление будет осуществляться скриптом запускаемой по крону

		return
	}

	if cnfMl.Attempts > di.DI.ConfirmAttempts {
		// превышено допустимое число ошибок подтверждения по текущей ссылке
		err = fmt.Errorf("Вы перешли по недействительной ссылке. Запросите повторно ссылку на подтверждение Email.")
		return
	}

	if cnfMl.ConfirmCode != cnfMl.ConfirmCodeSent {
		// код из ссылки не равен коду из базы
		err = fmt.Errorf("Вы перешли по неправильной ссылке. Попробуйте скопировать ссылку и вставить ее в адресную строку. Или запросите ссылку на подтверждение повторно.")

		// обновить информацию о количестве попыток подтверждения
		cnfMl.updateConfirmEmailAttempts()

		return
	}

	return
}

/*
Увеличивает счетчик ошибок подтверждения e-mail
*/
func (cnfMl *ConfirmEmail) updateConfirmEmailAttempts() {
	sql := "UPDATE users_email_confirm SET attempts = attempts + 1 WHERE user_id =?"

	di.DI.DB.Exec(sql, cnfMl.UserId)
}

/*
Создает код подтверждения
*/
func (cnfMl *ConfirmEmail) createConfirmCode() (err error) {
	cnfMl.TimestampAdd = time.Now().Unix()

	confirmCode, err := bcrypt.GenerateFromPassword(([]byte(strconv.FormatInt(cnfMl.TimestampAdd, 10))), 10)
	if err != nil {
		return
	}

	cnfMl.ConfirmCode = string(confirmCode)

	return
}
