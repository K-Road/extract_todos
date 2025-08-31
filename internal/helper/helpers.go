package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func String(s string) *string    { return &s }
func Int(i int) *int             { return &i }
func Bool(b bool) *bool          { return &b }
func Float64(f float64) *float64 { return &f }

func HashTodo(file, text string) string {
	s := fmt.Sprintf("%s:%s", file, text)
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
