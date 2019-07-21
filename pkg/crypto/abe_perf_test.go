package crypto

import (
	"fmt"
	"testing"

	"github.com/mariusgiger/iot-abe/pkg/utils"
)

func BenchmarkSetup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if _, _, err := Setup(); err != nil {
			fmt.Println(err.Error())
			b.FailNow()
		}
	}
}

func TestGenerateSKLength(t *testing.T) {
	pubKey, masterKey, err := Setup()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}

	MaxAttrs := 50
	AttrLength := 10
	attrs := []string{}

	fmt.Println("attr,keylen")
	for k := 0; k < MaxAttrs; k++ {
		attr := utils.RandString(AttrLength)
		attrs = append(attrs, attr)

		key, err := GenerateKey(pubKey, masterKey, attrs)
		if err != nil {
			fmt.Println(err.Error())
			t.FailNow()
		}

		fmt.Printf("%v, %v\n", len(attrs), len(key))
	}
}

func BenchmarkGenerateSK(b *testing.B) {
	pubKey, masterKey, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	MaxAttrs := 50
	AttrLength := 10
	attrs := []string{}

	for k := 1; k <= MaxAttrs; k++ {
		attr := utils.RandString(AttrLength)
		attrs = append(attrs, attr)

		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if key, err := GenerateKey(pubKey, masterKey, attrs); err != nil {
					fmt.Println(err.Error())
					_ = key
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkEncryptKeyBytesWithIncreasingPolicyLength(b *testing.B) {
	pubKey, _, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	MaxAttrs := 50
	AttrLength := 10
	NumBytes := 1 //* 1024 * 1024 //1MB
	msgStr := []byte(utils.RandString(NumBytes))
	policy := "("

	for k := 1; k <= MaxAttrs; k++ {
		attr := utils.RandString(AttrLength)
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " and "
		}

		policy = policy + attr + ")"

		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if encKey, encMsg, err := Encrypt(pubKey, policy, msgStr); err != nil {
					fmt.Println(err.Error())
					_ = encKey
					_ = encMsg
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkEncryptKeyWithIncreasingPolicyLengthNoAES(b *testing.B) {
	pubKey, _, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	MaxAttrs := 50
	AttrLength := 10
	policy := "("

	for k := 1; k <= MaxAttrs; k++ {
		attr := utils.RandString(AttrLength)
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " and "
		}

		policy = policy + attr + ")"

		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if encKey, key, err := EncryptKey(pubKey, policy); err != nil {
					fmt.Println(err.Error())
					_ = encKey
					_ = key
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkDecryptKeyWithIncreasingPolicyLength(b *testing.B) {
	pubKey, masterKey, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	MaxAttrs := 50
	AttrLength := 10
	NumBytes := 1 //* 1024 * 1024 //1MB
	msgStr := []byte(utils.RandString(NumBytes))
	attrs := []string{}
	policy := "("

	for k := 1; k <= MaxAttrs; k++ {
		attr := utils.RandString(AttrLength)
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " and "
		}
		attrs = append(attrs, attr)
		policy = policy + attr + ")"

		sk, err := GenerateKey(pubKey, masterKey, attrs)
		if err != nil {
			fmt.Println(err.Error())
			b.FailNow()
		}

		encKey, _, err := Encrypt(pubKey, policy, msgStr)
		if err != nil {
			fmt.Println(err.Error())
			b.FailNow()
		}

		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if key, err := DecryptKey(pubKey, sk, encKey); err != nil {
					fmt.Println(err.Error())
					_ = key
					b.FailNow()
				}
			}
		})
	}
}

func TestEncryptKeyLength(t *testing.T) {
	pubKey, _, err := Setup()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}

	MaxAttrs := 50
	AttrLength := 10
	NumBytes := 1 * 1024 * 1024 //1MB
	msgStr := []byte(utils.RandString(NumBytes))
	policy := "("

	fmt.Println("attrs,enc_keylen,cph_len")
	for k := 1; k <= MaxAttrs; k++ {
		attr := utils.RandString(AttrLength)
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " and "
		}

		policy = policy + attr + ")"

		encKey, encMsg, err := Encrypt(pubKey, policy, msgStr)
		if err != nil {
			fmt.Println(err.Error())
			t.FailNow()
		}

		fmt.Printf("%v, %v, %v\n", k, len(encKey), len(encMsg))
	}
}

func BenchmarkEncryptDataLen(b *testing.B) {
	pubKey, _, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	NumBytes := 1 * 1024 * 1024 //1MB
	for k := 100; k < NumBytes; k = k + 1000 {
		msgStr := utils.RandString(k)
		msg := []byte(msgStr)

		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if _, _, err := Encrypt(pubKey, "(admin and it_departement)", msg); err != nil {
					fmt.Println(err.Error())
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkEncryptPolicyLenAnd(b *testing.B) {
	pubKey, _, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	NumPolicyAttrs := 20
	AttrLength := 20
	msgStr := utils.RandString(1 * 1024 * 1024) //1MB
	msg := []byte(msgStr)
	policy := "("
	for k := 1; k <= NumPolicyAttrs; k++ {
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " and "
		}

		policy = policy + utils.RandString(AttrLength) + ")"
		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if _, _, err := Encrypt(pubKey, policy, msg); err != nil {
					fmt.Println(err.Error())
					b.FailNow()
				}
			}
		})
	}
}

func BenchmarkEncryptPolicyLenOr(b *testing.B) {
	pubKey, _, err := Setup()
	if err != nil {
		b.Fatalf("could not generate key: %v", err)
	}

	NumPolicyAttrs := 20
	AttrLength := 20
	msgStr := utils.RandString(1 * 1024 * 1024) //1MB
	msg := []byte(msgStr)
	policy := "("
	for k := 1; k <= NumPolicyAttrs; k++ {
		if len(policy) > 1 {
			policy = policy[:len(policy)-1] + " or "
		}

		policy = policy + utils.RandString(AttrLength) + ")"
		b.Run(fmt.Sprintf("len_%v", k), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if _, _, err := Encrypt(pubKey, policy, msg); err != nil {
					fmt.Println(err.Error())
					b.FailNow()
				}
			}
		})
	}
}
