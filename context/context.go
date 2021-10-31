package context

import (
	"context"
	"fmt"
	"net/http"
)

// func Server(store Store) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()

// 		data := make(chan string, 1)

// 		go func() {
// 			data <- store.Fetch()
// 		}()

// 		select {
// 		case d := <- data: // dataが来たらプリント
// 			fmt.Fprint(w, d)
// 		case <-ctx.Done(): // ctx がDoneになったら storeをCancel. context has a method Done() which returns a channel which gets sent a signal when the context is "done" or "cancelled".
// 			store.Cancel()
// 		}
//     }
// }

func Server(store Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

		if err != nil {
			return // TODO: log error however you like
		}

		fmt.Fprint(w, data)
    }
}

// type Store interface {
//     Fetch() string
//     Cancel()
// }

type Store interface {
    Fetch(ctx context.Context) (string, error)
}
