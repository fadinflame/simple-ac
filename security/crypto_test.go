package security

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	plain := []byte(`{"server":"vpn.example.com","password":"secret"}`)

	blob, err := Encrypt(plain, "master")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decrypt(blob, "master")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Errorf("round trip = %q, want %q", got, plain)
	}
}

func TestEncryptBlobLayout(t *testing.T) {
	plain := []byte("hello")
	blob, err := Encrypt(plain, "pw")
	if err != nil {
		t.Fatal(err)
	}
	// [16 salt][12 nonce][ciphertext + 16-byte GCM tag]
	if want := saltSize + 12 + len(plain) + 16; len(blob) != want {
		t.Errorf("blob length = %d, want %d", len(blob), want)
	}
	if bytes.Contains(blob, plain) {
		t.Error("plaintext must not appear in the blob")
	}
}

func TestEncryptIsRandomised(t *testing.T) {
	a, _ := Encrypt([]byte("same"), "pw")
	b, _ := Encrypt([]byte("same"), "pw")
	if bytes.Equal(a, b) {
		t.Error("two encryptions of the same input must differ (fresh salt and nonce)")
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	blob, _ := Encrypt([]byte("data"), "right")
	if _, err := Decrypt(blob, "wrong"); err == nil {
		t.Error("wrong password must fail")
	}
}

func TestDecryptDetectsTampering(t *testing.T) {
	blob, _ := Encrypt([]byte("data"), "pw")

	for _, idx := range []int{0, saltSize, len(blob) - 1} { // salt, nonce, tag
		tampered := append([]byte(nil), blob...)
		tampered[idx] ^= 0xff
		if _, err := Decrypt(tampered, "pw"); err == nil {
			t.Errorf("flipping byte %d should be detected", idx)
		}
	}
}

func TestDecryptRejectsShortInput(t *testing.T) {
	for _, n := range []int{0, 1, saltSize, saltSize + 11} {
		if _, err := Decrypt(make([]byte, n), "pw"); err == nil {
			t.Errorf("input of %d bytes should be rejected", n)
		}
	}
}

func TestEmptyPlaintextAndPassword(t *testing.T) {
	blob, err := Encrypt(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Decrypt(blob, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got %q, want empty", got)
	}
}
