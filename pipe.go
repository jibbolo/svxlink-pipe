package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"sync"

	"log"

	"github.com/jibbolo/svxlink-pipe/parser"
	"github.com/olahol/melody"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type record []byte

type Pipe struct {
	history  []record
	histLock sync.Mutex
	maxSize  int
	m        *melody.Melody
	httpPort string
	r        io.Reader
}

func NewPipe(port string, r io.Reader) *Pipe {
	return &Pipe{
		maxSize:  maxSize,
		m:        melody.New(),
		r:        r,
		httpPort: port,
	}
}

func (p *Pipe) SaveRecord(b record) {
	p.histLock.Lock()
	defer p.histLock.Unlock()
	p.history = append(p.history, b)
	if len(p.history) == (p.maxSize + 1) {
		p.history = p.history[1:len(p.history)]
	}
}

func (p *Pipe) WriteHistory(s *melody.Session) {
	p.histLock.Lock()
	defer p.histLock.Unlock()
	for _, record := range p.history {
		s.Write(record)
	}
}

func (p *Pipe) Scan() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("cannot scan from input: %v", err)
		}
	}()

	scanner := bufio.NewScanner(p.r)
	for scanner.Scan() {
		b := scanner.Bytes()

		if len(b) > maxRowBytes {
			// skip too long lines (saving memory space)
			continue
		}

		res, err := parser.Parse(b)
		if err != nil {
			fmt.Println("can't parse input:", err)
			continue
		}
		if len(res) > 0 {
			pipe.SaveRecord(res)
			p.m.Broadcast(res)
		}
	}
}

func (p *Pipe) NewRouter() chi.Router {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(templateIndex)
	})

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		p.m.HandleRequest(w, r)
	})

	p.m.HandleConnect(func(s *melody.Session) {
		pipe.WriteHistory(s)
	})
	return r
}

func (p *Pipe) Run() error {
	go p.Scan()
	return http.ListenAndServe(":"+p.httpPort, p.NewRouter())
}
