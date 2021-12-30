package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/captcha"
	petname "github.com/dustinkirkland/golang-petname"
	log "github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"

	emailverifier "github.com/AfterShip/email-verifier"
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
	guestVlanId   int
	corpVlanId    int
	ssid          string
	emailSender   string
	emailAdmin    string
	mailgunApiKey string
	baseUrl       string

	verifier = emailverifier.NewVerifier().EnableAutoUpdateDisposable()
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

type captchaResp struct {
	CaptchaId string `json:"captcha_id"`
}

type registerReq struct {
	Email         string `json:"email"`
	CaptchaId     string `json:"captcha_id"`
	CaptchaAnswer string `json:"captcha_answer"`
	CorpAccess    bool   `json:"corp_access"`
}

type registerResp struct {
	Success      bool              `json:"success"`
	Message      string            `json:"message"`
	InputErrors  map[string]string `json:"input_errors"`
	Email        string            `json:"email"`
	ValidForDays int               `json:"valid_for_days"`
	CorpAccess   bool              `json:"corp_access"`
}

type approveResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

var ctx = context.Background()

func redisHandler(c *redis.Client,
	f func(c *redis.Client, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(c, w, r) })
}

func getCaptcha(w http.ResponseWriter, r *http.Request) {
	resp := captchaResp{
		CaptchaId: captcha.New(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func sendRegisterResponse(w http.ResponseWriter, statusCode int, success bool,
	message string, inputErrors map[string]string, email string, validForDays int,
	corpAccess bool) {
	resp := registerResp{
		Success:      success,
		Message:      message,
		InputErrors:  inputErrors,
		Email:        email,
		ValidForDays: validForDays,
		CorpAccess:   corpAccess,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func sendApproveResponse(w http.ResponseWriter, statusCode int, success bool,
	message string, email string) {
	resp := approveResp{
		Success: success,
		Message: message,
		Email:   email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func registerAccount(c *redis.Client, w http.ResponseWriter, r *http.Request) {
	var rr registerReq

	inputErrors := make(map[string]string)

	err := json.NewDecoder(r.Body).Decode(&rr)
	if err != nil {
		log.Errorf("register: unable to decode JSON request: %s\n", err)
		sendRegisterResponse(w, http.StatusInternalServerError, false,
			"Internal server error, please try again.", inputErrors, rr.Email, 0,
			rr.CorpAccess)
		return
	}

	ret, err := verifier.Verify(rr.Email)
	if err != nil {
		switch {
		case !ret.HasMxRecords:
			log.Errorf("register: no MX record for domain, email: %s\n", rr.Email)
			inputErrors["email"] = "No MX record for domain"
		default:
			log.Errorf("register: unable to verify email: %s\n", err)
			sendRegisterResponse(w, http.StatusInternalServerError, false,
				"Internal server error, please try again.", inputErrors, rr.Email, 0,
				rr.CorpAccess)
			return
		}
	}

	if !ret.Syntax.Valid {
		log.Errorf("register: invalid email, email: %s\n", rr.Email)
		inputErrors["email"] = "Invalid email address"
	}

	if ret.Disposable {
		log.Errorf("register: use of disposable email, email: %s\n", rr.Email)
		inputErrors["email"] = "No disposable email please"
	}

	if ret.RoleAccount {
		log.Errorf("register: use of role email address, email: %s\n", rr.Email)
		inputErrors["email"] = "No role email address please"
	}

	if !captcha.VerifyString(rr.CaptchaId, rr.CaptchaAnswer) {
		log.Errorf("register: wrong captcha, email: %s\n", rr.Email)
		inputErrors["captchaAnswer"] = "Wrong CAPTCHA answer"
	}

	if len(inputErrors) > 0 {
		sendRegisterResponse(w, http.StatusBadRequest, false,
			"", inputErrors, rr.Email, 0, rr.CorpAccess)
		return
	}

	var pending string
	var vlanId int
	if rr.CorpAccess {
		pending = ":pending"
		vlanId = corpVlanId
	} else {
		pending = ""
		vlanId = guestVlanId
	}

	var pwd = ""
	emailKey := "guest:email:" + strings.ToLower(rr.Email) + pending
	duration := time.Duration(pwdExpiration) * 24 * time.Hour

	// Use existing password if exists, otherwise generate a new password.
	pwd, err = c.Get(ctx, emailKey).Result()
	if err != nil {
		pwd = petname.Generate(2, "_")
	}

	err = c.Set(ctx, emailKey, pwd, duration).Err()
	if err != nil {
		log.Errorf("register: failed to write key to redis: %s\n", err)
		sendRegisterResponse(w, http.StatusInternalServerError, false,
			"Internal server error, please try again.", inputErrors, rr.Email, 0,
			rr.CorpAccess)
		return
	}

	vlanKey := "guest:vlan:" + strings.ToLower(rr.Email)
	err = c.Set(ctx, vlanKey, strconv.Itoa(vlanId), duration).Err()
	if err != nil {
		log.Errorf("register: failed to write key to redis: %s\n", err)
		sendRegisterResponse(w, http.StatusInternalServerError, false,
			"Internal server error, please try again.", inputErrors, rr.Email, 0,
			rr.CorpAccess)
		return
	}

	if rr.CorpAccess {
		requestId := randstr.String(64)
		approvalKey := "guest:approval:" + requestId
		err = c.Set(ctx, approvalKey, strings.ToLower(rr.Email), duration).Err()
		if err != nil {
			log.Errorf("register: failed to write key to redis: %s\n", err)
			sendRegisterResponse(w, http.StatusInternalServerError, false,
				"Internal server error, please try again.", inputErrors, rr.Email, 0,
				rr.CorpAccess)
			return
		}

		approvalLink := baseUrl + "/api/v1/approve?id=" + requestId
		sendApprovalMail(rr.Email, approvalLink)
		sendRegisterResponse(w, http.StatusOK, true,
			"Account is under review.", inputErrors, rr.Email, pwdExpiration,
			rr.CorpAccess)
	} else {
		sendRegSuccessMail(pwd, rr.Email)
		sendRegisterResponse(w, http.StatusOK, true,
			"Account successfully registered.", inputErrors, rr.Email, pwdExpiration,
			rr.CorpAccess)
	}
}

func approveAccount(c *redis.Client, w http.ResponseWriter, r *http.Request) {
	requestId := r.URL.Query().Get("id")
	if requestId == "" {
		sendApproveResponse(w, http.StatusBadRequest, false,
			"Invalid request ID.", "")
		return
	}

	approvalKey := "guest:approval:" + requestId
	email, err := c.Get(ctx, approvalKey).Result()
	if err != nil {
		log.Errorf("approve: failed to get key %s in redis: %s\n",
			approvalKey, err)
		sendApproveResponse(w, http.StatusInternalServerError, false,
			"Internal server error.", "")
		return
	}

	oldKey := "guest:email:" + email + ":pending"
	newKey := "guest:email:" + email
	if err = c.Rename(ctx, oldKey, newKey).Err(); err != nil {
		log.Errorf("approve: failed to rename key %s to %s in redis: %s\n",
			oldKey, newKey, err)
		sendApproveResponse(w, http.StatusInternalServerError, false,
			"Internal server error.", "")
		return
	}

	pwd, err := c.Get(ctx, newKey).Result()
	if err != nil {
		log.Errorf("approve: failed to get key %s in redis: %s\n", err)
		sendApproveResponse(w, http.StatusInternalServerError, false,
			"Internal server error.", "")
		return
	}

	if err = c.Del(ctx, approvalKey).Err(); err != nil {
		log.Errorf("approve: failed to delete key %s in redis: %s\n",
			approvalKey, err)
		sendApproveResponse(w, http.StatusInternalServerError, false,
			"Internal server error.", "")
		return
	}

	sendRegSuccessMail(pwd, email)
	sendApproveResponse(w, http.StatusOK, true, "Request approved.", email)
}

func sendRegSuccessMail(pwd string, recipient string) {
	parts := strings.Split(emailSender, "@")
	mg := mailgun.NewMailgun(parts[1], mailgunApiKey)
	sender := emailSender
	subject := "Bektinet Guest Wi-Fi"
	body := fmt.Sprintf("Hello,\n\n"+
		"Thank you for registering with Bektinet. You may access the guest Wi-Fi by using "+
		"the information below.\n\n"+
		"SSID: "+ssid+"\n"+
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

func sendApprovalMail(requesterEmail string, approvalLink string) {
	parts := strings.Split(emailSender, "@")
	mg := mailgun.NewMailgun(parts[1], mailgunApiKey)
	sender := emailSender
	subject := "[Request] Bektinet Corp Wi-Fi"
	body := fmt.Sprintf("Hello,\n\n"+
		"The following user is requesting access to the corp Wi-Fi.\n\n"+
		"Email: %s\n\n"+
		"Click this link to approve the request: %s\n\n"+
		"Thank you.", requesterEmail, approvalLink)
	message := mg.NewMessage(sender, subject, body, emailAdmin)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Email sent to %s, ID: %s Resp: %s\n", emailAdmin, id, resp)
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
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
		cli.StringFlag{
			Name:        "ssid",
			Value:       "",
			Usage:       "ssid",
			EnvVar:      "SSID",
			Destination: &ssid,
		},
		cli.IntFlag{
			Name:        "guest-vlan-id",
			Value:       10,
			Usage:       "guest vlan id",
			EnvVar:      "GUEST_VLAN_ID",
			Destination: &guestVlanId,
		},
		cli.IntFlag{
			Name:        "corp-vlan-id",
			Value:       0,
			Usage:       "corp vlan id",
			EnvVar:      "CORP_VLAN_ID",
			Destination: &corpVlanId,
		},
		cli.StringFlag{
			Name:        "email-sender",
			Value:       "",
			Usage:       "email sender",
			EnvVar:      "EMAIL_SENDER",
			Destination: &emailSender,
		},
		cli.StringFlag{
			Name:        "email-admin",
			Value:       "",
			Usage:       "email admin",
			EnvVar:      "EMAIL_ADMIN",
			Destination: &emailAdmin,
		},
		cli.StringFlag{
			Name:        "mailgun-api-key",
			Value:       "",
			Usage:       "mailgun api key",
			EnvVar:      "MAILGUN_API_KEY",
			Destination: &mailgunApiKey,
		},
		cli.StringFlag{
			Name:        "base-url",
			Value:       "",
			Usage:       "base url",
			EnvVar:      "BASE_URL",
			Destination: &baseUrl,
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

		log.Infof("server listening on port %d\n", bindPort)
		router := mux.NewRouter()
		router.HandleFunc("/api/v1/captcha", getCaptcha).Methods("GET")
		router.Handle("/api/v1/register", redisHandler(rdb, registerAccount)).Methods("POST")
		router.Handle("/api/v1/approve", redisHandler(rdb, approveAccount)).Methods("GET")
		router.Methods("GET").PathPrefix("/captcha/").Handler(captcha.Server(captcha.StdWidth, captcha.StdHeight))

		spa := spaHandler{staticPath: "build", indexPath: "index.html"}
		router.PathPrefix("/").Handler(spa)

		srv := &http.Server{
			Handler:      router,
			Addr:         bindAddr + ":" + strconv.Itoa(bindPort),
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  30 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
