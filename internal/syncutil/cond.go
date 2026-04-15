package syncutil

import (
	"sync"
	"time"
)

func WaitWithDeadline(c *sync.Cond, d time.Time) {
	timer := time.AfterFunc(time.Until(d), c.Broadcast)
	defer timer.Stop()
	c.Wait()
}
