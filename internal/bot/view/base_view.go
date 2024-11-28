package view

import (
	"fmt"
	"text/template"
)

// baseView - базовая реализация view для составления сообщений для напоминаний и заметок
type baseView struct {
	text string
}

type record interface {
	// Fields возвращает значения всех доступных полей
	Fields() map[string]string
}

// textForRecord параметром принимает одну запись (заметки или напоминания).
// Вторым аргументом получает сообщение, которое нужно отформатировать.
// Возвращает готовое отформатированное сообщение.
// Сообщение должно иметь подходящие для парсера ключи, например ".ID"
func textForRecord(record record, message string) (string, error) {
	var res = ""

	tmpl, err := template.New("test").Parse(message)
	if err != nil {
		return "", fmt.Errorf("base view: error parsing template: %+v", err)
	}

	view := &baseView{}
	err = tmpl.Execute(view, record.Fields())
	if err != nil {
		return "", fmt.Errorf("base view: error executing template: %+v", err)
	}

	fields := record.Fields()

	res += fmt.Sprintf(message, fields)

	return view.text, nil
}

func (s *baseView) Write(p []byte) (n int, err error) {
	s.text += string(p)
	return len(p), nil
}
