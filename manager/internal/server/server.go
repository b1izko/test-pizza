package server

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/b1izko/test-pizza/internal/logger"
	"github.com/b1izko/test-pizza/manager/internal/handler"
	"github.com/b1izko/test-pizza/manager/store"
	"github.com/b1izko/test-pizza/manager/store/admin"
	"github.com/b1izko/test-pizza/manager/store/order"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/pkg/errors"
)

// Config for server
type Config struct {
	ListenAddress string
	DBUrl         string
	DBName        string
	Username      string
	Password      string
}

// Server wrapper
type Server struct {
	*Config
	run int32

	srv   http.Server
	store *store.Store
	ch    *amqp.Channel
	close chan struct{}
	wg    sync.WaitGroup
}

// New server
func New(cfg *Config, channel *amqp.Channel) *Server {
	return &Server{
		Config: cfg,
		ch:     channel,
	}
}

// Start server
func (s *Server) Start() error {
	if atomic.AddInt32(&s.run, 1) != 1 {
		return errors.New("already started")
	}

	store, err := store.New(s.DBUrl, s.Username, s.Password, s.DBName)
	if err != nil {
		return err
	}

	s.store = store

	s.srv.Addr = s.ListenAddress
	s.srv.Handler = handler.New(store).Get()
	s.close = make(chan struct{})

	q, err := s.ch.QueueDeclare(
		"go",  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if logger.IsError(err, "Failed to declare a queue") {
		return err
	}

	msgs, err := s.ch.Consume(
		q.Name, // queue
		"go",   // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if logger.IsError(err, "Failed to register a consumer") {
		return err
	}

	var forever chan struct{}

	go func() {
		err := s.store.Connect()
		if logger.IsError(err, "Start error") {
			panic(err)
		}

		for d := range msgs {
			//debug
			log.Printf("Received a message: %s", d.Body)

			var order order.Model
			err = order.UnmarshalJSON(d.Body)
			if logger.IsError(err, "Failed to parse message") {
				panic(err)
			}

			err := order.Save(s.store)
			if logger.IsError(err, "Failed to save order") {
				panic(err)
			}
		}
	}()

	<-forever

	s.wg.Add(2)
	go s.loop()
	go s.serve()
	s.init()

	return nil
}

func (s *Server) init() {
	err := s.store.Connect()
	if logger.IsError(err, "Init error ") {
		panic(err)
	}

	root := &admin.Model{
		Login:    "root",
		Password: "67b231db6f904b2fa827881ae9f6913c",
		Status:   admin.StatusRoot,
		LastAuth: time.Now(),
	}

	err = root.Save(s.store)
	if logger.IsError(err, "Init error") {
		panic(err)
	}
}

func (s *Server) serve() {
	defer s.wg.Done()

	err := s.srv.ListenAndServe()
	if logger.IsError(err, "Listen error") {
		panic(err)
	}

	select {
	case <-s.close:
	default:
		close(s.close)
	}
}

func (s *Server) loop() {
	defer s.wg.Done()

main:
	for {
		select {
		case <-s.close:
			break main
		}
	}

	s.srv.Close()
}

// Close server
func (s *Server) Close() error {
	if atomic.LoadInt32(&s.run) == 0 {
		err := errors.New("not started")
		logger.IsError(err, "Listen error")
		return err
	}

	select {
	case <-s.close:
	default:
		close(s.close)
	}

	s.wg.Wait()
	return nil
}
