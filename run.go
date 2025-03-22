package admin

type serve interface {
	ListenAndServe() error
}
