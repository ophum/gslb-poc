package main

import (
	"fmt"
	"net/http"
	"slices"
	"sync"

	"github.com/labstack/echo/v4"
)

func main() {
	servers := map[string]*http.Server{}
	mu := sync.RWMutex{}
	e := echo.New()
	e.GET("/servers", func(c echo.Context) error {
		mu.RLock()
		defer mu.RUnlock()

		addrs := []string{}
		for _, v := range servers {
			addrs = append(addrs, v.Addr)
		}

		slices.Sort(addrs)

		return c.JSON(http.StatusOK, addrs)
	})

	e.POST("/servers", func(c echo.Context) error {
		var req struct {
			Addr string `json:"addr"`
			Port int    `json:"port"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()

		k := fmt.Sprintf("%s:%d", req.Addr, req.Port)
		_, ok := servers[k]
		if ok {
			return c.String(http.StatusOK, "Already started\n")
		}
		s := &http.Server{
			Addr: k,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("listen " + k + "\n"))
			}),
		}
		servers[k] = s

		go s.ListenAndServe()
		return c.String(http.StatusOK, k+" started\n")
	})
	e.DELETE("/servers/:addr", func(c echo.Context) error {
		addr := c.Param("addr")

		mu.Lock()
		defer mu.Unlock()
		s, ok := servers[addr]
		if !ok {
			return c.NoContent(http.StatusNoContent)
		}

		if err := s.Shutdown(c.Request().Context()); err != nil {
			return err
		}

		return c.String(http.StatusOK, addr+" stopped\n")
	})

	if err := e.Start(":8081"); err != nil {
		panic(err)
	}
}
