# AirPipe

Transfer files between terminal and phone with a QR code. No apps. No accounts. End-to-end encrypted.
```
$ airpipe send config.yaml
```

![demo](demo.gif)

Scan → file downloads. Done.

## Install
```bash
# macOS
brew install sanyam-g/tap/airpipe

# Linux (amd64)
curl -sL https://github.com/Sanyam-G/Airpipe/releases/latest/download/airpipe-linux-amd64 -o airpipe
chmod +x airpipe
sudo mv airpipe /usr/local/bin/

# From source
go install github.com/Sanyam-G/Airpipe/cmd/airpipe@latest
```

## Usage

**Send file (server → phone):**
```bash
airpipe send ./error.log
```

**Receive file (phone → server):**
```bash
airpipe receive ./downloads
```

## How it works

1. CLI generates encryption key locally
2. Key embedded in URL fragment (`#...`) — never sent to server
3. File encrypted locally, streamed through relay
4. Phone decrypts in browser
5. Relay only sees encrypted bytes — zero knowledge

## Self-host relay
```bash
docker run -p 8080:8080 ghcr.io/sanyam-g/airpipe-relay
airpipe --relay wss://your-server:8080 send file.txt
```

## Security

- **Encryption:** NaCl secretbox (XSalsa20-Poly1305)
- **Key exchange:** None — key in URL fragment, never touches server
- **Tokens:** Random, single-use, 10-minute expiry

## License

MIT
