package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//TODO run memory diagnostic: https://golangcode.com/print-the-current-memory-usage/

func TestSetup(t *testing.T) {
	pubKey, masterKey, err := Setup()
	require.NoError(t, err)
	assert.NotZero(t, pubKey)
	assert.NotZero(t, masterKey)
}

func TestABE(t *testing.T) {
	pubKey, masterKey, err := Setup()
	require.NoError(t, err)

	privKey, err := GenerateKey(pubKey, masterKey, []string{"admin", "it_departement"})
	require.NoError(t, err)
	assert.NotZero(t, privKey)

	message := []byte("some secret message")
	cphData, aesData, err := Encrypt(pubKey, "(admin and it_departement)", message)
	require.NoError(t, err)

	decryptedMessage, err := Decrypt(pubKey, privKey, cphData, aesData)
	require.NoError(t, err)
	assert.Equal(t, message, decryptedMessage)
}

func TestABEKeyReuse(t *testing.T) {
	pubKey, masterKey, err := Setup()
	require.NoError(t, err)

	privKey, err := GenerateKey(pubKey, masterKey, []string{"admin", "it_departement"})
	require.NoError(t, err)
	assert.NotZero(t, privKey)

	encryptedKey, symKey, err := EncryptKey(pubKey, "(admin and it_departement)")
	require.NoError(t, err)
	require.NotNil(t, encryptedKey)
	require.NotNil(t, symKey)

	message := []byte("some secret message")
	cipher, err := EncryptWithKey(message, symKey)
	require.NoError(t, err)
	require.NotNil(t, cipher)

	//TODO this does not work atm...most probably element serialization gives the wrong result
	//it overflows heavily
	decryptedMessage, err := Decrypt(pubKey, privKey, encryptedKey, cipher)
	require.NoError(t, err)
	assert.Equal(t, message, decryptedMessage)
}

func TestABEForNotMatchingPolicy(t *testing.T) {
	pubKey, masterKey, err := Setup()
	require.NoError(t, err)

	privKey, err := GenerateKey(pubKey, masterKey, []string{"admin", "it_departement"})
	require.NoError(t, err)
	assert.NotZero(t, privKey)

	message := []byte("some secret message")
	cphData, aesData, err := Encrypt(pubKey, "(ceo and it_departement)", message)
	require.NoError(t, err)

	decryptedMessage, err := Decrypt(pubKey, privKey, cphData, aesData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot decrypt, attributes in key do not satisfy policy")
	assert.Equal(t, []byte(nil), decryptedMessage)
}
