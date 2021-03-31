package paste

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTTL(t *testing.T) {
	a := assert.New(t)

	expiration := OneTime

	p := New(expiration)

	a.Equalf(p.TTL, expiration, "Expiration should be: %s", expiration)
}

func TestNewDefaultStoragePath(t *testing.T) {
	a := assert.New(t)

	expiration := OneTime

	p := New(expiration)

	a.Equalf(p.StoragePath, DefaultStoragePath, "Storage path should be the default one: %s", DefaultStoragePath)
}

func TestNewWithStoragePath(t *testing.T) {
	a := assert.New(t)

	expiration := OneTime
	storagePath := "/tmp/foo"

	p := New(expiration, storagePath)

	a.Equalf(p.StoragePath, storagePath, "Storage path should be: %s", storagePath)
}

func TestPopulateRecoveryKeyNotEmpty(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	p.populateRecoveryKey()

	a.NotEmpty(p.RecoveryKey)
}

func TestPopulateRecoveryKeyLength(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	p.populateRecoveryKey()

	a.LessOrEqualf(len(p.RecoveryKey), MaxRecoveryKeyLength, "Recovery key length cannot be greater than: %s", MaxRecoveryKeyLength)
	a.GreaterOrEqualf(len(p.RecoveryKey), MinRecoveryKeyLength, "Recovery key length cannot be lesser than: %s", MinRecoveryKeyLength)
}

func TestPopulateRecoveryKeyRandomness(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)
	p2 := New(OneTime)

	p.populateRecoveryKey()
	p2.populateRecoveryKey()

	a.NotEqualf(p.RecoveryKey, p2.RecoveryKey, "Generated recovery keys should be different: %s != %s", p.RecoveryKey, p2.RecoveryKey)
}

func TestPopulateEncryptionKey(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	p.RecoveryKey = "abcdef123"

	p.populateEncryptionKey()

	a.NotEmpty(p.EncryptionKey, "Encryption key is empty")
}

func TestPopulateEncryptionKeyLength(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	p.RecoveryKey = "abcdef123"

	p.populateEncryptionKey()

	a.Equal(len(p.EncryptionKey), 32, "Encryption key length must be 32 bytes")
}

func TestPopulateEncryptionKeyIsSHA256(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)
	p.RecoveryKey = "abcdef123"
	expectedEncryptionKey := []uint8{159, 76, 18, 29, 96, 207, 85, 58, 216, 225, 183, 63, 108, 2, 173, 38, 137, 180, 84, 7, 87, 18, 209, 102, 95, 114, 104, 91, 201, 48, 68, 194}

	p.populateEncryptionKey()

	a.Equal(p.EncryptionKey, expectedEncryptionKey, "Encryption key is not a SHA256 of recovery key")
}

func TestPopulateEncryptionKeyNotExistingRecoveryKey(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	a.Error(p.populateEncryptionKey())
	a.EqualError(p.populateEncryptionKey(), "Error generating encryption key: Invalid recovery key")
}

func TestLoadRecoveryKey(t *testing.T) {
	a := assert.New(t)

	recoveryKey := "abcdef123"
	p := Load(recoveryKey)

	a.Equalf(p.RecoveryKey, recoveryKey, "Recovery key should be: %s", recoveryKey)
}

func TestLoadEncryptionKeyPresent(t *testing.T) {
	a := assert.New(t)

	recoveryKey := "abcdef123"
	p := Load(recoveryKey)

	a.NotEmpty(p.EncryptionKey, "Encryption key has not been populated")
}

func TestLoadDefaultStoragePath(t *testing.T) {
	a := assert.New(t)

	recoveryKey := "abcdef123"

	p := Load(recoveryKey)

	a.Equalf(p.StoragePath, DefaultStoragePath, "Storage path should be the default one: %s", DefaultStoragePath)
}

func TestLoadWithStoragePath(t *testing.T) {
	a := assert.New(t)
	recoveryKey := "abcdef123"
	storagePath := "/tmp/foo"

	p := Load(recoveryKey, storagePath)

	a.Equalf(p.StoragePath, storagePath, "Storage path should be: %s", storagePath)
}

func TestGenerateFilenamePrefix(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)
	p.RecoveryKey = "abcdef123"

	prefix, err := p.generateFilenamePrefix()

	a.NotEmpty(prefix)
	a.Nil(err)
}

func TestGenerateFilenamePrefixNotExistingRecoveryKey(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	_, err := p.generateFilenamePrefix()

	a.Error(err)
}

func TestPopulateFilenameNotEmpty(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)
	p.RecoveryKey = "abcdef123"

	p.populateFilename()

	a.NotEmpty(p.Filename, "Filename should not be empty")
}

func TestPopulateFilenameIfAlreadyGenerated(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)
	p.RecoveryKey = "abcdef123"
	filename := "abcdef123.1t"
	p.Filename = filename

	p.populateFilename()

	a.Equal(p.Filename, filename, "A new filename has been generated and it already was present")
}

func TestPopulateFilenameNotExistingRecoveryKey(t *testing.T) {
	a := assert.New(t)

	p := New(OneTime)

	a.Error(p.populateFilename())
}

func TestEncrypt(t *testing.T) {
	a := assert.New(t)

	expectedEncryptedContentLength := 37
	p := New(OneTime)
	p.RecoveryKey = "abcdef123"
	p.EncryptionKey = []uint8{159, 76, 18, 29, 96, 207, 85, 58, 216, 225, 183, 63, 108, 2, 173, 38, 137, 180, 84, 7, 87, 18, 209, 102, 95, 114, 104, 91, 201, 48, 68, 194}
	p.Content = "foobarbaz"

	p.Encrypt()

	a.Equal(len(p.EncryptedContent), expectedEncryptedContentLength, "Encrypted content length different from expected: %s", expectedEncryptedContentLength)
}
