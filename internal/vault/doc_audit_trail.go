// Package vault provides the AuditTrail type for recording secret access
// and modification events during a vaultpull session.
//
// AuditTrail captures:
//   - The operation performed (read, write, delete)
//   - The secret path accessed
//   - The namespace filter in use
//   - The keys returned or modified
//   - Any error that occurred
//
// Usage:
//
//	trail := vault.NewAuditTrail()
//	trail.Record("read", "secret/myapp", "prod", []string{"DB_URL"}, nil)
//	fmt.Print(trail.Summary())
//
// The trail is intended for in-process use during a sync run. For persistent
// audit logging, use the internal/audit package with its Shipper and Rotator.
package vault
