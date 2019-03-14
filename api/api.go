//go:generate swagger generate spec

// Package api Payment API.
//
// This package contains all the realated features af the API. Which includes
// the logic, the data layer and the controllers.
//
// Terms Of Service:
//
// There are no TOS at this moment, use at your own risk we take no
// responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v1
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Extension:
//     x-go-name
//
// swagger:meta
package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/globalsign/mgo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Now() *time.Time {
	now := time.Now()
	return &now
}

var (
	mongo        *mgo.Session
	mongoHealthy bool
	config       *APIConfig
	store        PaymentStore
)

func Config() *APIConfig {
	return config
}

// NotFound is the default handler that is called when an unknown route is
// called. It will return the following body:
//   {
//     "data": {
//       "error": "Not found",
//       "code": "not_found"
//     },
//     "code": 404,
//     "status": "error"
//   }
func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewJSENDData(ErrNotFound))
}

// Ping is a simple route that will return 204 No Content in case everything
// is working fine
func Ping(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

// datasourceHealthy is a middleware aiming on blocking the api calls if the
// data source is not healthy: If it's not healthy it will return the following
// body:
//   {
//     "data": {
//       "error": "Maintenance is being done on the API",
//       "code": "undergoing_maintenance"
//     },
//     "code": 503,
//     "status": "error"
//   }
func datasourceHealthy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch config.DBType {
		case DatabaseTypeMongo:
			if mongo == nil || !mongoHealthy {
				handleError(w, r, ErrAPIMaintainance)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// paymentContext is a middleware that will try to fetch the payment from the
// data source and inject it in the request's context.
//
// On failure it will stop the chain and return the error in the body.
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

// swagger:route GET /payments payments listPayments
//
// Lists payments with pagination
//
// This will show a list of payments stored in the database
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: paymentList
func ListPayments(w http.ResponseWriter, r *http.Request) {
	limit, offset := readLimOff(r)
	ret, err := store.GetMany(limit, offset)
	if err != nil {
		handleError(w, r, err)
		return
	}
	render.Render(w, r, NewJSENDData(ret, http.StatusOK))
}

// swagger:route GET /payments/{id} payments getPayment
//
// Retrieves a single payment
//
//     Consumes:
//     	- application/json
//
//     Produces:
//     	- application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: singlePayment
//       404: reqError
//       400: reqError
func GetPayment(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewJSENDData(
		r.Context().Value("payment"),
		http.StatusOK,
	))
}

// SavePaymentReq is the payload for a payment creation request
type SavePaymentReq struct {
	*Payment
}

func NewSavePaymentReq() *SavePaymentReq {
	return &SavePaymentReq{NewPayment()}
}

func (p *SavePaymentReq) Bind(req *http.Request) error {
	// TODO: Validate
	return nil
}

// SavePayment will read the request's body and create or update a payment in
// the data source
// swagger:route POST /payments/{id} payments savePayment
//
// Creates or update a payment. When id is specified, updates the given payment
//
// Responses:
//    201: singlePayment
//		200: singlePayment
func SavePayment(w http.ResponseWriter, r *http.Request) {
	code := http.StatusCreated
	payload := NewSavePaymentReq()
	if err := render.Bind(r, payload); err != nil {
		handleError(w, r, ErrInvalidInput)
		return
	}
	if pCtx, ok := r.Context().Value("payment").(*Payment); ok {
		code = http.StatusOK
		payload.Payment.ID = pCtx.ID
		payload.CreatedAt = pCtx.CreatedAt
		payload.UpdatedAt = pCtx.UpdatedAt
	}
	if err := store.Save(payload.Payment); err != nil {
		handleError(w, r, err)
		return
	}
	render.Render(w, r, NewJSENDData(payload, code))
}

// DeletePayment removes a payment from the datasource
// swagger:route DELETE /payments/{id} payments deletePayment
//
// Deletes a pet from the store.
//
// Responses:
//    default: reqError
//        204:
func DeletePayment(w http.ResponseWriter, r *http.Request) {
	payment := r.Context().Value("payment").(*Payment)
	if err := store.Delete(payment.ID); err != nil {
		handleError(w, r, err)
		return
	}
	render.NoContent(w, r)
}

// Routes initializes the multiplexer and returns the http.Handler
func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType(strings.Split(APIV1ContentTypes, ",")...))
	if config != nil {
		cors := cors.New(cors.Options{
			AllowedOrigins:   config.Cors.AllowedOrigins,
			AllowedMethods:   config.Cors.AllowedMethods,
			AllowedHeaders:   config.Cors.AllowedHeaders,
			AllowCredentials: true,
			MaxAge:           300,
		})
		r.Use(cors.Handler)
		if config.DevMode {
			logrus.Info("Dev mode enabled")
			r.Use(middleware.Logger)
		}
	}
	r.Use(datasourceHealthy)
	r.NotFound(NotFound)
	r.Route(APIV1Prefix, func(r chi.Router) {
		r.Get("/ping", Ping)
		r.Route("/payments", func(r chi.Router) {
			r.Use()
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
	})
	return r
}

// Start wil take care of creating and starting the server with the given config
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
		if mongo != nil {
			logrus.Info("Closing connection to Mongo")
			mongo.Close()
		}
		close(done)
	}()
	var err error
	logrus.Infof("Starting server on %s", config.GetHostURL())
	if config.GetTLS() {
		logrus.Info("TLS Enabled")
		logrus.Debug("Loading cert: ", config.TLSCert)
		logrus.Debug("Loading key: ", config.TLSKey)
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
