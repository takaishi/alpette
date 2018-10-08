package stns

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/STNS/STNS/model"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
)

type stnsTC struct {
	info *credentials.ProtocolInfo
	stnsAddress string
	stnsPort string
}

const rs3Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (tc *stnsTC) randString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = rs3Letters[int(rand.Int63()%int64(len(rs3Letters)))]
	}
	return string(b)
}

func getUsername() (string, error) {
	if os.Getenv("SSH_USER") == "" {
		return "", errors.New("require SSH_USER")
	}
	return os.Getenv("SSH_USER"), nil
}

func privateKeyPath() string {
	if os.Getenv("SSH_PRIVATE_KEY_PATH") != "" {
		return os.Getenv("SSH_PRIVATE_KEY_PATH")
	} else {
		return fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
	}
}

func publicKeyPath() string {
	if os.Getenv("SSH_PUBLIC_KEY_PATH") != "" {
		return os.Getenv("SSH_PUBLIC_KEY_PATH")
	} else {
		return fmt.Sprintf("%s/.ssh/id_rsa.pub", os.Getenv("HOME"))
	}
}

func (tc *stnsTC) ClientHandshake(ctx context.Context, addr string, rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	username, err := getUsername()
	log.Printf("[DEBUG] username: %s\n", string(username))

	if err != nil {
		return nil, nil, err
	}
	rawConn.Write([]byte(username))

	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		log.Printf("[ERROR] Read error: %s\n", err)
		return nil, nil, err
	}
	log.Printf("[DEBUG] privateKeyPath: %s\n", privateKeyPath())
	log.Printf("[DEBUG] buf: %s\n", string(buf[:n]))
	key, err := tc.readPrivateKey(privateKeyPath())
	if err != nil {
		log.Printf("[ERROR] Failed to read private key: %s\n", err)
		return nil, nil, err
	}

	decrypted, err := tc.Decrypt(string(buf[:n]), key)
	if err != nil {
		log.Printf("[ERROR] Failed to decrypt: %s\n", err)
		return nil, nil, err
	}
	h := sha256.Sum256([]byte(decrypted))

	rawConn.Write([]byte(fmt.Sprintf("%x\n", h)))

	r := make([]byte, 64)
	n, err = rawConn.Read(r)
	if err != nil {
		log.Printf("[ERROR] Read error: %s\n", err)
		return nil, nil, err
	}
	r = r[:n]
	if string(r) != "ok" {
		log.Println("[ERROR] Failed to authenticate")
		return nil, nil, errors.New("Failed to authenticate")
	} else {
		log.Println("[INFO] success to authenticate")
	}

	return rawConn, nil, err
}

func (tc *stnsTC) getPubKeyFromSTNS(name string) ([]byte, error) {
	var user_resp []*model.User

	resp, err := http.Get(fmt.Sprintf("http://localhost:1104/v1/users?name=%s", name))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &user_resp)
	return []byte(user_resp[0].Keys[0]), nil
}

func (tc *stnsTC) ServerHandshake(rawConn net.Conn) (_ net.Conn, _ credentials.AuthInfo, err error) {
	// 乱数を生成する
	s := tc.randString()
	h := sha256.Sum256([]byte(s))

	// ユーザー名を読み込み&STNSからPublicKeyを取得
	buf := make([]byte, 2014)
	n, err := rawConn.Read(buf)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Read error: %s\n", err))
	}
	rawKey, err := tc.getPubKeyFromSTNS(string(buf[:n]))
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] rawKey = %s\n", string(rawKey))

	pubKey, err := tc.ParsePublicKey(rawKey)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to parse: %s\n", err))
	}
	encrypted, err := tc.Encrypt(s, pubKey)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to encrypt: %s\n", err))
	}
	rawConn.Write([]byte(encrypted))

	buf = make([]byte, 2014)
	n, err = rawConn.Read(buf)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Read error: %s\n", err))
	}
	buf = buf[:n]
	if strings.TrimRight(string(buf), "\n") == fmt.Sprintf("%x", h) {
		rawConn.Write([]byte("ok"))
		log.Println("[INFO] Success!!!")
	} else {
		rawConn.Write([]byte("ng"))
		log.Println("[INFO] Failed!!!")
		return nil, nil, errors.New(fmt.Sprintf("Failed to authenticate: invalid key"))
	}

	return rawConn, nil, err
}

func (tc *stnsTC) Info() credentials.ProtocolInfo {
	return *tc.info
}

func (tc *stnsTC) Clone() credentials.TransportCredentials {
	info := *tc.info
	return &stnsTC{
		info: &info,
	}
}

func (tc *stnsTC) OverrideServerName(serverNameOverride string) error {
	return nil
}

func NewServerCreds(stnsAddress string, stnsPort string) credentials.TransportCredentials {
	return &stnsTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion:  "1.0",
		},
		stnsAddress: stnsAddress,
		stnsPort: stnsPort,
	}
}

func NewClientCreds() credentials.TransportCredentials {
	return &stnsTC{
		info: &credentials.ProtocolInfo{
			SecurityProtocol: "ssh",
			SecurityVersion:  "1.0",
		},
	}
}
