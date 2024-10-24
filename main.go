package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"strings"

	smtpServer "github.com/emersion/go-smtp"
	"github.com/mattmattox/go-smtp-relay/pkg/config"
	"github.com/mattmattox/go-smtp-relay/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logging.SetupLogging()

// Prometheus metrics
var (
	emailsReceived = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "smtp_emails_received_total",
			Help: "Total number of emails received by the SMTP relay.",
		},
	)
	emailsForwarded = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "smtp_emails_forwarded_total",
			Help: "Total number of emails successfully forwarded.",
		},
	)
	emailsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "smtp_emails_failed_total",
			Help: "Total number of emails that failed to forward.",
		},
	)
)

// Backend defines the backend to handle SMTP sessions.
type Backend struct{}

// NewSession creates a new session for each connection.
func (bkd *Backend) NewSession(c *smtpServer.Conn) (smtpServer.Session, error) {
	session := &Session{
		from: "",
		to:   []string{},
	}
	return session, nil
}

// Session stores the state of an SMTP session.
type Session struct {
	from string
	to   []string
	data strings.Builder
}

// Mail handles the MAIL FROM command.
func (s *Session) Mail(from string, opts *smtpServer.MailOptions) error {
	log.Infof("Mail from: %s", from)
	// Override the from address using the config value.
	s.from = config.CFG.FromAddress
	return nil
}

// Rcpt handles the RCPT TO command.
func (s *Session) Rcpt(to string, opts *smtpServer.RcptOptions) error {
	log.Infof("Rcpt to: %s", to)
	s.to = append(s.to, to)
	return nil
}

// Data handles the DATA command and stores the email content.
func (s *Session) Data(r io.Reader) error {
	log.Info("Receiving email data...")
	_, err := io.Copy(&s.data, r)
	if err != nil {
		return err
	}

	// Forward the email using the SMTP server details from the config.
	err = forwardEmail(s.from, s.to, s.data.String())
	if err != nil {
		log.Errorf("Failed to forward email: %v", err)
		return err
	}

	log.Info("Email forwarded successfully.")
	return nil
}

// Reset resets the session state.
func (s *Session) Reset() {
	log.Info("Resetting session")
	s.data.Reset()
	s.from = ""
	s.to = []string{}
}

// Logout logs the client out.
func (s *Session) Logout() error {
	log.Info("Logging out")
	return nil
}

// forwardEmail forwards the email using SMTP with authentication.
func forwardEmail(from string, to []string, body string) error {
	smtpHost := config.CFG.EmailServerHost
	smtpPort := config.CFG.EmailServerPort
	username := config.CFG.EmailServerUser
	password := config.CFG.EmailServerPass

	auth := smtp.PlainAuth("", username, password, smtpHost)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\n%s", from, strings.Join(to, ","), body))

	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// healthHandler responds with 200 OK to indicate the service is healthy.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// setupMetrics registers Prometheus metrics.
func setupMetrics() {
	prometheus.MustRegister(emailsReceived)
	prometheus.MustRegister(emailsForwarded)
	prometheus.MustRegister(emailsFailed)
}

func main() {
	// Load configuration
	config.LoadConfiguration()

	// Set up Prometheus metrics
	setupMetrics()

	// Set up HTTP server for metrics and health check endpoints
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/healthz", healthHandler)

	// Start the HTTP server for metrics and health check in a goroutine
	go func() {
		metricsPort := config.CFG.MetricsPort
		log.Infof("Starting metrics and health check server on :%d...", metricsPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil); err != nil {
			log.Fatalf("Error starting metrics server: %v", err)
		}
	}()

	backend := &Backend{}
	server := smtpServer.NewServer(backend)

	server.Addr = fmt.Sprintf(":%d", config.CFG.SmtpPort) // Use port from config.
	server.Domain = "localhost"
	server.AllowInsecureAuth = true // No TLS or authentication required.

	// Start the SMTP server
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Failed to bind to port %d: %v", config.CFG.SmtpPort, err)
	}

	log.Infof("Starting SMTP server on port %d...", config.CFG.SmtpPort)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error starting SMTP server: %v", err)
	}
}
