package context

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func paymentService(ctx context.Context) (string, error) {
	select {
	case <-time.After(100 * time.Millisecond):
		return "Payment completed", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func userService(ctx context.Context) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond):
		return "User fetched", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func orderService(ctx context.Context) (string, error) {
	select {
	case <-time.After(10 * time.Second):
		return "Order fetched", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	results := make(chan string, 3)

	var wg sync.WaitGroup

	services := []func(context.Context) (string, error){
		paymentService,
		userService,
		orderService,
	}

	for _, service := range services {
		wg.Add(1)

		go func(fn func(context.Context) (string, error)) {
			defer wg.Done()

			result, err := fn(ctx)
			if err != nil {
				return
			}
			results <- result

		}(service)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var responses []string

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Client disconnected...", ctx.Err())
			return
		case result, ok := <-results:
			if !ok {
				fmt.Fprintln(w, responses)
				return
			}
			responses = append(responses, result)
		}
	}
}

func main_context() {
	mux := http.NewServeMux()

	// 2. Wrap your handler with http.TimeoutHandler to enforce a 300ms execution timeout
	timeoutHandler := http.TimeoutHandler(
		http.HandlerFunc(handler),
		300*time.Millisecond,
		"Request timed out", // Custom message returned to the client on timeout
	)
	mux.Handle("/api", timeoutHandler)

	// 3. Configure the underlying server with network/connection-level timeouts
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,        // Protects against Slowloris attacks
		WriteTimeout:      500 * time.Millisecond, // Slightly longer than execution timeout to ensure response delivery
		IdleTimeout:       30 * time.Second,
	}

	// 4. Start the configured server
	srv.ListenAndServe()
}
