package gcrypto

import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
)

func MD5(s string) string    { h := md5.Sum([]byte(s)); return hex.EncodeToString(h[:]) }
func SHA1(s string) string   { h := sha1.Sum([]byte(s)); return hex.EncodeToString(h[:]) }
func SHA256(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }

func MD5Bytes(b []byte) string    { h := md5.Sum(b); return hex.EncodeToString(h[:]) }
func SHA1Bytes(b []byte) string   { h := sha1.Sum(b); return hex.EncodeToString(h[:]) }
func SHA256Bytes(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }

func Base64Encode(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }
func Base64Decode(s string) string {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(data)
}

func SM3(s string) string      { return SM3Bytes([]byte(s)) }
func SM3Bytes(b []byte) string { h := sm3.Sm3Sum(b); return hex.EncodeToString(h[:]) }

func SM2Sign(prikey []byte, msg []byte) ([]byte, error) {
	// 创建sm2对象
	sm2Prikey, err := ToSM2(prikey)
	if err != nil {
		return nil, err
	}
	return sm2Prikey.Sign(rand.Reader, msg, nil)
}

func SM2Verify(pubkey []byte, msg []byte, sign []byte) bool {
	sm2PubKey, err := UnmarshalPubkey(pubkey)
	if err != nil {
		return false
	}
	return sm2PubKey.Verify(msg, sign)
}

// sm2加密
func SM2Encrypt(pubkey []byte, msg []byte) ([]byte, error) {
	sm2PubKey, err := UnmarshalPubkey(pubkey)
	if err != nil {
		return nil, err
	}
	return sm2.EncryptAsn1(sm2PubKey, msg, rand.Reader)
}

func SM2Decrypt(prikey []byte, ciphertext []byte) ([]byte, error) {
	sm2Prikey, err := ToSM2(prikey)
	if err != nil {
		return nil, err
	}
	return sm2.DecryptAsn1(sm2Prikey, ciphertext)
}

var errInvalidPubkey = errors.New("invalid sm2 public key")

func UnmarshalPubkey(pub []byte) (*sm2.PublicKey, error) {
	x, y := elliptic.Unmarshal(sm2.P256Sm2(), pub)
	if x == nil {
		return nil, errInvalidPubkey
	}
	return &sm2.PublicKey{Curve: sm2.P256Sm2(), X: x, Y: y}, nil
}

func FromSM2Pub(pub *sm2.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(sm2.P256Sm2(), pub.X, pub.Y)
}

func toSM2(d []byte, strict bool) (*sm2.PrivateKey, error) {
	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = sm2.P256Sm2()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(sm2.P256Sm2().Params().N) >= 0 {
		return nil, errors.New("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, errors.New("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

func ToSM2(d []byte) (*sm2.PrivateKey, error) {
	return toSM2(d, true)
}
