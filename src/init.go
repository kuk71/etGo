package main

import (
	"database/sql"
	"et/src/di"
)

func start() {
	dbConnect()
}

func dbConnect() {
	var err error

	// основные настройки к базе
	dsn := di.DI.DbUser + ":" + di.DI.DbPassword
	dsn += "@tcp(" + di.DI.DbHost + ":" + di.DI.DbPort + ")"
	dsn += "/" + di.DI.DbName + "?"

	dsn += "&charset=" + di.DI.DbCharset
	// отказываемся от prapared statements
	// параметры подставляются сразу
	dsn += "&interpolateParams=true"

	di.DI.DB, err = sql.Open(di.DI.DbDriver, dsn)

	if err != nil {
		panic(err)
	}

	di.DI.DB.SetMaxOpenConns(10)

	err = di.DI.DB.Ping() // пинг базы на предмет ошибки подключения
	if err != nil {
		panic(err)
	}
}
