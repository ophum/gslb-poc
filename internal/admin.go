package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func startAdminServer(store *ServerStore) (shutdown func(ctx context.Context) error, err error) {
	e := echo.New()
	e.GET("/servers", func(c echo.Context) error {
		return c.JSON(http.StatusOK, store.List())
	})
	e.POST("/servers", func(c echo.Context) error {
		var req struct {
			Addr string `json:"addr"`
			Port int    `json:"port"`
		}
		if err := c.Bind(&req); err != nil {
			return err
		}

		isNew := false
		sv, err := store.Get(req.Addr)
		if err != nil {
			// only not found
			sv = &Server{
				Addr: req.Addr,
			}
			isNew = true
		} else if sv.Port == req.Port {
			// unchanged
			return c.String(http.StatusOK, fmt.Sprintf("%s:%d unchanged\n", req.Addr, req.Port))
		}

		sv.Port = req.Port
		sv.HealthOK = false
		store.Set(sv)

		msg := fmt.Sprintf("%s:%d ", req.Addr, req.Port)
		if isNew {
			msg += "created\n"
		} else {
			msg += "updated\n"
		}
		return c.String(http.StatusOK, msg)
	})
	httpServer := &http.Server{
		Addr:    ":5000",
		Handler: e,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(err)
		}
	}()

	return httpServer.Shutdown, nil
}
