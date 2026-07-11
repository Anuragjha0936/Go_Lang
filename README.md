# Website Health Checker CLI

A lightweight command-line application written in Go that checks whether a website or server is reachable by attempting a TCP connection to a specified domain and port.

## Features

- Check if a website/server is reachable
- Supports custom ports
- Default port is `80`
- Fast TCP connectivity check with configurable timeout
- Simple command-line interface
- Cross-platform (Windows, macOS, Linux)

---

## Tech Stack

- Go
- urfave/cli/v3
- Go net package

---

## Installation

### Clone the repository

```bash
git clone https://github.com/<your-username>/website-health-checker.git

cd website-health-checker
```

### Install dependencies

```bash
go mod tidy
```

### Run

```bash
go run main.go --domain google.com
```

or

```bash
go run main.go -d google.com
```

---

## Usage

### Check a website on the default HTTP port (80)

```bash
go run main.go --domain google.com
```

### Check a custom port

```bash
go run main.go --domain google.com --port 443
```

or

```bash
go run main.go -d google.com -p 443
```

---

## Example Output

### Website is reachable

```text
[UP] google.com is reachable
From: 192.168.1.10:54321
To: 142.250.183.206:80
```

### Website is unreachable

```text
[DOWN] example.com is unreachable
Error: dial tcp: i/o timeout
```

---

## Project Structure

```
.
├── main.go
├── go.mod
├── go.sum
└── README.md
```

---

## How It Works

1. Accepts a domain name from the command line.
2. Uses port `80` if no port is provided.
3. Attempts a TCP connection using `net.DialTimeout`.
4. Waits up to **5 seconds** for a response.
5. Reports whether the destination is reachable.

---

## Future Improvements

- HTTPS health checks
- HTTP status code validation
- Ping latency measurement
- Check multiple domains from a file
- Export results as JSON or CSV
- Continuous monitoring mode
- DNS lookup information
- Colored terminal output

---

## Author

**Anurag Jha**

- GitHub: https://github.com/<your-username>
- LinkedIn: https://linkedin.com/in/<your-linkedin>

---

## License

This project is licensed under the MIT License.
