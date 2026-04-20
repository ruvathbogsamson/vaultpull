package vault

import (
	"fmt"
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of secrets for a given path.
type Snapshot struct {
	Path      string
	Secrets   map[string]string
	CreatedAt time.Time
}

// RollbackManager stores snapshots and allows restoring a previous secret state.
type RollbackManager struct {
	mu        sync.Mutex
	snapshots map[string][]Snapshot
	maxDepth  int
}

// NewRollbackManager creates a RollbackManager that retains up to maxDepth
// snapshots per secret path.
func NewRollbackManager(maxDepth int) *RollbackManager {
	if maxDepth <= 0 {
		maxDepth = 5
	}
	return &RollbackManager{
		snapshots: make(map[string][]Snapshot),
		maxDepth:  maxDepth,
	}
}

// Save records a snapshot for the given path. Older snapshots beyond maxDepth
// are evicted.
func (r *RollbackManager) Save(path string, secrets map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}

	snap := Snapshot{Path: path, Secrets: copy, CreatedAt: time.Now()}
	r.snapshots[path] = append(r.snapshots[path], snap)

	if len(r.snapshots[path]) > r.maxDepth {
		r.snapshots[path] = r.snapshots[path][len(r.snapshots[path])-r.maxDepth:]
	}
}

// Latest returns the most recently saved snapshot for path, or an error if
// none exists.
func (r *RollbackManager) Latest(path string) (Snapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	snaps, ok := r.snapshots[path]
	if !ok || len(snaps) == 0 {
		return Snapshot{}, fmt.Errorf("rollback: no snapshot found for path %q", path)
	}
	return snaps[len(snaps)-1], nil
}

// Previous returns the snapshot one step before the latest, enabling a single
// undo operation.
func (r *RollbackManager) Previous(path string) (Snapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	snaps, ok := r.snapshots[path]
	if !ok || len(snaps) < 2 {
		return Snapshot{}, fmt.Errorf("rollback: no previous snapshot for path %q", path)
	}
	return snaps[len(snaps)-2], nil
}

// Flush removes all stored snapshots for a path.
func (r *RollbackManager) Flush(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.snapshots, path)
}
