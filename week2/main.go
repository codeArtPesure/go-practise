package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof"
    "os"
    "os/signal"
    "syscall"
    "time"

    "golang.org/x/sync/errgroup"
)

func main() {
    g, ctx := errgroup.WithContext(context.Background())

    mux := http.NewServeMux()

    mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("pong"))
    })

    // 模拟单个服务错误退出
    serverOut := make(chan struct{})
    mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
        serverOut <- struct{}{}
    })

    server := http.Server{
        Handler: mux,
        Addr:    ":8888",
    }


    g.Go(func() error {
        return server.ListenAndServe()
    })


    g.Go(func() error {
        select {
        case <-ctx.Done():
            log.Println("week2 exit...")
        case <-serverOut:
            log.Println("server will out...")
        }

        timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()

        log.Println("shutting down server...")
        return server.Shutdown(timeoutCtx)
    })


    g.Go(func() error {
        quit := make(chan os.Signal, 0)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

        select {
        case <-ctx.Done():
            return ctx.Err()
        case sig := <-quit:
            return fmt.Errorf("get os signal: %v", sig)
        }
    })

    fmt.Printf("week2 exiting: %+v\n", g.Wait())
}
