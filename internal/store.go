package internal

import (
	"errors"
	"slices"
	"sync"
)

type Server struct {
	Addr     string
	Port     int
	HealthOK bool
}

func (s *Server) DeepCopy() *Server {
	return &Server{
		Addr:     s.Addr,
		Port:     s.Port,
		HealthOK: s.HealthOK,
	}
}

type ServerStore struct {
	servers []*Server
	mu      sync.RWMutex
}

func NewServerStore() *ServerStore {
	return &ServerStore{
		servers: []*Server{},
		mu:      sync.RWMutex{},
	}
}

func (s *ServerStore) List() []*Server {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*Server, 0, len(s.servers))
	for _, sv := range s.servers {
		res = append(res, sv.DeepCopy())
	}
	return res
}

func (s *ServerStore) Set(sv *Server) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := slices.IndexFunc(s.servers, func(v *Server) bool {
		return v.Addr == sv.Addr
	})
	if i == -1 {
		s.servers = append(s.servers, sv.DeepCopy())
		return
	}

	s.servers[i] = sv.DeepCopy()
}

func (s *ServerStore) Get(addr string) (*Server, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	i := slices.IndexFunc(s.servers, func(v *Server) bool {
		return v.Addr == addr
	})
	if i == -1 {
		return nil, errors.New("not found")
	}

	return s.servers[i].DeepCopy(), nil
}
