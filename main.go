package main

import "log"

// https://github.com/jackc/pgx/pull/545
// https://www.postgresql.org/docs/current/libpq-connect.html + Ctrl+F = target_session_attrs
const connString = "postgres://postgres:example@127.0.0.1:5436/postgres?pool_max_conns=1000" // ?target_session_attrs=any
const goroutines = 5

func main() {
	s := newStorage(newConn(connString))

	for i := 0; i < goroutines; i++ {
		go func(s *storage, pid int) {
			p := newDaemon(s)
			p.run(pid)
		}(s, i)
	}

	log.Fatal(newHandler(s).init().run())
}
