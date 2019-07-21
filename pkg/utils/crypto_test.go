package utils

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mariusgiger/iot-abe/pkg/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/suite"
)

// CryptoUtilsTestSuite contains crypto utils related tests
type CryptoUtilsTestSuite struct {
	suite.Suite
}

func (ts *CryptoUtilsTestSuite) TestEncryptDecrypt() {
	//arrange
	message := []byte("Some secret message.")

	prv, _, err := GenerateKeyPair()
	ts.Require().NoError(err)

	//act
	ct, err := EncryptMessage(&prv.PublicKey, message)
	ts.Require().NoError(err)

	decryptedMessage, err := DecryptMessage(prv, ct)
	ts.Require().NoError(err)

	//assert
	ts.Equal(message, decryptedMessage)
}

func (ts *CryptoUtilsTestSuite) TestVerifySignature() {
	//arrange
	privKey, _, err := GenerateKeyPair()
	ts.Require().NoError(err)
	chainID := big.NewInt(3)
	signedTx := ts.generateSignedTx(privKey, chainID)

	//act
	verified, err := VerifySignature(signedTx, chainID)

	//assert
	ts.Require().NoError(err)
	ts.True(verified)
}

func (ts *CryptoUtilsTestSuite) TestRecoverPubKey() {
	//arrange
	privKey, expectedAddr, err := GenerateKeyPair()
	ts.Require().NoError(err)
	expectedPubKeyECDSA, ok := privKey.Public().(*ecdsa.PublicKey)
	ts.Require().True(ok)
	expectedPubKeyBytes := crypto.FromECDSAPub(expectedPubKeyECDSA)

	chainID := big.NewInt(3)
	signedTx := ts.generateSignedTx(privKey, chainID)
	ts.Require().Equal(chainID, signedTx.ChainId())

	//act
	pubKey, err := RecoverPubKey(signedTx, chainID)

	//assert
	ts.Require().NoError(err)
	publicKeyBytes := crypto.FromECDSAPub(pubKey)
	addr := crypto.PubkeyToAddress(*pubKey)
	ts.Equal(expectedAddr, addr)
	ts.Equal(expectedPubKeyBytes, publicKeyBytes)
}

func (ts *CryptoUtilsTestSuite) generateSignedTx(privKey *ecdsa.PrivateKey, chainID *big.Int) *types.Transaction {
	amount := big.NewInt(100000000000000000)
	nonce := uint64(1)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(core.GWei * 10)

	to := common.HexToAddress("0xa9a0E7C567f5fE4f9C7f684b3398FD74041385BF")

	signer := types.NewEIP155Signer(chainID)
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privKey)
	ts.Require().NoError(err)

	signedTx, err := tx.WithSignature(signer, signature)
	ts.Require().NoError(err)

	return signedTx
}

func TestCryptoUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoUtilsTestSuite))
}
