package statistics

import (
	"fmt"
	"net/http"
)

const template = "Количество обращений за поледнюю минуту %d"

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	resultMessage := fmt.Sprintf(template, s.counter.Count())

	if _, err := w.Write([]byte(resultMessage)); err != nil {
		// TODO: add logger
	}
}
