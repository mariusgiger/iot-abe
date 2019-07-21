//Package crypto wraps the libbswabe C library
//refer to: http://hms.isi.jhu.edu/acsc/cpabe/
package crypto

/*

#cgo CFLAGS: -I/usr/local/include/pbc
#cgo darwin CFLAGS: -I/usr/local/opt/openssl/include -I/usr/local/Cellar/glib/2.60.4/include/glib-2.0 -I/usr/local/Cellar/glib/2.60.4/lib/glib-2.0/include -I/usr/local/include/pbc -I/usr/local/opt/openssl/include
#cgo linux CFLAGS: -I/usr/include/glib-2.0 -I/usr/lib/x86_64-linux-gnu/glib-2.0/include -I/usr/lib/arm-linux-gnueabihf/glib-2.0/include/ -w
#cgo darwin LDFLAGS: -L/usr/local/Cellar/glib/2.60.4/lib -L/usr/local/opt/gettext/lib -lglib-2.0 -lintl -Wl,-rpath /usr/local/lib -lgmp -Wl,-rpath /usr/local/lib -lpbc -lbswabe -lcrypto -lcrypto
#cgo linux LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lglib-2.0 -lbswabe -lcrypto -Wl,-rpath /usr/local/lib -lpbc -lgmp
#cgo linux/arm LDFLAGS: -L/usr/lib/arm-linux-gnueabihf/glib-2.0/include/lib/ -lglib-2.0
#cgo linux LDFLAGS: -L/usr/local/lib

#include <glib.h>
#include <pbc.h>
#include <pbc_random.h>
#include "./libbswabe/bswabe.h"
#include "./libbswabe/core.c"
#include "./libpolicies/common.h"
#include "./libpolicies/policy_lang.h"

static char**makeCharArray(int size) {
        return calloc(sizeof(char*), size);
}

static void setArrayString(char **a, char *s, int n) {
        a[n] = s;
}

static void freeCharArray(char **a, int size) {
        int i;
        for (i = 0; i < size; i++)
                free(a[i]);
        free(a);
}

gint
comp_string(gconstpointer a, gconstpointer b)
{
	return strcmp(a, b);
}

static GSList* sort_glist(GSList* alist)
{
	return g_slist_sort(alist, comp_string);
}

static char** copyAttrs(GSList* alist)
{
	int i = 0;
	char** attrs    = 0;
	GSList* ap;

	int n = g_slist_length(alist);
	attrs = malloc((n + 1) * sizeof(char*));
	for(ap = alist; ap; ap = ap->next )
		attrs[i++] = ap->data;
	attrs[i] = 0;

	return attrs;
}

static struct element_s newArr()
{
	struct element_s m;
	return m;
}

static GByteArray* encrypt(GByteArray* pt, unsigned char* key)
{
	struct element_s k;
	struct pairing_s pairing;
	char* pairingBuf;

	pairingBuf = strdup(TYPE_A_PARAMS);
	pairing_init_set_buf(&pairing, pairingBuf, strlen(pairingBuf));
	element_init_GT(&k, &pairing);
	element_from_bytes(&k, key);

	GByteArray* ct = aes_128_cbc_encrypt(pt, &k);

	element_clear(&k);
	pairing_clear(&pairing);
	free(pairingBuf);

	return ct;
}
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

// helper for memcopying byte arras to char arrays
func memcpy(dest *C.uchar, src []byte) int {
	n := len(src)
	C.memcpy(unsafe.Pointer(dest), unsafe.Pointer(&src[0]), C.size_t(n))
	return n
}

// Setup sets up ABE
// returns the public key and the master key
func Setup() ([]byte, []byte, error) {
	var (
		pub *C.bswabe_pub_t
		msk *C.bswabe_msk_t
	)
	C.bswabe_setup(&pub, &msk)
	defer C.bswabe_pub_free(pub)
	defer C.bswabe_msk_free(msk)

	pubSerialized := C.bswabe_pub_serialize(pub)
	pubData := C.GoBytes(unsafe.Pointer(pubSerialized.data), C.int(pubSerialized.len))
	defer C.g_byte_array_free(pubSerialized, C.int(1))

	mskSerialized := C.bswabe_msk_serialize(msk)
	mskData := C.GoBytes(unsafe.Pointer(mskSerialized.data), C.int(mskSerialized.len))
	defer C.g_byte_array_free(mskSerialized, C.int(1))

	return pubData, mskData, nil
}

// GenerateKey generates a new private key with associated attributes
func GenerateKey(pubKey []byte, masterKey []byte, attrs []string) ([]byte, error) {
	var (
		pub *C.bswabe_pub_t
		msk *C.bswabe_msk_t
		prv *C.bswabe_prv_t

		pubArr *C.GByteArray
		mskArr *C.GByteArray

		alist *C.GSList
	)

	//parse public key
	pubArr = C.g_byte_array_new()
	C.g_byte_array_set_size(pubArr, C.uint(len(pubKey)))
	memcpy(pubArr.data, pubKey)
	pub = C.bswabe_pub_unserialize(pubArr, 1)
	defer C.bswabe_pub_free(pub)

	//parse master key
	mskArr = C.g_byte_array_new()
	C.g_byte_array_set_size(mskArr, C.uint(len(masterKey)))
	memcpy(mskArr.data, masterKey)
	msk = C.bswabe_msk_unserialize(pub, mskArr, 1)
	defer C.bswabe_msk_free(msk)

	//parse attributes
	for _, attr := range attrs {
		C.parse_attribute(&alist, C.CString(attr))
	}
	alist = C.sort_glist(alist)
	defer C.g_slist_free(alist)
	cAttrs := C.copyAttrs(alist) //TODO is this cleared in bswabe_prv_free?

	//generate private key
	prv = C.bswabe_keygen(pub, msk, cAttrs)
	defer C.bswabe_prv_free(prv)

	//serialize private key
	prvSerialized := C.bswabe_prv_serialize(prv)
	defer C.g_byte_array_free(prvSerialized, C.int(1))
	prvData := C.GoBytes(unsafe.Pointer(prvSerialized.data), C.int(prvSerialized.len))

	return prvData, nil
}

// Encrypt encrypts a messaged with provided policy
// returns the encrypted key (asymmetric), the encrypted message (symmetric)
func Encrypt(pubKey []byte, policy string, message []byte) ([]byte, []byte, error) {
	var (
		pub    *C.bswabe_pub_t
		pubArr *C.GByteArray
	)

	//parse policy
	cPolicy := C.parse_policy_lang(C.CString(policy))
	defer C.free(unsafe.Pointer(cPolicy))

	//parse public key
	pubArr = C.g_byte_array_new()
	C.g_byte_array_set_size(pubArr, C.uint(len(pubKey)))
	memcpy(pubArr.data, pubKey)
	pub = C.bswabe_pub_unserialize(pubArr, 1)
	defer C.bswabe_pub_free(pub)

	m := C.newArr()

	//encrypt key with policy
	cph := C.bswabe_enc(pub, &m, cPolicy)
	if cph == nil {
		err := C.bswabe_error()
		return nil, nil, fmt.Errorf("could not encrypt: %v", C.GoString(err))
	}
	defer C.element_clear(&m)
	defer C.bswabe_cph_free(cph)

	//serialize encrypted key
	cphBuf := C.bswabe_cph_serialize(cph)
	defer C.g_byte_array_free(cphBuf, 1)
	cphData := C.GoBytes(unsafe.Pointer(cphBuf.data), C.int(cphBuf.len))

	//convert plain text
	plt := C.g_byte_array_new()
	msgBytes := []byte(message)
	C.g_byte_array_set_size(plt, C.uint(len(msgBytes)))
	memcpy(plt.data, msgBytes)
	defer C.g_byte_array_free(plt, 1)

	//encrypt plain text
	aesBuf := C.aes_128_cbc_encrypt(plt, &m)
	defer C.g_byte_array_free(aesBuf, 1)
	aesData := C.GoBytes(unsafe.Pointer(aesBuf.data), C.int(aesBuf.len))

	return cphData, aesData, nil
}

// EncryptKey creates a new symmetric key and encrypts it with policy
// returns the encrypted key and the symmetric key
func EncryptKey(pubKey []byte, policy string) ([]byte, []byte, error) {
	var (
		pub    *C.bswabe_pub_t
		pubArr *C.GByteArray
	)

	//parse policy
	cPolicy := C.parse_policy_lang(C.CString(policy))
	defer C.free(unsafe.Pointer(cPolicy))

	//parse public key
	pubArr = C.g_byte_array_new()
	C.g_byte_array_set_size(pubArr, C.uint(len(pubKey)))
	memcpy(pubArr.data, pubKey)
	pub = C.bswabe_pub_unserialize(pubArr, 1)
	defer C.bswabe_pub_free(pub)

	m := C.newArr()
	defer C.element_clear(&m)

	//encrypt key with policy
	cph := C.bswabe_enc(pub, &m, cPolicy)
	if cph == nil {
		err := C.bswabe_error()
		return nil, nil, fmt.Errorf("could not encrypt: %v", C.GoString(err))
	}
	defer C.bswabe_cph_free(cph)

	//serialize encrypted key
	cphBuf := C.bswabe_cph_serialize(cph)
	defer C.g_byte_array_free(cphBuf, 1)
	cphData := C.GoBytes(unsafe.Pointer(cphBuf.data), C.int(cphBuf.len))

	keyArr := C.g_byte_array_new()
	keyLen := C.element_length_in_bytes(&m)
	C.g_byte_array_set_size(keyArr, C.uint(keyLen))
	C.element_to_bytes(keyArr.data, &m)
	keyArrBytes := C.GoBytes(unsafe.Pointer(keyArr.data), C.int(keyArr.len))

	return cphData, keyArrBytes, nil
}

// EncryptWithKey encrypts with provided symmetric key
// returns the encrypted message (symmetric)
func EncryptWithKey(message []byte, key []byte) ([]byte, error) {
	//convert plain text
	var plt = C.g_byte_array_new()
	C.g_byte_array_set_size(plt, C.uint(len(message)))
	memcpy(plt.data, message)
	defer C.g_byte_array_free(plt, 1)

	//encrypt plain text
	cKey := (*C.uchar)(unsafe.Pointer(&key))
	var aesBuf = C.encrypt(plt, cKey)
	defer C.g_byte_array_free(aesBuf, 1)

	aesData := C.GoBytes(unsafe.Pointer(aesBuf.data), C.int(aesBuf.len))
	fmt.Printf("%v\n", string(aesData))
	return aesData, nil
}

// Decrypt decrypts a cipher text
func Decrypt(pubKey []byte, prvKey []byte, cphData []byte, aesData []byte) ([]byte, error) {
	var (
		pub    *C.bswabe_pub_t
		pubArr *C.GByteArray

		prv    *C.bswabe_prv_t
		prvArr *C.GByteArray

		aesBuf *C.GByteArray
		plt    *C.GByteArray
		cphBuf *C.GByteArray
		cph    *C.bswabe_cph_t
	)

	//parse public key
	pubArr = C.g_byte_array_new()
	C.g_byte_array_set_size(pubArr, C.uint(len(pubKey)))
	memcpy(pubArr.data, pubKey)
	pub = C.bswabe_pub_unserialize(pubArr, 1)
	defer C.bswabe_pub_free(pub)

	//parse private key
	prvArr = C.g_byte_array_new()
	C.g_byte_array_set_size(prvArr, C.uint(len(prvKey)))
	memcpy(prvArr.data, prvKey)
	prv = C.bswabe_prv_unserialize(pub, prvArr, 1)
	defer C.bswabe_prv_free(prv)

	//parse encrypte key
	cphBuf = C.g_byte_array_new()
	C.g_byte_array_set_size(cphBuf, C.uint(len(cphData)))
	memcpy(cphBuf.data, cphData)
	cph = C.bswabe_cph_unserialize(pub, cphBuf, 1)
	defer C.bswabe_cph_free(cph)

	//parse aes cipher text
	aesBuf = C.g_byte_array_new()
	C.g_byte_array_set_size(aesBuf, C.uint(len(aesData)))
	memcpy(aesBuf.data, aesData)
	defer C.g_byte_array_free(aesBuf, C.int(1))

	m := C.newArr()
	//decrypt key
	res := C.bswabe_dec(pub, prv, cph, &m)
	if int(res) != 1 {
		err := C.bswabe_error()
		return nil, fmt.Errorf("could not decrypt: %v", C.GoString(err))
	}

	//decrypt cipher text
	//TODO this heavily overflows if decryption is not possible...
	plt = C.aes_128_cbc_decrypt(aesBuf, &m)

	//defer C.g_byte_array_free(plt, C.int(1))
	defer C.element_clear(&m)

	//convert clear text to go byte array
	msgData := C.GoBytes(unsafe.Pointer(plt.data), C.int(plt.len))

	return msgData, nil
}

// DecryptKey decrypts an abe encrypted key
func DecryptKey(pubKey []byte, prvKey []byte, cphData []byte) ([]byte, error) {
	var (
		pub    *C.bswabe_pub_t
		pubArr *C.GByteArray

		prv    *C.bswabe_prv_t
		prvArr *C.GByteArray

		cphBuf *C.GByteArray
		cph    *C.bswabe_cph_t
	)

	//parse public key
	pubArr = C.g_byte_array_new()
	C.g_byte_array_set_size(pubArr, C.uint(len(pubKey)))
	memcpy(pubArr.data, pubKey)
	pub = C.bswabe_pub_unserialize(pubArr, 1)
	defer C.bswabe_pub_free(pub)

	//parse private key
	prvArr = C.g_byte_array_new()
	C.g_byte_array_set_size(prvArr, C.uint(len(prvKey)))
	memcpy(prvArr.data, prvKey)
	prv = C.bswabe_prv_unserialize(pub, prvArr, 1)
	defer C.bswabe_prv_free(prv)

	//parse encrypte key
	cphBuf = C.g_byte_array_new()
	C.g_byte_array_set_size(cphBuf, C.uint(len(cphData)))
	memcpy(cphBuf.data, cphData)
	cph = C.bswabe_cph_unserialize(pub, cphBuf, 1)
	defer C.bswabe_cph_free(cph)

	m := C.newArr()
	//decrypt key
	res := C.bswabe_dec(pub, prv, cph, &m)
	if int(res) != 1 {
		err := C.bswabe_error()
		return nil, fmt.Errorf("could not decrypt: %v", C.GoString(err))
	}

	keyArr := C.g_byte_array_new()
	keyLen := C.element_length_in_bytes(&m)
	C.g_byte_array_set_size(keyArr, C.uint(keyLen))
	C.element_to_bytes(keyArr.data, &m)
	keyArrBytes := C.GoBytes(unsafe.Pointer(keyArr.data), C.int(keyArr.len))
	defer C.element_clear(&m)

	return keyArrBytes, nil
}

func dump(data []byte, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}
