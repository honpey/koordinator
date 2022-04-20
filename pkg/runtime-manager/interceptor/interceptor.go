package interceptor

type Interceptor interface {
	Name() string
	Setup() error
}
