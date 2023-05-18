package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pinebit/go-boot/boot"
)

//go:embed testdata/add.wasm
var addWasm []byte

func main() {
	portFlag := flag.Int("port", 8080, "specify server port")
	flag.Parse()

	wasmService := NewWasmService()

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		a := r.URL.Query().Get("a")
		b := r.URL.Query().Get("b")
		ai, err := strconv.Atoi(a)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad 'a' parameter"))
			return
		}
		bi, err := strconv.Atoi(b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad 'b' parameter"))
			return
		}
		res, err := wasmService.Run(r.Context(), addWasm, uint64(ai), uint64(bi))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("result: %d", res)))
		log.Printf("handled request to sum %d and %d, the result is %d", ai, bi, res)
	})

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", *portFlag),
		ReadTimeout:       httpReadTimeoutMs * time.Millisecond,
		ReadHeaderTimeout: httpReadTimeoutMs * time.Millisecond,
		WriteTimeout:      httpWriteTimeoutMs * time.Millisecond,
	}
	httpServerService := boot.NewHttpServer(httpServer)

	appServices := boot.Sequentially(wasmService, httpServerService)

	app := boot.NewApplicationForService(appServices, shutdownTimeoutSeconds*time.Second)
	if err := app.Run(context.Background()); err != nil {
		fmt.Println("server error:", err)
	}
}
