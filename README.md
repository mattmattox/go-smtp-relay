# GO-SMTP-RELAY

This project is a **SMTP relay server** designed to forward emails to a configured SMTP server using authentication. It logs SMTP interactions, provides **Prometheus metrics** for monitoring, and exposes a **health check endpoint** to verify the server's status. The relay server listens on port 25 for incoming emails with no authentication required, overwrites the `From` address, and forwards the email to the configured SMTP server for delivery.

## Overview

The project consists of the following key components:

### **`pkg/version/version.go`**
- Contains version information for the application.
- Provides the `GetVersion` function to return the version as JSON.

### **`pkg/config/config.go`**
- Manages the configuration of the application.
- Defines the `AppConfig` structure to load settings from environment variables and command line flags.

### **`pkg/logging/logging.go`**
- Handles structured logging using [Logrus](https://github.com/sirupsen/logrus).
- Includes the `SetupLogging` function to initialize the logger and configure it based on the `debug` flag.

### **`main.go`**
- The main entry point of the SMTP relay server.
- Sets up the SMTP backend, Prometheus metrics, and HTTP server for health checks.
- Starts both the **SMTP server** and **metrics/health check server**.

## Getting Started

1. **Set up Configuration:**
   - Use environment variables or command-line flags to configure the application.
     Example environment variables:
     ```bash
     export EMAIL_SERVER_HOST="mail.support.tools"
     export EMAIL_SERVER_PORT=587
     export EMAIL_SERVER_USER="no-reply@support.tools"
     export EMAIL_SERVER_PASS="yourpassword"
     export FROM_ADDRESS="no-reply@support.tools"
     export SERVER_PORT=25
     export DEBUG=true
     ```

2. **Run the Application:**
   ```bash
   sudo go run main.go
   ```

3. **Check Version:**
   - Use the `-version` flag to display version information:
     ```bash
     go run main.go -version
     ```

## Prometheus Metrics and Health Checks

The SMTP relay server exposes the following endpoints:

- **Metrics Endpoint:** [http://localhost:2112/metrics](http://localhost:2112/metrics)  
  Provides detailed metrics about received, forwarded, and failed emails.
  - `smtp_emails_received_total`: Total number of emails received.
  - `smtp_emails_forwarded_total`: Total number of emails successfully forwarded.
  - `smtp_emails_failed_total`: Total number of emails that failed to forward.

- **Health Check Endpoint:** [http://localhost:2112/healthz](http://localhost:2112/healthz)  
  Returns `200 OK` if the service is running and healthy.

## Usage

### **Configure Environment Variables:**
Make sure to set the appropriate SMTP server credentials and configuration settings.

### **Run the Application:**
- Run the application directly using Go:
  ```bash
  sudo go run main.go
  ```

- **Using Docker:**
  ```bash
  docker run -d -p 25:25 -p 2112:2112 \
    -e EMAIL_SERVER_HOST=mail.support.tools \
    -e EMAIL_SERVER_PORT=587 \
    -e EMAIL_SERVER_USER=no-reply@support.tools \
    -e EMAIL_SERVER_PASS=yourpassword \
    -e FROM_ADDRESS=no-reply@support.tools \
    -e SERVER_PORT=25 \
    --name go-smtp-relay \
    cube8021/go-smtp-relay
  ```

### **Testing SMTP Relay with Telnet:**
1. Open a telnet connection:
   ```bash
   telnet localhost 25
   ```

2. Send a test email:
   ```
   EHLO localhost
   MAIL FROM:<sender@example.com>
   RCPT TO:<recipient@example.com>
   DATA
   Subject: Test Email

   This is a test message.
   .
   QUIT
   ```

## Security

- **Do not expose the SMTP relay to the public internet** without proper security controls. The relay server does not enforce authentication for incoming SMTP connections and could be exploited if left open to the public. So make sure to run the relay server behind a firewall or VPN.

## Dependencies

- [Logrus](https://github.com/sirupsen/logrus): A structured logging library for Go.
- [Prometheus Go Client](https://github.com/prometheus/client_golang): Used for exposing application metrics.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.
