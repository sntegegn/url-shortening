package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/sntegegn/url-shortening/internal/mailer"
)

type config struct {
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	addr string
}

type application struct {
	config      config
	mailer      mailer.Mailer
	logger      *slog.Logger
	urls        map[string]string
	formDecoder *form.Decoder
}

func main() {
	var cfg config
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "3d7205e4cb840f", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "5c0655e3404e1c", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "URLShortner <no-reply@urlshortner.selam.com>", "SMTP sender")

	flag.StringVar(&cfg.addr, "addr", "localhost:4000", "Server address")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	formDecoder := form.NewDecoder()

	app := application{
		config:      cfg,
		mailer:      mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		logger:      logger,
		urls:        map[string]string{},
		formDecoder: formDecoder,
	}

	srv := http.Server{
		Addr:         cfg.addr,
		Handler:      app.route(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("Listening on: ", "addr", cfg.addr)

	err := srv.ListenAndServe()
	app.logger.Error(err.Error())
	os.Exit(1)

	//app.shortenURL("https://stackoverflow.com/questions/28322997/how-to-get-a-list-of-values-into-a-flag-in-golang")
	//fmt.Println(app.urls)
	//app.sendEmail()

}
