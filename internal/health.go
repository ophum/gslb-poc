package internal

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"
)

// check tcp port
func startHealthChecker(store *ServerStore) (shutdown func(ctx context.Context) error, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("start health checker")
	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			for _, s := range store.List() {
				newSv := s.DeepCopy()
				isDirty := false
				c, err := net.Dial("tcp", net.JoinHostPort(s.Addr, strconv.Itoa(s.Port)))
				if err != nil {
					if s.HealthOK {
						newSv.HealthOK = false
						isDirty = true
					}
				} else {
					c.Close()
					if !s.HealthOK {
						newSv.HealthOK = true
						isDirty = true
					}
				}

				if isDirty {
					log.Printf("update server health: %s %d %t -> %t", s.Addr, s.Port, s.HealthOK, newSv.HealthOK)
					store.Set(newSv)
				}
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				return
			}
		}
	}()
	return func(ctx context.Context) error {
		cancel()
		return nil
	}, nil
}
