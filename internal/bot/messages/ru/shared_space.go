package messages

const SharedSpacesNotFoundMessage = "Это раздел совместных пространств. Пока что у тебя еще нет ни одного пространства. Чтобы создать, нажми на кнопку"

// Creation
const SharedSpaceNameMessage = "Напиши название нового пространства"
const SharedSpaceCreationSuccessMessage = "Пространство <b>%s</b> успешно создано!"

// Participants
const SharedSpaceMessage = "<b>%d. %s</b>\n\nУчастники:\n%+v\nЗаметок: %d\nНапоминаний: %d\n\nСоздано: %+v\n\n"
const AddParticipantMessage = "Пришли username или ссылку пользователя, которого хочешь добавить в пространство"
const InvintationMessage = "Пользователь @%s приглашает вас в совместное пространство <b>%s</b>"
const UserNotRegisteredMessage = "⚠️Пользователь с таким username не зарегистрирован в боте. Попросите пользователя написать боту и повторите попытку"
const InvalidUserLinkMessage = "Невалидная ссылка. Пришлите другую ссылку и повторите попытку"
const SuccessfullySentInvintationMessage = "✅Приглашение было успешно отправлено"

// Records
const NoNotesInSharedSpaceMessage = "В пространстве <b>%s</b> пока не создано ни одной заметки. Чтобы создать, пришли текст или нажми кнопку"
const NoRemindersInSharedSpaceMessage = "В пространстве <b>%s</b> пока не создано ни одного напоминания"

// Notes
const AskNoteTextMessage = "Напиши текст заметки"
const SuccessfullyAddedNoteMessage = "Заметка успешно добавлена в совместное пространство <b>%s</b>!"
