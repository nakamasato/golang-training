package context

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubStore struct {
    response string
}

func (s *StubStore) Fetch() string {
    return s.response
}

type SpyStore struct {
    response string
	t *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		// intentionally make the process slow
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				s.t.Log("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()
    select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}

// Our SpyResponseWriter implements http.ResponseWriter
type SpyResponseWriter struct {
    written bool
}

func (s *SpyResponseWriter) Header() http.Header {
    s.written = true
    return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
    s.written = true
    return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
    s.written = true
}

func TestServer(t *testing.T) {
    data := "hello, world"
    // svr := Server(&StubStore{data})

    // request := httptest.NewRequest(http.MethodGet, "/", nil)
    // response := httptest.NewRecorder()

    // svr.ServeHTTP(response, request)

    // if response.Body.String() != data {
    //     t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    // }

    t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
        store := &SpyStore{response: data, t: t}
        svr := Server(store)
        request := httptest.NewRequest(http.MethodGet, "/", nil)
        response := httptest.NewRecorder()
        svr.ServeHTTP(response, request)
        if response.Body.String() != data {
            t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
        }
    })

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		cancellingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingCtx)
		response := &SpyResponseWriter{}
		svr.ServeHTTP(response, request)
		if response.written {
			t.Error("a response should not have been written")
		}
	})
}
