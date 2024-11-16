package internal

import (
	"context"
	"log"

	"github.com/miekg/dns"
)

func startDNSServer(store *ServerStore) (shutdown func(ctx context.Context) error, err error) {
	roundrobin := 0

	sv := &dns.Server{
		Addr:       ":8053",
		Net:        "udp",
		TsigSecret: nil,
		ReusePort:  true,
		Handler: dns.HandlerFunc(func(w dns.ResponseWriter, m *dns.Msg) {
			res := new(dns.Msg)
			res.SetReply(m)
			res.Authoritative = true

			for _, q := range m.Question {
				log.Println("name", q.Name, "qclass", q.Qclass, "qtype", q.Qtype)
				switch q.Qtype {
				case dns.TypeA:
					if q.Name != "www.example.com." {
						continue
					}
					ans := []dns.RR{}
					for _, s := range store.List() {
						if !s.HealthOK {
							continue
						}

						rr, err := dns.NewRR(q.Name + " 10 A " + s.Addr)
						if err != nil {
							continue
						}
						ans = append(ans, rr)
					}
					if len(ans) < 1 {
						res.Answer = append(res.Answer, ans...)
						continue
					}

					i := roundrobin % len(ans)
					ans[0], ans[i] = ans[i], ans[0]
					res.Answer = append(res.Answer, ans...)

					roundrobin++

				}
			}
			w.WriteMsg(res)
		}),
	}

	go func() {
		log.Println("start dns server :8053")
		if err := sv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	return sv.ShutdownContext, nil
}
