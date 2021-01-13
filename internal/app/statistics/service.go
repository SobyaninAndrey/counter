package statistics

type Counter interface {
	Count() int
}

type Service struct {
	counter Counter
}

func New(counter Counter) *Service {
	return &Service{counter: counter}
}
