// Package audit provides structured audit logging and log rotation
// for vaultpull operations.
//
// Use NewLogger to record secret sync events with timestamps and optional
// error context. Use NewRotator with a RotationConfig to manage log file
// size and age, automatically archiving old logs to a backup directory.
package audit
