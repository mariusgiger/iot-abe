package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/pkg/errors"
)

// VerifySignature checks the signature of a transaction
func VerifySignature(tx *types.Transaction, networkID *big.Int) (bool, error) {
	signer := types.NewEIP155Signer(networkID)
	hash := signer.Hash(tx)

	signature, err := extractSignature(tx, networkID)
	if err != nil {
		return false, err
	}

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false, err
	}

	signatureNoRecoverID := signature[:len(signature)-1]
	verified := crypto.VerifySignature(sigPublicKey, hash.Bytes(), signatureNoRecoverID)

	return verified, nil
}

// RecoverPubKey recovers the public key of a transaction
func RecoverPubKey(tx *types.Transaction, networkID *big.Int) (*ecdsa.PublicKey, error) {
	signature, err := extractSignature(tx, networkID)
	if err != nil {
		return nil, err
	}

	signer := types.NewEIP155Signer(networkID)
	hash := signer.Hash(tx)

	return crypto.SigToPub(hash.Bytes(), signature)
}

// SenderFromTx retrieves the sender from a transaction
func SenderFromTx(tx *types.Transaction, networkID *big.Int) (common.Address, error) {
	signer := types.NewEIP155Signer(networkID)
	return signer.Sender(tx)
}

// GenerateKeyPair generates a new ecdsa key pair
func GenerateKeyPair() (*ecdsa.PrivateKey, common.Address, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, common.Address{}, err
	}

	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("error casting public key to ECDSA")
	}

	//cutting off 4 bytes of the hex prefix '0x' and EC uncompressed point prefix '04'
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return key, address, nil
}

// EncryptMessage asymmetrically based on ecies
func EncryptMessage(pubKey *ecdsa.PublicKey, message []byte) ([]byte, error) {
	//TODO should a shared key be generated
	//TODO further research about ecies
	//TODO this might be further improved if the requester adds a different publicKey to the contract

	// prvKey, err := ecies.GenerateKey(rand.Reader, ethCrypto.S256(), nil)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not generate ecies privateKey")
	// }

	//
	// skLen := ecies.MaxSharedKeyLength(eciesPub) / 2

	// sharedKey, err := prvKey.GenerateShared(eciesPub, skLen, skLen)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not generate ecies sharedKey")
	// }

	eciesPub := ecies.ImportECDSAPublic(pubKey)
	encryptedKey, err := ecies.Encrypt(rand.Reader, eciesPub, message, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not encrypt message")
	}

	return encryptedKey, nil
}

// DecryptMessage decrypts a message based on ecies
func DecryptMessage(prv *ecdsa.PrivateKey, ct []byte) ([]byte, error) {
	fmt.Println("got here")
	eciesPrv := ecies.ImportECDSA(prv)
	message, err := eciesPrv.Decrypt(ct, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not decrypt message")
	}

	return message, nil
}

func extractSignature(tx *types.Transaction, networkID *big.Int) ([]byte, error) {
	var signature []byte
	vTx, r, s := tx.RawSignatureValues()
	v := big.NewInt(vTx.Int64())
	//Ethereum magic, refer to: https://github.com/ethereum/go-ethereum/blob/4c181e4fb98bb88503cccd6147026b6c2b7b56f6/core/types/transaction_signing.go#L195
	chainIDMul := new(big.Int).Mul(networkID, big.NewInt(2))
	if networkID.Sign() != 0 {
		v = v.Sub(v, big.NewInt(35))
		v = v.Sub(v, chainIDMul)
	} else {
		v = v.Sub(v, big.NewInt(27))
	}

	signature = append(signature, r.Bytes()...)
	signature = append(signature, s.Bytes()...)
	if v.Cmp(big.NewInt(0)) == 0 {
		signature = append(signature, byte(0))
	} else {
		signature = append(signature, v.Bytes()...)
	}

	if len(signature) != 65 {
		return nil, errors.New("signature length invalid")
	}

	return signature, nil
}
