package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

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
}

func NewPipe() *Pipe {
	return &Pipe{
		maxSize: maxSize,
		m:       melody.New(),
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

func (p *Pipe) Scan(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()

		if len(b) > maxRowBytes {
			// skip too long lines (saving memory space)
			continue
		}

		res, err := parser.Parse(b)
		if err != nil {
			fmt.Println(err)
			continue
		}
		encodedRes, err := json.Marshal(res)
		if err != nil {
			fmt.Println("can't marshal:", err)
			continue
		}
		pipe.SaveRecord(encodedRes)
		p.m.Broadcast(encodedRes)
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
