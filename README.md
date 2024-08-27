# Go-Based Server Monitoring Tool

<p style="text-align: center;">
  <img src="https://i.imgur.com/b67VT1v.png" alt="Moniping" width="300"/>
</p>

## Overview

This project is a Go-based service designed to monitor the availability of servers by periodically pinging them and sending alerts via email when a server goes down or comes back up. The tool is configurable via a JSON file, supports concurrent monitoring of multiple servers, and includes features such as email notifications, logging, and rate limiting.

## Features

- **Server Monitoring**: Periodically pings a list of servers and tracks their status.
- **Email Notifications**: Sends alert emails when a server goes down or comes back up.
- **Concurrency**: Supports concurrent monitoring of multiple servers using a worker pool.
- **Rate Limiting**: Controls the rate of ping requests to avoid overwhelming the network.
- **Graceful Shutdown**: Handles system signals for graceful shutdown of the monitoring service.
- **Customizable Configuration**: Easily configure servers, email settings, and other parameters via a JSON configuration file.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Setting Up as a Systemd Service (Linux)](#setting-up-as-a-systemd-service)
- [Usage](#usage)
- [Logging](#logging)
- [Health Check Endpoint](#health-check-endpoint)
- [License](#license)

## Installation

### Prerequisites

- Go 1.15+ installed on your machine.
- An SMTP server for sending email alerts.

### Steps

1. **Clone the Repository:**
   ```sh
   https://github.com/pentesterest-volkan/Moniping.git
   cd Moniping
   ```

2. **Install Dependencies:**
   ```sh
   go mod tidy
   ```

3. **Build the Project:**
   ```sh
   go build -o moniping
   ```

4. **Run the Application:**
   ```sh
   ./moniping -config=config.json
   ```


### Configuration

- The application is configured via a JSON file. Below is an example configuration:

  ```json
  {
  "ClientName": "Client Name",
  "Servers": [
  {
  "IP": "Server IP",
  "Name": "Server Name",
  "MaxRetries": 3,
  "Recipients": ["client@client.com"],
  "PingTimeout": 5,
  "PingCount": 3
  }
  ],
  "Email": {
  "From": "monitoring@server.com",
  "Password": "smtp-server-password",
  "SMTPHost": "smtp.server.com",
  "SMTPPort": "587",
  "DefaultRecipients": ["admin@admin.com"],
  "InsecureSkipVerify": true
  },
  "WorkerCount": 5,
  "RateLimit": 10
  }
  ```

### Setting Up as a Systemd Service

Systemd is a system and service manager for Linux operating systems. To set up your Go application as a systemd service, follow these steps:

1. **Create a Systemd Service File:**
Create a service file in the /etc/systemd/system/ directory. For example, create a file named server-monitor.service:
   ```sh
   sudo nano /etc/systemd/system/moniping.service
   ```
2. **Define the Service Configuration:**

   ```ini
   [Unit]
   Description=Go-Based Server Monitoring Service
   After=network.target
   
   [Service]
   ExecStart=/path/to/your/server-monitor -config=/path/to/your/config.json
   WorkingDirectory=/path/to/your
   Restart=always
   User=your-username
   Group=your-group
   Environment=GODEBUG=netdns=go
   
   [Install]
   WantedBy=multi-user.target
   ```

3. **Reload Systemd and Start the Service:**

   ```sh
   sudo systemctl daemon-reload
   ```
   ```sh
   sudo systemctl start moniping
   ```
   ```sh
   sudo systemctl enable moniping
   ```

4. **Check the Service Status:**

You can check the status of your service with:
   ```sh
   sudo systemctl status moniping
   ```

### Usage
- To start the monitoring service, run the compiled binary with the -config flag pointing to your JSON configuration file:
   ```sh
   ./moniping -config=config.json
   ```
- -config: Path to the JSON configuration file.

### Logging

- The application uses logrus for logging. Logs are printed to the console and optionally saved to a file. The log file is automatically trimmed to keep only the last 1,000 lines.

### Health Check Endpoint
A basic health check endpoint is available at /health on port 8080. You can use this to verify that the service is running:
   ```sh
   curl http://localhost:8080/health
   ```

Expected response:
   ```sh
   HTTP/1.1 200 OK
   ```

### License
This project is licensed under the MIT License.

