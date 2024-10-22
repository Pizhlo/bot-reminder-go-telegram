package elastic

const (
	NoteIndex     = "notes"
	ReminderIndex = "reminders"
)

// Структура для сохранения и поиска в ElasticSearch
type Data struct {
	Model interface{} // Note / Reminder
	Index string
}
