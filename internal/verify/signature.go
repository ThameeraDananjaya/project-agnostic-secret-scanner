package verify

import (
	"crypto/ed25519"
	"encoding/base64"
)

const MaxDocumentBytes = 1 << 20

type DetachedSignature struct {
	TrustDomain string `json:"trustDomain"`
	Algorithm   string `json:"algorithm"`
	KeyID       string `json:"keyId"`
	Value       string `json:"value"`
}

type SignatureVerifier interface {
	Verify(message []byte, signature DetachedSignature) error
}

type PublicKey struct {
	KeyID string
	Bytes []byte
}

type Ed25519Verifier struct {
	domain string
	keys   map[string]ed25519.PublicKey
}

func NewEd25519Verifier(domain string, keys []PublicKey) (*Ed25519Verifier, error) {
	if !IsToken(domain) || len(keys) == 0 {
		return nil, ErrInvalidReference
	}
	trusted := make(map[string]ed25519.PublicKey, len(keys))
	for _, key := range keys {
		if !IsToken(key.KeyID) || len(key.Bytes) != ed25519.PublicKeySize {
			return nil, ErrInvalidReference
		}
		if _, duplicate := trusted[key.KeyID]; duplicate {
			return nil, ErrInvalidReference
		}
		trusted[key.KeyID] = append(ed25519.PublicKey(nil), key.Bytes...)
	}
	return &Ed25519Verifier{domain: domain, keys: trusted}, nil
}

func (v *Ed25519Verifier) Verify(message []byte, signature DetachedSignature) error {
	if v == nil || signature.TrustDomain != v.domain || signature.Algorithm != "ed25519" || !IsToken(signature.KeyID) {
		return ErrInvalidSignature
	}
	key, ok := v.keys[signature.KeyID]
	if !ok || len(key) != ed25519.PublicKeySize {
		return ErrInvalidSignature
	}
	decoded, err := base64.RawURLEncoding.DecodeString(signature.Value)
	if err != nil || len(decoded) != ed25519.SignatureSize || base64.RawURLEncoding.EncodeToString(decoded) != signature.Value {
		return ErrInvalidSignature
	}
	if !ed25519.Verify(key, message, decoded) {
		return ErrInvalidSignature
	}
	return nil
}

type DocumentBinding struct {
	SchemaFamily   string            `json:"schemaFamily"`
	SchemaVersion  string            `json:"schemaVersion"`
	AdapterVersion string            `json:"adapterVersion"`
	Digest         string            `json:"digest"`
	Signature      DetachedSignature `json:"signature"`
}

func SignedBindingMessage(binding DocumentBinding) ([]byte, error) {
	if !IsToken(binding.SchemaFamily) || !IsToken(binding.AdapterVersion) || !IsDigest(binding.Digest) ||
		!IsToken(binding.Signature.TrustDomain) || !IsToken(binding.Signature.KeyID) {
		return nil, ErrInvalidReference
	}
	if _, _, err := parseVersion(binding.SchemaVersion); err != nil {
		return nil, err
	}
	p, _ := newPreimage("PSCAN-SIGNED-DOCUMENT-BINDING-V1")
	p.add(binding.Signature.TrustDomain)
	p.add(binding.SchemaFamily)
	p.add(binding.SchemaVersion)
	p.add(binding.AdapterVersion)
	p.add(binding.Signature.KeyID)
	p.add(binding.Digest)
	return p.bytes(), nil
}

func VerifyDocument(raw []byte, binding DocumentBinding, verifier SignatureVerifier) error {
	if len(raw) == 0 || len(raw) > MaxDocumentBytes || verifier == nil {
		return ErrInvalidReference
	}
	if DigestBytes(raw) != binding.Digest {
		return ErrBindingMismatch
	}
	message, err := SignedBindingMessage(binding)
	if err != nil {
		return err
	}
	return verifier.Verify(message, binding.Signature)
}
