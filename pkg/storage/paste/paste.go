package paste

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// TTL time to live for pastes
type TTL string

// possible values for PasteTTL
const (
	OneTime              TTL = "1t"
	OneDay               TTL = "1d"
	TwoDay               TTL = "2d"
	MinRecoveryKeyLength int = 6
	MaxRecoveryKeyLength int = 12
)

const DefaultStoragePath string = "/var/sll/pastes"

type Paste interface {
	Save(content string) error
	Retrieve() error
}

type paste struct {
	TTL              TTL
	Filename         string
	RecoveryKey      string
	EncryptionKey    []byte
	Content          string
	EncryptedContent []byte
	StoragePath      string
}

type Error struct {
	Code    int
	Message string
}

// Implement error interface
func (e Error) Error() string {
	return e.Message
}

func (p *paste) populateRecoveryKey() error {
	// do nothing is recovery key has already been generated
	if p.RecoveryKey != "" {
		return nil
	}

	// variable recovery key length
	variableLength := MinRecoveryKeyLength
	randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(MaxRecoveryKeyLength-MinRecoveryKeyLength)))
	if err != nil {
		return err
	}
	variableLength = variableLength + int(randomInt.Int64())

	// recovery key generation
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, variableLength)
	for i := 0; i < variableLength; i++ {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return err
		}
		ret[i] = letters[randomInt.Int64()]
	}

	p.RecoveryKey = string(ret)

	return nil
}

func (p *paste) populateEncryptionKey() error {
	if p.RecoveryKey == "" {
		return NewError(400, "Error generating encryption key: Invalid recovery key")
	}
	encryptionKey := sha256.Sum256([]byte(p.RecoveryKey))
	// convert array to slice
	p.EncryptionKey = encryptionKey[:]

	return nil
}

func (p *paste) populateKeys() error {
	err := p.populateRecoveryKey()
	if err != nil {
		return err
	}

	err = p.populateEncryptionKey()
	if err != nil {
		return err
	}

	return nil
}

func (p paste) generateFilenamePrefix() (string, error) {
	if p.RecoveryKey == "" {
		return "", NewError(400, "Error generating filename: Invalid recovery key")
	}

	data := []byte(p.RecoveryKey)
	h := md5.Sum(data)

	hash := hex.EncodeToString(h[:])

	return string(hash), nil
}

func (p *paste) populateFilename() error {
	if p.Filename == "" {
		prefix, err := p.generateFilenamePrefix()
		if err != nil {
			return err
		}

		p.Filename = prefix + "-" + string(p.TTL) + ".paste"
	}

	return nil
}

func (p *paste) populateTTL() {
	data := strings.Split(p.Filename, ".")
	data = strings.Split(data[0], "-")

	p.TTL = TTL(data[1])
}

func (p *paste) Encrypt() error {
	// https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/

	text := []byte(p.Content)
	// convert array to slice

	// generate a new aes cipher using our key
	c, err := aes.NewCipher(p.EncryptionKey)
	if err != nil {
		return err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	p.EncryptedContent = gcm.Seal(nonce, nonce, text, nil)

	return nil
}

func (p *paste) Decrypt() error {
	c, err := aes.NewCipher(p.EncryptionKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(p.EncryptedContent) < nonceSize {
		return err
	}

	nonce, ciphertext := p.EncryptedContent[:nonceSize], p.EncryptedContent[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	p.Content = string(plaintext)

	return nil
}

func (p *paste) GetInfo() error {
	ps := string(os.PathSeparator)

	filenamePrefix, err := p.generateFilenamePrefix()
	if err != nil {
		return err
	}

	matches, _ := filepath.Glob(DefaultStoragePath + ps + filenamePrefix + "-*.paste")
	for _, match := range matches {
		p.Filename = filepath.Base(match)
	}

	if p.Filename == "" {
		return NewError(404, "Paste not found")
	}

	p.populateTTL()

	return nil
}

func (p *paste) Save(content string) error {
	ps := string(os.PathSeparator)
	p.Content = content

	err := p.populateKeys()
	if err != nil {
		return err
	}

	if err := p.Encrypt(); err != nil {
		return NewError(500, "Could not encrypt paste")
	}

	if err := p.populateFilename(); err != nil {
		return err
	}

	// the WriteFile method returns an error if unsuccessful
	if err := ioutil.WriteFile(DefaultStoragePath+ps+p.Filename, p.EncryptedContent, 0664); err != nil {
		return NewError(500, "Could not save paste")
	}

	return nil
}

func (p *paste) Retrieve() error {
	ps := string(os.PathSeparator)

	if err := p.GetInfo(); err != nil {
		return err
	}

	pastePath := DefaultStoragePath + ps + p.Filename

	ciphertext, err := ioutil.ReadFile(pastePath)
	if err != nil {
		return NewError(404, "Paste not found")
	}

	p.EncryptedContent = ciphertext

	if err := p.Decrypt(); err != nil {
		return NewError(500, "Error decrypting paste "+p.RecoveryKey)
	}

	if p.TTL == OneTime {
		if err := os.Remove(pastePath); err != nil {
			return NewError(500, "Error destroying one time paste "+p.RecoveryKey)
		}
	}

	return nil
}

func New(ttl TTL, storagePath ...string) paste {
	sp := DefaultStoragePath
	if len(storagePath) > 0 {
		sp = storagePath[0]
	}

	p := paste{TTL: ttl, StoragePath: sp}

	return p
}

func Load(recoveryKey string, storagePath ...string) paste {
	sp := DefaultStoragePath
	if len(storagePath) > 0 {
		sp = storagePath[0]
	}

	p := paste{RecoveryKey: recoveryKey, StoragePath: sp}
	p.populateEncryptionKey()

	return p
}

func NewError(code int, message string) *Error {
	e := &Error{Code: code, Message: message}
	return e
}
