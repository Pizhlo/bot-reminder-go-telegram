package messages

// SAVE
const SuccessfullyCreatedNoteMessage = "Отличная заметка, я уже сохранил!👍\n\n"

// SEARCH
const SearchNotesByTextMessage = "Введи текст, который нужно найти"
const SearchNotesByDateChooseMessage = "Выбери, как именно искать заметки:\n\n<b>* По одной дате</b> - будут найдены все заметки, созданные только в эту дату\n<b>* По диапазону дат </b> - будут найдены заметки, созданные в промежуток между двумя датами"
const UserDoesntHaveNotesMessage = "У тебя пока нет заметок. Чтобы создать, просто пришли мне текст/фото, и я сохраню это!"
const NoNotesFoundByTextMessage = "У тебя нет заметок с таким текстом"
const NoNotesFoundByDateMessage = "У тебя нет заметок, созданных %s"
const NoNotesFoundByTwoDatesMessage = "У тебя нет заметок, созданных в промежутке между %s и %s"
const SearchOneDateMessage = "Выбери, за какую дату искать заметки"
const FirstDateBeforeSecondMessage = "Вторая дата раньше первой, выбери другую:"
const SecondDateInFutureMessage = "Вторая дата еще не наступила, выбери другую:"
const FirstDateInFutureMessage = "Первая дата еще не наступила, выбери другую:"

const SearchByTwoDatesFirstDateMessage = "Выбери первую дату:"
const SearchByTwoDatesSecondDateMessage = "Выбери вторую дату:"

// DELETE
const ConfirmDeleteNotesMessage = "Ой-ой... Ты точно хочешь удалить ВСЕ заметки?😥"
const AllNotesDeletedMessage = "Все заметки успешно удалены!👌"
const NotDeleteMessage = "Я отменил операцию😌"
const NoteDeletedSuccessMessage = "Заметка под номером <b>%d</b> успешно удалена!🥳"
const NoNoteFoundByNumberMessage = "У тебя нет заметки под номером %d🤔"
