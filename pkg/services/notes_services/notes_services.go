package notesservices

type NotesServices struct {
}

type Notes interface {
	CacheTodo(todo interface{}) error
}

func NewNotesServices() *NotesServices {
	return &NotesServices{}
}
