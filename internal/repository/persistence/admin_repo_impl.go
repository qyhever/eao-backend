package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"eao/internal/model"
	"eao/internal/repository"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

var (
	adminIDSequence atomic.Int64
)

type AdminRepository struct {
	db             *sql.DB
	mu             sync.RWMutex
	nextID         int64
	adminsByID     map[int64]model.Admin
	adminIDsByName map[string]int64
}

func NewAdminRepository(db *sql.DB) repository.AdminRepository {
	return &AdminRepository{
		db:             db,
		nextID:         1,
		adminsByID:     make(map[int64]model.Admin),
		adminIDsByName: make(map[string]int64),
	}
}

func (r *AdminRepository) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	if r.db != nil {
		return r.findOne(ctx, `SELECT id, username, password_hash, display_name, status, created_at, updated_at, deleted_at
FROM admins
WHERE username = ? AND deleted_at IS NULL`, username)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.adminIDsByName[username]
	if !ok {
		return nil, nil
	}

	admin := r.adminsByID[id]
	if admin.DeletedAt != nil {
		return nil, nil
	}
	return cloneAdmin(admin), nil
}

func (r *AdminRepository) FindByID(ctx context.Context, id int64) (*model.Admin, error) {
	if r.db != nil {
		return r.findOne(ctx, `SELECT id, username, password_hash, display_name, status, created_at, updated_at, deleted_at
FROM admins
WHERE id = ? AND deleted_at IS NULL`, id)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	admin, ok := r.adminsByID[id]
	if !ok {
		return nil, nil
	}
	if admin.DeletedAt != nil {
		return nil, nil
	}
	return cloneAdmin(admin), nil
}

func (r *AdminRepository) Upsert(ctx context.Context, admin model.Admin) error {
	if r.db != nil {
		now := time.Now()
		if admin.Status == "" {
			admin.Status = "active"
		}
		if admin.CreatedAt.IsZero() {
			admin.CreatedAt = now
		}
		admin.UpdatedAt = now

		_, err := r.db.ExecContext(ctx, `INSERT INTO admins (id, username, password_hash, display_name, status, created_at, updated_at, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NULL)
ON DUPLICATE KEY UPDATE
	username = VALUES(username),
	password_hash = VALUES(password_hash),
	display_name = VALUES(display_name),
	status = VALUES(status),
	updated_at = VALUES(updated_at),
	deleted_at = NULL`,
			admin.ID,
			admin.Username,
			admin.PasswordHash,
			admin.Name,
			admin.Status,
			admin.CreatedAt,
			admin.UpdatedAt,
		)
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if admin.Status == "" {
		admin.Status = "active"
	}
	now := time.Now()
	if admin.CreatedAt.IsZero() {
		admin.CreatedAt = now
	}
	admin.UpdatedAt = now
	if existing, ok := r.adminsByID[admin.ID]; ok && existing.Username != admin.Username {
		delete(r.adminIDsByName, existing.Username)
	}
	admin.DeletedAt = nil
	r.adminsByID[admin.ID] = admin
	r.adminIDsByName[admin.Username] = admin.ID
	if admin.ID >= r.nextID {
		r.nextID = admin.ID + 1
	}
	return nil
}

func (r *AdminRepository) Create(ctx context.Context, admin model.Admin) (*model.Admin, error) {
	if r.db != nil {
		now := time.Now()
		if admin.Status == "" {
			admin.Status = "active"
		}
		admin.ID = nextAdminID()
		admin.CreatedAt = now
		admin.UpdatedAt = now

		_, err := r.db.ExecContext(ctx, `INSERT INTO admins (id, username, password_hash, display_name, status, created_at, updated_at, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NULL)`,
			admin.ID,
			admin.Username,
			admin.PasswordHash,
			admin.Name,
			admin.Status,
			admin.CreatedAt,
			admin.UpdatedAt,
		)
		if isDuplicateAdminUsernameError(err) {
			return nil, repository.ErrAdminUsernameAlreadyExists
		}
		if err != nil {
			return nil, err
		}
		return r.FindByID(ctx, admin.ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.adminIDsByName[admin.Username]; exists {
		return nil, repository.ErrAdminUsernameAlreadyExists
	}
	if admin.Status == "" {
		admin.Status = "active"
	}
	now := time.Now()
	admin.ID = r.nextID
	r.nextID++
	admin.CreatedAt = now
	admin.UpdatedAt = now
	r.adminsByID[admin.ID] = admin
	r.adminIDsByName[admin.Username] = admin.ID
	return cloneAdmin(admin), nil
}

func (r *AdminRepository) Update(ctx context.Context, admin model.Admin) (*model.Admin, error) {
	if r.db != nil {
		admin.UpdatedAt = time.Now()
		_, err := r.db.ExecContext(ctx, `UPDATE admins
SET username = ?, password_hash = ?, display_name = ?, status = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL`,
			admin.Username,
			admin.PasswordHash,
			admin.Name,
			admin.Status,
			admin.UpdatedAt,
			admin.ID,
		)
		if isDuplicateAdminUsernameError(err) {
			return nil, repository.ErrAdminUsernameAlreadyExists
		}
		if err != nil {
			return nil, err
		}
		return r.FindByID(ctx, admin.ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.adminsByID[admin.ID]
	if !ok || current.DeletedAt != nil {
		return nil, nil
	}
	if existingID, exists := r.adminIDsByName[admin.Username]; exists && existingID != admin.ID {
		return nil, repository.ErrAdminUsernameAlreadyExists
	}
	if current.Username != admin.Username {
		delete(r.adminIDsByName, current.Username)
	}
	admin.CreatedAt = current.CreatedAt
	admin.UpdatedAt = time.Now()
	admin.DeletedAt = nil
	r.adminsByID[admin.ID] = admin
	r.adminIDsByName[admin.Username] = admin.ID
	return cloneAdmin(admin), nil
}

func (r *AdminRepository) BatchSoftDelete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	if r.db != nil {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
		args := make([]any, 0, len(ids))
		for _, id := range ids {
			args = append(args, id)
		}
		_, err := r.db.ExecContext(ctx, fmt.Sprintf(`UPDATE admins
SET deleted_at = NOW(), updated_at = NOW()
WHERE id IN (%s) AND deleted_at IS NULL`, placeholders), args...)
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, id := range ids {
		admin, ok := r.adminsByID[id]
		if !ok || admin.DeletedAt != nil {
			continue
		}
		admin.DeletedAt = &now
		admin.UpdatedAt = now
		r.adminsByID[id] = admin
	}
	return nil
}

func (r *AdminRepository) UpdateStatus(ctx context.Context, id int64, status string) (*model.Admin, error) {
	if r.db != nil {
		_, err := r.db.ExecContext(ctx, `UPDATE admins
SET status = ?, updated_at = NOW()
WHERE id = ? AND deleted_at IS NULL`, status, id)
		if err != nil {
			return nil, err
		}
		return r.FindByID(ctx, id)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	admin, ok := r.adminsByID[id]
	if !ok || admin.DeletedAt != nil {
		return nil, nil
	}
	admin.Status = status
	admin.UpdatedAt = time.Now()
	r.adminsByID[id] = admin
	return cloneAdmin(admin), nil
}

func (r *AdminRepository) findOne(ctx context.Context, query string, arg any) (*model.Admin, error) {
	var admin model.Admin
	var deletedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Name,
		&admin.Status,
		&admin.CreatedAt,
		&admin.UpdatedAt,
		&deletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		admin.DeletedAt = &deletedAt.Time
	}
	return &admin, nil
}

func cloneAdmin(admin model.Admin) *model.Admin {
	copied := admin
	return &copied
}

func nextAdminID() int64 {
	for {
		now := time.Now().UnixNano()
		current := adminIDSequence.Load()
		next := now
		if current >= next {
			next = current + 1
		}
		if adminIDSequence.CompareAndSwap(current, next) {
			return next
		}
	}
}

func isDuplicateAdminUsernameError(err error) bool {
	var mysqlErr *mysqlDriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
