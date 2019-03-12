package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	config = NewAPIConfig()
	store  PaymentStore
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewJSENDData(ErrNotFound))
}

func paymentContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paymentID, err := uuid.Parse(chi.URLParam(r, "paymentID"))
		if err != nil {
			render.Render(w, r, NewJSENDData(ErrInvalidInput))
			return
		}
		payment, err := store.GetByID(paymentID)
		if err != nil {
			handleError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), "payment", payment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if apiErr, ok := err.(*APIError); ok {
		render.Render(w, r, NewJSENDData(apiErr))
		return
	}
	render.Render(w, r, NewJSENDData(ErrSomethingWentWrong(err)))
}

func ListPayments(w http.ResponseWriter, r *http.Request) {
	limit, offset := readLimOff(r)
	ret, err := store.GetMany(limit, offset)
	if err != nil {
		handleError(w, r, err)
		return
	}
	render.Render(w, r, NewJSENDData(ret, http.StatusOK))
}

func GetPayment(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewJSENDData(
		r.Context().Value("payment"),
		http.StatusOK,
	))
}

type SavePaymentReq struct {
	*Payment
}

func NewSavePaymentReq() *SavePaymentReq {
	return &SavePaymentReq{NewPayment()}
}

func (p *SavePaymentReq) Bind(req *http.Request) error {
	// TODO: Validate
	p.CreatedAt = nil
	p.UpdatedAt = nil
	return nil
}

func SavePayment(w http.ResponseWriter, r *http.Request) {
	var payment *Payment
	p := NewSavePaymentReq()
	if err := render.Bind(r, p); err != nil {
		handleError(w, r, ErrInvalidInput)
		return
	}
	if pCtx := r.Context().Value("payment"); pCtx == nil {
		payment = p.Payment
	} else {
		payment = pCtx.(*Payment)
		p.CreatedAt = payment.CreatedAt
		p.UpdatedAt = payment.UpdatedAt
	}
	if err := store.Save(p.Payment); err != nil {
		handleError(w, r, err)
		return
	}
	render.Render(w, r, NewJSENDData(p))
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	payment := r.Context().Value("payment").(*Payment)
	if err := store.Delete(payment.ID); err != nil {
		handleError(w, r, err)
		return
	}
	render.NoContent(w, r)
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
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType(strings.Split(APIV1ContentTypes, ",")...))
	if !config.DevMode {
		r.Use(middleware.Logger)
	}
	r.NotFound(NotFound)
	r.Route("/payments", func(r chi.Router) {
		r.Get(URLRoot, ListPayments)
		r.Post(URLRoot, SavePayment)
		r.Route("/{paymentID}", func(r chi.Router) {
			r.Use(paymentContext)
			r.Get(URLRoot, GetPayment)
			r.Put(URLRoot, SavePayment)
			r.Post(URLRoot, SavePayment)
			r.Delete(URLRoot, DeletePayment)
		})

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
