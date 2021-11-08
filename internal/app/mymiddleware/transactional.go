package mymiddleware

import (
	"context"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
	"go.uber.org/zap"
	"net/http"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Transactional(handler basedbhandler.Transactioner, log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("start transaction")
			tx, err := handler.NewTx(context.Background())
			if err != nil {
				log.Error("TransactionMiddleware: can't start transaction", zap.Error(err))
				return
			}
			defer func() {
				if r := recover(); r != nil {
					log.Error("TransactionMiddleware: Panic. Try to rollback", zap.Error(err))
					err = tx.Rollback(context.Background())
					fmt.Println(err)
				}
			}()
			ctx := context.WithValue(r.Context(), "tx", tx)

			sw := statusWriter{ResponseWriter: w}
			next.ServeHTTP(&sw, r.WithContext(ctx))
			if sw.status > http.StatusNoContent {
				fmt.Println("rollback transaction")
				if err := tx.Rollback(ctx); err != nil {
					log.Error("TransactionMiddleware: Can't commit", zap.Error(err))
				}
			} else {
				fmt.Println("commit transaction")
				if err := tx.Commit(ctx); err != nil {
					log.Error("TransactionMiddleware: Can't rollback", zap.Error(err))
				}
			}
		})
	}
}
