package errors

import "errors"

// ошибка о том, что данный пользователь уже пригласил этого пользователя в выбранное совместное пространство
var ErrInvitationExists = errors.New("invitation already exists")

// ошибка о том, что пользователь уже состоит в этом пространстве
var ErrUserAlreadyExists = errors.New("user already exists")
