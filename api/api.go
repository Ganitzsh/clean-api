package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

var (
	config = NewAPIConfig()
	store  DocumentStore
)

type JSENDData struct {
	Data   interface{} `json:"data"`
	Code   int         `json:"code"`
	Status string      `json:"status"`
}

func postProcess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data interface{}
		var status string
		var code int = http.StatusInternalServerError
		var ctx context.Context
		next.ServeHTTP(w, r)
		ctx = r.Context()
		if ctxCode, ok := ctx.Value(APIContextKeyCode).(int); ok {
			code = ctxCode
		}
		if err, ok := ctx.Value(APIContextKeyError).(error); ok {
			status = "error"
			if apiErr, ok := err.(*APIError); ok {
				if apiErr.DataError {
					status = "fail"
				}
				data = apiErr
			} else {
				data = &APIError{Message: err.Error()}
			}
		} else {
			status = "success"
			data = ctx.Value(APIContextKeyData)
		}
		render.WriteJSON(w, JSENDData{
			Data:   data,
			Code:   code,
			Status: status,
		})
	})
}

// func success(r *http.Request, code int, data interface{}) {
// 	ctx := NewAPIContext(r.Context()).
// 		SetCode(code).
// 		SetData(data)
// 	newReq := r.WithContext(ctx)
// 	*r = *newReq
// }

// func failure(r *http.Request, err error, code int) {
// 	ctx := NewAPIContext((*r).Context())
// 	if err == nil {
// 		ctx.
// 			SetCode(code).
// 			SetError(ErrSomethingWentWrong)
// 	} else {
// 		ctx.
// 			SetCode(code).
// 			SetError(err)
// 	}
// 	newReq := r.WithContext(ctx)
// 	*r = *newReq
// }

func success(w http.ResponseWriter, code int, data interface{}) {
	render.WriteJSON(w, JSENDData{
		Data:   data,
		Code:   code,
		Status: "success",
	})
}

func failure(w http.ResponseWriter, err error, code int) {
	var data interface{}
	status := "error"
	if apiErr, ok := err.(*APIError); ok {
		if apiErr.DataError {
			status = "fail"
		}
		data = apiErr
	} else {
		data = &APIError{Message: err.Error()}
	}
	w.WriteHeader(code)
	render.WriteJSON(w, JSENDData{
		Data:   data,
		Code:   code,
		Status: status,
	})
}

func readLimOff(r *http.Request) (lim int, off int) {
	if r == nil {
		return 0, 0
	}
	val, err := strconv.Atoi(r.URL.Query().Get("lim"))
	if err == nil {
		lim = val
	}
	val, err = strconv.Atoi(r.URL.Query().Get("off"))
	if err == nil {
		off = val
	}
	return lim, off
}

func ListPayments(w http.ResponseWriter, r *http.Request) {
	limit, offset := readLimOff(r)
	ret, err := store.GetMany(limit, offset)
	if err != nil {
		failure(w, err, http.StatusBadRequest)
	}
	success(w, http.StatusOK, ret)
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	failure(w, ErrNotImplemented, http.StatusBadRequest)
}

func SavePayment(w http.ResponseWriter, r *http.Request) {
	failure(w, ErrNotImplemented, http.StatusBadRequest)
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	failure(w, ErrNotImplemented, http.StatusBadRequest)
}

func Routes() http.Handler {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Use(middleware.AllowContentType(strings.Split(APIV1ContentTypes, ",")...))
	if !config.DevMode {
		r.Use(middleware.Logger)
	}
	r.Route("/payments", func(r chi.Router) {
		// r.Use(postProcess)
		r.Get("/", ListPayments)
	})
	return r
}

func Start() error {
	srv := http.Server{
		Addr:    config.GetFullHost(),
		Handler: Routes(),
	}
	done := make(chan bool)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.Fatalf("Could not shutdown: %v", err)
		}
		close(done)
	}()
	var err error
	logrus.Infof("Starting server on %s", srv.Addr)
	if config.GetTLS() {
		err = srv.ListenAndServeTLS(config.TLSCert, config.TLSKey)
	} else {
		err = srv.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		return fmt.Errorf("Error: %s", err)
	}
	logrus.Info("Shutting down...")
	<-done
	return nil
}
