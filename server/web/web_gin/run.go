package web_gin

type server interface {
	ListenAndServe() error
}
