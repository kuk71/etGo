package www

import "html/template"

// структура для передачи в главный шаблон
type MainParams struct {
	Content template.HTML // этот тип не экранируется при выводе в шаблон
}
