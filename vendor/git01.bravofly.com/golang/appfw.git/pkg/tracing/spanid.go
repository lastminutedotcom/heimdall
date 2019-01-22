package tracing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	// this is due to sha256 hash.Sum string being generated at 64 bytes length, so that "runtime error: slice bounds out of range" is avoided
	maxSpanIDRange = 57
)

func generateSpanID(seed int) string {
	if seed == 0 {
		return "0000000"
	}

	spanID := sha256.New()
	spanID.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	hash := hex.EncodeToString(spanID.Sum(nil))
	return strings.Replace(strings.ToLower(fmt.Sprintf("%s", hash[seed:seed+7])), " ", "0", -1)
}

// BaseSpanID gets the first call spanID: you want to call this when there is no spanID set so far, so that the base will be generated
func BaseSpanID() string {
	return generateSpanID(0)
}

// NewSpanID creates a new spanID: you want to call this when there is already a spanID set for the request/response and you will create a new one
func NewSpanID() string {
	return generateSpanID(rand.Intn(maxSpanIDRange))
}
