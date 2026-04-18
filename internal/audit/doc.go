// Package audit provides structured, append-only audit logging for
// vaultpull sync operations.
//
// Each sync attempt — successful or not — is recorded as a JSON line
// containing the secret path, number of keys written, whether the run
// was a dry-run, and an optional error message.
//
// Usage:
//
//	f, _ := os.OpenFile("vaultpull.audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
//	l := audit.NewLogger(f)
//	l.Log(audit.Entry{Operation: "sync", SecretPath: path, KeysWritten: n})
package audit
