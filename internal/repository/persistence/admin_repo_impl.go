package persistence

import (
	"context"
	"database/sql"
	"errors"
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
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByUsername(ctx context.Context, username string) (*model.Admin, error) {
	if r.db != nil {
		return r.findOne(ctx, `SELECT id, username, password_hash, display_name, status, created_at, updated_at
FROM admins
WHERE username = ?`, username)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.adminIDsByName[username]
	if !ok {
		return nil, nil
	}

	admin := r.adminsByID[id]
	return cloneAdmin(admin), nil
}

func (r *AdminRepository) FindByID(ctx context.Context, id int64) (*model.Admin, error) {
	if r.db != nil {
		return r.findOne(ctx, `SELECT id, username, password_hash, display_name, status, created_at, updated_at
FROM admins
WHERE id = ?`, id)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	admin, ok := r.adminsByID[id]
	if !ok {
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

		_, err := r.db.ExecContext(ctx, `INSERT INTO admins (id, username, password_hash, display_name, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
	username = VALUES(username),
	password_hash = VALUES(password_hash),
	display_name = VALUES(display_name),
	status = VALUES(status),
	updated_at = VALUES(updated_at)`,
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
	r.adminsByID[admin.ID] = admin
	r.adminIDsByName[admin.Username] = admin.ID
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

		_, err := r.db.ExecContext(ctx, `INSERT INTO admins (id, username, password_hash, display_name, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
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

func (r *AdminRepository) findOne(ctx context.Context, query string, arg any) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.Name,
		&admin.Status,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
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
