package transfer

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sanyamgarg/airpipe/internal/crypto"
)

type Receiver struct {
	relayURL string
	token    string
	key      []byte
	conn     *websocket.Conn
}

func NewReceiver(relayURL, token string, key []byte) *Receiver {
	return &Receiver{relayURL: relayURL, token: token, key: key}
}

func (r *Receiver) Connect() error {
	url := fmt.Sprintf("%s/ws/%s", r.relayURL, r.token)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to relay: %w", err)
	}
	r.conn = conn

	readyMsg := NewReadyMessage()
	encryptedReady, err := crypto.EncryptChunk(EncodeMessage(readyMsg), r.key)
	if err != nil {
		return fmt.Errorf("failed to encrypt ready message: %w", err)
	}

	if err := r.conn.WriteMessage(websocket.BinaryMessage, encryptedReady); err != nil {
		return fmt.Errorf("failed to send ready message: %w", err)
	}

	return nil
}

func (r *Receiver) ReceiveFile(destDir string, progressFn func(received, total int64)) (string, error) {
	var metadata Metadata
	var file *os.File
	var bytesReceived int64
	var destPath string

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	for {
		r.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		_, encryptedData, err := r.conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read message: %w", err)
		}

		decrypted, err := crypto.DecryptChunk(encryptedData, r.key)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt message: %w", err)
		}

		msg, err := DecodeMessage(decrypted)
		if err != nil {
			return "", fmt.Errorf("failed to decode message: %w", err)
		}

		switch msg.Type {
		case MsgTypeMetadata:
			metadata, err = ParseMetadata(msg.Payload)
			if err != nil {
				return "", fmt.Errorf("failed to parse metadata: %w", err)
			}
			destPath = filepath.Join(destDir, metadata.Filename)
			file, err = os.Create(destPath)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}

		case MsgTypeChunk:
			if file == nil {
				return "", fmt.Errorf("received chunk before metadata")
			}
			n, err := file.Write(msg.Payload)
			if err != nil {
				return "", fmt.Errorf("failed to write chunk: %w", err)
			}
			bytesReceived += int64(n)
			if progressFn != nil {
				progressFn(bytesReceived, metadata.Size)
			}

		case MsgTypeComplete:
			return destPath, nil

		case MsgTypeError:
			return "", fmt.Errorf("sender error: %s", string(msg.Payload))
		}
	}
}

func (r *Receiver) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
