//go:generate mockery --case=underscore --name Storage

package counter

import (
	"container/list"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Storage interface {
	Save(events []time.Time) error
	Get() (*list.List, error)
}

type Counter struct {
	mx       sync.RWMutex
	storage  Storage
	events   *list.List
	lifetime time.Duration
}

func New(persistentStorage Storage, lifetime time.Duration) *Counter {
	return &Counter{
		storage:  persistentStorage,
		events:   list.New(),
		lifetime: lifetime,
	}
}

func (c *Counter) Load() error {
	events, err := c.storage.Get()
	if err != nil {
		return errors.Wrap(err, "can not load events")
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	c.events = events

	return nil
}

func (c *Counter) AddEvent() {
	c.mx.Lock()
	defer c.mx.Unlock()

	currentTime := time.Now()
	c.clearBefore(currentTime.Add(-c.lifetime))
	c.events.PushBack(currentTime)
}

func (c *Counter) Cancel() error {
	c.mx.Lock()

	c.clearBefore(time.Now().Add(-c.lifetime))

	data := make([]time.Time, 0, c.events.Len())
	for e := c.events.Front(); e != nil; e = e.Next() {
		if dt, ok := e.Value.(time.Time); ok {
			data = append(data, dt)
		}
	}

	c.mx.Unlock()

	if err := c.storage.Save(data); err != nil {
		return errors.Wrap(err, "can not save")
	}

	return nil
}

func (c *Counter) Count() int {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.clearBefore(time.Now().Add(-c.lifetime))
	return c.events.Len()
}

func (c *Counter) clearBefore(beforeDt time.Time) {
	var next *list.Element
	for e := c.events.Front(); e != nil; e = next {
		if dt, ok := e.Value.(time.Time); ok && dt.After(beforeDt) {
			break
		}
		next = e.Next()
		c.events.Remove(e)
	}
}
