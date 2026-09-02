// Package verify validates public references without taking custody of the
// referenced evidence or any signing key.
package verify

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	tokenPattern  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/+-]{0,127}$`)
	uuidPattern   = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

var (
	ErrInvalidReference  = errors.New("reference is invalid")
	ErrUnsupportedSchema = errors.New("schema is unsupported")
	ErrInvalidSignature  = errors.New("signature is invalid")
	ErrBindingMismatch   = errors.New("binding does not match")
	ErrRevoked           = errors.New("reference is revoked")
	ErrEvidenceConflict  = errors.New("evidence chain conflicts")
	ErrEvidenceRollback  = errors.New("evidence chain rolled back")
)

func IsDigest(value string) bool { return digestPattern.MatchString(value) }
func IsToken(value string) bool  { return tokenPattern.MatchString(value) }
func IsUUID(value string) bool   { return uuidPattern.MatchString(value) }

// DigestBytes returns the lowercase SHA-256 of the exact admitted bytes.
func DigestBytes(value []byte) string {
	digest := sha256.Sum256(value)
	return hex.EncodeToString(digest[:])
}

type preimage struct{ value []byte }

func newPreimage(domain string) (*preimage, error) {
	if !IsToken(domain) {
		return nil, ErrInvalidReference
	}
	p := &preimage{}
	p.add(domain)
	return p, nil
}

func (p *preimage) add(value string) {
	var size [8]byte
	binary.BigEndian.PutUint64(size[:], uint64(len([]byte(value))))
	p.value = append(p.value, size[:]...)
	p.value = append(p.value, []byte(value)...)
}

func (p *preimage) addInt(value int64) { p.add(strconv.FormatInt(value, 10)) }
func (p *preimage) bytes() []byte      { return append([]byte(nil), p.value...) }

func parseVersion(value string) (int, int, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" ||
		(len(parts[0]) > 1 && parts[0][0] == '0') ||
		(len(parts[1]) > 1 && parts[1][0] == '0') {
		return 0, 0, ErrUnsupportedSchema
	}
	major, majorErr := strconv.Atoi(parts[0])
	minor, minorErr := strconv.Atoi(parts[1])
	if majorErr != nil || minorErr != nil || major < 1 || minor < 0 {
		return 0, 0, ErrUnsupportedSchema
	}
	return major, minor, nil
}

func parseCanonicalTime(value string) (time.Time, error) {
	if !strings.HasSuffix(value, "Z") {
		return time.Time{}, ErrInvalidReference
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil || parsed.Location() != time.UTC || parsed.Format(time.RFC3339Nano) != value {
		return time.Time{}, ErrInvalidReference
	}
	return parsed, nil
}

func ParseCanonicalTime(value string) (time.Time, error) { return parseCanonicalTime(value) }
