// Package audit provides structured audit logging for vaultpull operations.
//
// It supports newline-delimited JSON entries via Shipper, automatic log
// rotation via Rotator, and low-level append logging via Logger.
//
// Typical usage:
//
//	shipper, err := audit.NewShipper("/var/log/vaultpull/audit.log", audit.DefaultRotationConfig())
//	if err != nil { ... }
//	shipper.Ship(audit.Entry{
//		Operation: "sync",
//		Path:      "secret/data/myapp",
//		Namespace: "MYAPP",
//		Keys:      []string{"MYAPP_DB"},
//	})
package audit
