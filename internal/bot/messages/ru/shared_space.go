package messages

const SharedSpacesNotFoundMessage = "Это раздел совместных пространств. Пока что у тебя еще нет ни одного пространства. Чтобы создать, нажми на кнопку"

// Creation
const SharedSpaceNameMessage = "Напиши название нового пространства"
const SharedSpaceCreationSuccessMessage = "Пространство <b>%s</b> успешно создано!"

// Participants
const SharedSpaceMessage = "<b>%d. %s</b>\n\n%s\nЗаметок: %d\nНапоминаний: %d\n\nСоздано: %+v\n\n"
const AddParticipantMessage = "Пришли username, контакт или ссылку на пользователя, которого хочешь добавить в пространство"
const RemoveParticipantMessage = "Выбери, какого участника хочешь исключить"
const InvitationsMessage = "Пользователь @%s приглашает вас в совместное пространство <b>%s</b>"
const UserNotRegisteredMessage = "⚠️Пользователь с таким username не зарегистрирован в боте. Попросите пользователя написать боту и повторите попытку"
const InvalidUserLinkMessage = "Невалидная ссылка. Пришлите другую ссылку и повторите попытку"
const SuccessfullySentInvitationsMessage = "✅Приглашение было успешно отправлено"
const UserAlreadyInvitedMessage = "Пользователь уже приглашен в совместное пространство"
const UserAcceptedInvitationMessage = "✅Пользователь <b>%s</b> принял приглашение в совместное пространство <b>%s</b>"
const UserRejecteddInvitationMessage = "❌Пользователь <b>%s</b> отклонил приглашение в совместное пространство <b>%s</b>"
const InvitationAcceptedMessage = "✅Приглашение в пространство <b>%s</b> успешно принято"
const InvitationRejectedMessage = "✅Приглашение в пространство <b>%s</b> успешно отклонено"
const UserWasRemovedMessage = "⚠️Вы были исключены из пространства <b>%s</b> пользователем %s"
const UserSuccesfullyRemoved = "✅Пользователь <b>%s</b> успешно исключен из пространства <b>%s</b>"

// Records
const NoNotesInSharedSpaceMessage = "В пространстве <b>%s</b> пока не создано ни одной заметки. Чтобы создать, просто пришли текст"
const NoRemindersInSharedSpaceMessage = "В пространстве <b>%s</b> пока не создано ни одного напоминания"

// Notes
const AskNoteTextMessage = "Напиши текст заметки"
const SuccessfullyAddedNoteMessage = "Заметка успешно добавлена в совместное пространство <b>%s</b>!"
const UserAddedNoteMessage = "Пользователь @%s добавил новую заметку в пространство <b>%s</b>!"
