package roll

import (
	"dndroller/internal/model"
	"sync"
	"time"
)

type rollMap struct {
	sync.Map
}

func newRollMap() *rollMap {
	return &rollMap{sync.Map{}}
}
func (x *rollMap) Store(id int, c *model.Roll) {
	x.Map.Store(id, c)
}

func (x *rollMap) StoreWithTll(id int, c *model.Roll, ttl time.Duration) {
	x.Store(id, c)
	go func() {
		t := time.NewTimer(ttl)
		<-t.C
		x.Map.Delete(id)
	}()
}

func (x *rollMap) Load(id int) (*model.Roll, bool) {
	v, ok := x.Map.Load(id)
	if !ok {
		return nil, false
	}
	c, ok := v.(*model.Roll)
	return c, ok
}

func (x *rollMap) Delete(id int) {
	x.Map.Delete(id)
}
