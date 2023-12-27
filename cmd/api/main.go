package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sntegegn/url-shortening/internal/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	addr string
	dsn  string
}

type application struct {
	config      config
	logger      *slog.Logger
	formDecoder *form.Decoder
	URLModel    models.URLModelInterface
	tracer      trace.Tracer
}

func main() {
	ctx := context.Background()

	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "Server address")

	flag.StringVar(&cfg.dsn, "dsn", "postgres://urlshortening:pa55word@my-db/urlshortening?sslmode=disable", "PostgresSQL dsn")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	formDecoder := form.NewDecoder()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool successfully established")

	tp, err := tracerProvider(ctx, "jaeger:4318", "/v1/traces")
	if err != nil {
		logger.Error(err.Error(), "tp", tp)
		log.Fatal(err)
	}
	defer tp.Shutdown(ctx)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tr := otel.GetTracerProvider().Tracer("shortner")

	app := application{
		config:      cfg,
		logger:      logger,
		formDecoder: formDecoder,
		URLModel:    &models.URLModel{DB: db},
		tracer:      tr,
	}

	prometheus.Register(totalRequest)
	prometheus.Register(httpDuration)
	prometheus.Register(expandCount)

	srv := http.Server{
		Addr:         cfg.addr,
		Handler:      app.route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("Listening on: ", "addr", cfg.addr)

	err = srv.ListenAndServe()
	app.logger.Error(err.Error())
	os.Exit(1)
}
