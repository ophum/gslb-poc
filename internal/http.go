package internal

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func startHTTPServer(store *ServerStore) (shutdown func(ctx context.Context) error, err error) {
	e := echo.New()
	e.GET("/servers", func(c echo.Context) error {
		return c.JSON(http.StatusOK, store.List())
	})
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: e,
	}

	go func() {
		log.Println("start http server :8080")
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			log.Println(err)
		}
	}()

	return func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	}, nil
}
