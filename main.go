package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	log "github.com/sirupsen/logrus"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/urfave/cli"
)

var (
	bindAddr      string
	bindPort      int
	redisAddr     string
	redisPwd      string
	logLevel      string
	pwdExpiration int
	vlanId        int
	emailSender   string
	mailgunApiKey string
)

var ctx = context.Background()

func redisHandler(c *redis.Client,
	f func(c *redis.Client, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(c, w, r) })
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	tmpl.Execute(w, data)
}

func terms(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/terms.html"))
	tmpl.Execute(w, nil)
}

func register(c *redis.Client, w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if _, err := mail.ParseAddress(email); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("register: invalid email, email: %s\n", email)
		tmpl := template.Must(template.ParseFiles("templates/invalid_email.html"))
		tmpl.Execute(w, nil)
		return
	}

	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("register: wrong captcha, email: %s\n", email)
		tmpl := template.Must(template.ParseFiles("templates/wrong_captcha.html"))
		tmpl.Execute(w, nil)
		return
	}

	pwd := namesgenerator.GetRandomName(1)
	duration := time.Duration(pwdExpiration) * 24 * time.Hour

	emailKey := "guest:email:" + strings.ToLower(email)
	err := c.Set(ctx, emailKey, pwd, duration).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error, please try again.")
		log.Errorf("register: failed to write key to redis: %s\n", err)
		return
	}

	vlanKey := "guest:vlan:" + strings.ToLower(email)
	err = c.Set(ctx, vlanKey, strconv.Itoa(vlanId), duration).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error, please try again.")
		log.Errorf("register: failed to write key to redis: %s\n", err)
		return
	}

	sendMail(pwd, email)

	tmpl := template.Must(template.ParseFiles("templates/success.html"))
	data := struct {
		Email string
	}{
		email,
	}
	tmpl.Execute(w, data)
}

func sendMail(pwd string, recipient string) {
	parts := strings.Split(emailSender, "@")

	mg := mailgun.NewMailgun(parts[1], mailgunApiKey)

	sender := emailSender
	subject := "Bektinet Guest Wi-Fi"
	body := fmt.Sprintf("Hello,\n\n"+
		"Thank you for registering with Bektinet. You may access the guest Wi-Fi by using "+
		"the information below.\n\n"+
		"SSID: bektinet-wpa\n"+
		"Username: %s\n"+
		"Password: %s\n\n"+
		"You can use the Wi-Fi for up to %d days. It will expire after that, and you will "+
		"need to register again from the website.\n\n"+
		"We hope you enjoy your stay with us.", recipient, pwd, pwdExpiration)

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Email sent to %s, ID: %s Resp: %s\n", recipient, id, resp)
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "bind-addr",
			Value:       "0.0.0.0",
			Usage:       "address to bind to",
			EnvVar:      "BIND_ADDR",
			Destination: &bindAddr,
		},
		&cli.IntFlag{
			Name:        "bind-port",
			Value:       8080,
			Usage:       "port number to bind to",
			EnvVar:      "BIND_PORT",
			Destination: &bindPort,
		},
		cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "log level",
			EnvVar:      "LOG_LEVEL",
			Destination: &logLevel,
		},
		cli.StringFlag{
			Name:        "redis-addr",
			Value:       "localhost:6379",
			Usage:       "redis address",
			EnvVar:      "REDIS_ADDR",
			Destination: &redisAddr,
		},
		cli.StringFlag{
			Name:        "redis-pwd",
			Value:       "",
			Usage:       "redis password",
			EnvVar:      "REDIS_PWD",
			Destination: &redisPwd,
		},
		cli.IntFlag{
			Name:        "pwd-expiration",
			Value:       3,
			Usage:       "password expiration in days",
			EnvVar:      "PWD_EXPIRATION",
			Destination: &pwdExpiration,
		},
		cli.IntFlag{
			Name:        "vlan-id",
			Value:       0,
			Usage:       "vlan id",
			EnvVar:      "VLAN_ID",
			Destination: &vlanId,
		},
		cli.StringFlag{
			Name:        "email-sender",
			Value:       "",
			Usage:       "email sender",
			EnvVar:      "EMAIL_SENDER",
			Destination: &emailSender,
		},
		cli.StringFlag{
			Name:        "mailgun-api-key",
			Value:       "",
			Usage:       "mailgun api key",
			EnvVar:      "MAILGUN_API_KEY",
			Destination: &mailgunApiKey,
		},
	}

	app.Action = func(c *cli.Context) error {
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)

		rdb := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPwd,
			DB:       0, // use default DB
		})

		log.Infof("welcome listening on port %d\n", bindPort)
		router := mux.NewRouter()
		router.HandleFunc("/", index).Methods("GET")
		router.HandleFunc("/terms", terms).Methods("GET")
		router.Handle("/register", redisHandler(rdb, register)).Methods("POST")
		router.Methods("GET").PathPrefix("/captcha/").Handler(captcha.Server(captcha.StdWidth, captcha.StdHeight))
		log.Fatal(http.ListenAndServe(bindAddr+":"+strconv.Itoa(bindPort), router))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
