package internal

import (
	"context"
	"log"
)

func Start(store *ServerStore) (shutdown func(ctx context.Context) error, err error) {
	healthShutdown, err := startHealthChecker(store)
	if err != nil {
		return nil, err
	}
	adminShutdown, err := startAdminServer(store)
	if err != nil {
		return nil, err
	}
	httpShutdown, err := startHTTPServer(store)
	if err != nil {
		return nil, err
	}
	dnsShutdown, err := startDNSServer(store)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		log.Println("admin shutdown...")
		if err := adminShutdown(ctx); err != nil {
			return err
		}
		log.Println("http shutdown...")
		if err := httpShutdown(ctx); err != nil {
			return err
		}
		log.Println("dns shutdown...")
		if err := dnsShutdown(ctx); err != nil {
			return err
		}
		log.Println("health checker shutdown...")
		if err := healthShutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil

}
