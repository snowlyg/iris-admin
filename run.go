package admin

import "context"

type serve interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}
