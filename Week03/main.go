package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

var S *http.Server

type Handler struct {
}

func (t Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	// S.Shutdown(context.Background())
	fmt.Fprintf(w, "Hello there!\n")
}

func main() {
	S = &http.Server{
		Addr:           ":8081",
		Handler:        &Handler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("receive panic", err)
			}
		}()
		return S.ListenAndServe()
	})

	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("receive panic", err)
			}
		}()
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		S.Shutdown(ctx)
		return nil
	})

	g.Go(func() error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("receive panic", err)
			}
		}()
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt)
		select {
		case sig := <-signalChan:
			fmt.Println("get sign Server ...", sig)
		case <-ctx.Done():
			fmt.Println("get shut down")
		}
		return errors.New("receive channel")
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("finish")
}
