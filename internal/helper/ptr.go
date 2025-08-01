package helper

func String(s string) *string    { return &s }
func Int(i int) *int             { return &i }
func Bool(b bool) *bool          { return &b }
func Float64(f float64) *float64 { return &f }
