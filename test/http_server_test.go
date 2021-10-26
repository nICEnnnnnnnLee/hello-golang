package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func FuncHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Func handler...")
}

type StructHandler struct {
	content string
}

func (handler *StructHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, handler.content)
}

func CheckValid(url string, expectedResponse string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	result, _ := io.ReadAll(res.Body)
	if !bytes.Equal(result, []byte(expectedResponse)) {
		tips := fmt.Sprintf("Received: %s, while expect: %s\n", string(result), expectedResponse)
		return errors.New(tips)
	}
	return err
}
func TestSimple(t *testing.T) {
	var err error
	go func() {
		http.HandleFunc("/func", FuncHandler)
		http.Handle("/struct", &StructHandler{content: "Struct handler..."})
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second * 1)
	err = CheckValid("http://localhost:8000/func", "Func handler...")
	if err != nil {
		t.Error(err)
	}
	err = CheckValid("http://localhost:8000/struct", "Struct handler...")
	if err != nil {
		t.Error(err)
	}
}

func TestServerMux(t *testing.T) {
	var err error
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/func", FuncHandler)
		mux.Handle("/struct", &StructHandler{content: "Struct handler..."})
		err := http.ListenAndServe(":8001", mux)
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(time.Second * 1)
	err = CheckValid("http://localhost:8001/func", "Func handler...")
	if err != nil {
		t.Error(err)
	}
	err = CheckValid("http://localhost:8001/struct", "Struct handler...")
	if err != nil {
		t.Error(err)
	}

}

func TestStopServerByOSSignal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/func", FuncHandler)
	mux.Handle("/struct", &StructHandler{content: "Struct handler..."})

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	// 创建系统信号接收器
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server:", err)
		}
	}()

	log.Println("Starting HTTP server...")
	err := server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			log.Print("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected")
		}
	}
}
