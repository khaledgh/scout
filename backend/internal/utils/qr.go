package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GenerateQRToken creates a signed token for QR-based attendance check-in.
// Format: "<memberID>.<unixTimestamp>.<hmac>"
func GenerateQRToken(memberID uint, secret string) string {
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	msg := fmt.Sprintf("%d.%s", memberID, ts)
	mac := hmacSHA256(msg, secret)
	return fmt.Sprintf("%s.%s", msg, mac)
}

// VerifyQRToken validates the token and returns the memberID and whether it is within ttl.
func VerifyQRToken(token, secret string, ttl time.Duration) (uint, bool) {
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 {
		return 0, false
	}
	memberID64, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, false
	}
	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, false
	}
	if time.Since(time.Unix(ts, 0)) > ttl {
		return 0, false
	}
	msg := fmt.Sprintf("%s.%s", parts[0], parts[1])
	expected := hmacSHA256(msg, secret)
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return 0, false
	}
	return uint(memberID64), true
}

func hmacSHA256(msg, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}
