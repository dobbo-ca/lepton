// Package db provides the database layer using GORM with SQLite.
package db

import (
	"context"
	"fmt"

	"github.com/dobbo-ca/lepton/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB wraps a GORM database connection.
type DB struct {
	db *gorm.DB
}

// Open opens (or creates) a SQLite database at the given path and runs migrations.
func Open(dsn string) (*DB, error) {
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	d := &DB{db: gdb}
	if err := d.Migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return d, nil
}

// Migrate runs auto-migrations for all domain types.
func (d *DB) Migrate() error {
	return d.db.AutoMigrate(
		&domain.Company{},
		&domain.Agent{},
		&domain.Issue{},
		&domain.Run{},
		&domain.Routine{},
		&domain.Project{},
		&domain.Goal{},
		&domain.SecretEntry{},
	)
}

// ---- Agent repository ----

type AgentRepo struct{ db *gorm.DB }

func (d *DB) Agents() *AgentRepo { return &AgentRepo{db: d.db} }

func (r *AgentRepo) List(ctx context.Context, companyID string) ([]domain.Agent, error) {
	var agents []domain.Agent
	return agents, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&agents).Error
}

func (r *AgentRepo) Get(ctx context.Context, id string) (*domain.Agent, error) {
	var a domain.Agent
	if err := r.db.WithContext(ctx).First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AgentRepo) Create(ctx context.Context, a *domain.Agent) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AgentRepo) Update(ctx context.Context, a *domain.Agent) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AgentRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Agent{}, "id = ?", id).Error
}

// ---- Issue repository ----

type IssueRepo struct{ db *gorm.DB }

func (d *DB) Issues() *IssueRepo { return &IssueRepo{db: d.db} }

func (r *IssueRepo) List(ctx context.Context, companyID string) ([]domain.Issue, error) {
	var issues []domain.Issue
	return issues, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&issues).Error
}

func (r *IssueRepo) Get(ctx context.Context, id string) (*domain.Issue, error) {
	var i domain.Issue
	if err := r.db.WithContext(ctx).First(&i, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *IssueRepo) Create(ctx context.Context, i *domain.Issue) error {
	return r.db.WithContext(ctx).Create(i).Error
}

func (r *IssueRepo) Update(ctx context.Context, i *domain.Issue) error {
	return r.db.WithContext(ctx).Save(i).Error
}

// ---- Run repository ----

type RunRepo struct{ db *gorm.DB }

func (d *DB) Runs() *RunRepo { return &RunRepo{db: d.db} }

func (r *RunRepo) List(ctx context.Context, companyID string) ([]domain.Run, error) {
	var runs []domain.Run
	return runs, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&runs).Error
}

func (r *RunRepo) Get(ctx context.Context, id string) (*domain.Run, error) {
	var run domain.Run
	if err := r.db.WithContext(ctx).First(&run, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *RunRepo) Create(ctx context.Context, run *domain.Run) error {
	return r.db.WithContext(ctx).Create(run).Error
}

func (r *RunRepo) Update(ctx context.Context, run *domain.Run) error {
	return r.db.WithContext(ctx).Save(run).Error
}

// ---- Routine repository ----

type RoutineRepo struct{ db *gorm.DB }

func (d *DB) Routines() *RoutineRepo { return &RoutineRepo{db: d.db} }

func (r *RoutineRepo) List(ctx context.Context, companyID string) ([]domain.Routine, error) {
	var routines []domain.Routine
	return routines, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&routines).Error
}

func (r *RoutineRepo) Create(ctx context.Context, routine *domain.Routine) error {
	return r.db.WithContext(ctx).Create(routine).Error
}

// ---- Project repository ----

type ProjectRepo struct{ db *gorm.DB }

func (d *DB) Projects() *ProjectRepo { return &ProjectRepo{db: d.db} }

func (r *ProjectRepo) List(ctx context.Context, companyID string) ([]domain.Project, error) {
	var projects []domain.Project
	return projects, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&projects).Error
}

// ---- Goal repository ----

type GoalRepo struct{ db *gorm.DB }

func (d *DB) Goals() *GoalRepo { return &GoalRepo{db: d.db} }

func (r *GoalRepo) List(ctx context.Context, companyID string) ([]domain.Goal, error) {
	var goals []domain.Goal
	return goals, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&goals).Error
}

// ---- SecretEntry repository ----

// SecretEntryRepo provides access to encrypted secret storage.
// It is intentionally kept internal — callers should use the secrets.SecretStore
// interface rather than this repo directly.
type SecretEntryRepo struct{ db *gorm.DB }

func (d *DB) SecretEntries() *SecretEntryRepo { return &SecretEntryRepo{db: d.db} }

func (r *SecretEntryRepo) List(ctx context.Context, companyID string) ([]domain.SecretEntry, error) {
	var entries []domain.SecretEntry
	return entries, r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&entries).Error
}

func (r *SecretEntryRepo) Get(ctx context.Context, id string) (*domain.SecretEntry, error) {
	var e domain.SecretEntry
	if err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *SecretEntryRepo) Create(ctx context.Context, e *domain.SecretEntry) error {
	return r.db.WithContext(ctx).Create(e).Error
}

// UpdateEncryptedVal replaces the encrypted value for the given secret entry.
func (r *SecretEntryRepo) UpdateEncryptedVal(ctx context.Context, id string, encVal []byte) error {
	return r.db.WithContext(ctx).Model(&domain.SecretEntry{}).
		Where("id = ?", id).
		Updates(map[string]any{"encrypted_val": encVal}).Error
}

func (r *SecretEntryRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.SecretEntry{}, "id = ?", id).Error
}
