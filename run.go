package admin

type server interface {
	ListenAndServe() error
}
