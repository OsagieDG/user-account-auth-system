package health

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func PostgreSQLHealthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("PostgreSQL database is unreachable"))
			if err != nil {
				log.Println("Error writing response:", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("PostgreSQL database is reachable"))
		if err != nil {
			log.Println("Error writing response:", err)
		}
	}
}
