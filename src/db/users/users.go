/*
структура для работы с моделью таблицы users
*/

package dbUsers

import (
	"et/src/di"
)

/*
Структура для работы с данными пользователя
*/
type UsersDB struct {
	Id       int64
	Email    string
	Password string
	Role     int
}

///*
//Проверка пароля Заменить хеш пароль на пароль
//*/
//func (user *UsersDB) CheckPasswordHash() bool {
//	err := bcrypt.CompareHashAndPassword([]byte(user.hashPassword), []byte(user.Password))
//	return err == nil
//}

/*
Проверяет есть ли email в базе
*/
func (user *UsersDB) IsEmailAlreadyThere() (isNotFound bool, err error) {
	sql := "SELECT count(*) FROM users WHERE email = ?"

	rows, err := di.DI.DB.Query(sql, user.Email)

	if err != nil {
		return
	}
	defer rows.Close()

	rows.Next()
	var rowCount int
	err = rows.Scan(&rowCount)

	if err != nil {
		return
	}

	if rowCount == 0 {
		isNotFound = true
	}

	return
}

/*
Добавляет нового пользователя в БД
*/
func (user *UsersDB) AddNewUser() (err error) {
	sql := "INSERT INTO users (`email`, `password`, `role`) VALUES (?, ?, ?)"

	res, err := di.DI.DB.Exec(sql, user.Email, user.Password, user.Role)
	if err != nil {
		return
	}

	user.Id, err = res.LastInsertId()
	if err != nil {
		return
	}

	return
}

/*
Отмечат E_Mail как подтвержденный
*/
func (user *UsersDB) UpdateConfirmEmail() (err error) {
	sql := "UPDATE users SET email_confirm = 1 WHERE id =?"

	_, err = di.DI.DB.Exec(sql, user.Id)

	return
}
