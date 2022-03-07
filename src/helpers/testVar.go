package helpers

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

/*
проверяет переменную на минимальную и максимальную длинну, а также обрезает концевые пробелв
str - проверяемая строка
minLen - минимальная длинна // если значение 0 - параметр не контролируется
maxLen - максимальная длинна  // если значение 0 - параметр не контролируется
trim - true - пробелы нужно обрезать

Возвращаемые ошибки
0 - отсутствие ошибок
1 - строка короче минимальной длинны
2- строка длиннее максимально допустимой длинны
*/
func TestString(str string, minLen, maxLen int, trim bool) (string, int) {
	errCode := 0

	if trim {
		str = strings.TrimSpace(str)
	}

	if (minLen != 0) && (utf8.RuneCountInString(str) < minLen) {
		errCode = 1
	}

	if (maxLen != 0) && (utf8.RuneCountInString(str) > maxLen) {
		errCode = 2
	}

	return str, errCode
}

/**
Проверяет валиднось e-mail адреса
mail - проверяемый адрес
возвращает mail - зачишенный об пробелов
код ошибки:
0 - ошибок нет
1 - данные не переданы т.е. trim(mail) == ""
2 - не правильня длинна поля
3 - не правильный формат mail
*/
func TestMail(mail string) (string, int) {
	// emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	emailRegex := regexp.MustCompile("\\w+([\\.-]?\\w+)*@\\w+([\\.-]?\\w+)*\\.\\w{2,4}")

	mail = strings.TrimSpace(mail)

	if utf8.RuneCountInString(mail) == 0 {
		return mail, 1
	}

	if utf8.RuneCountInString(mail) < 3 && utf8.RuneCountInString(mail) > 254 {
		return mail, 2
	}

	if !emailRegex.MatchString(mail) {
		return mail, 3
	}

	return mail, 0
}
