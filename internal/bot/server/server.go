package server

type Server struct {
	NoteEditor     noteEditor
	ReminderEditor reminderEditor
	UserEditor     userEditor
}

type noteEditor interface{}
type reminderEditor interface{}
type userEditor interface{}

func New(noteEditor noteEditor, reminderEditor reminderEditor, userEditor userEditor) *Server {
	return &Server{noteEditor, reminderEditor, userEditor}
}
