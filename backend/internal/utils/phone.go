package utils

import (
	"regexp"
	"strings"
)

var nonDigit = regexp.MustCompile(`[^\d]`)

// NormalizePhone normalises a Lebanese phone number to E.164 (+961XXXXXXXX).
// Accepts: 03123456 | 70123456 | +96103123456 | 0096103123456 | 961-03-123-456
func NormalizePhone(phone string) string {
	digits := nonDigit.ReplaceAllString(phone, "")

	switch {
	case strings.HasPrefix(digits, "00961") && len(digits) == 13:
		return "+" + digits[2:]
	case strings.HasPrefix(digits, "961") && len(digits) == 11:
		return "+" + digits
	case len(digits) == 8:
		return "+961" + digits
	default:
		return "+" + digits
	}
}
