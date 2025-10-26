package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maixuanbach174/online-course-app/internal/education/adapters/postgresql/database"
	"github.com/maixuanbach174/online-course-app/internal/education/domain/module"
	"github.com/pkg/errors"
)

type ModuleRepository struct {
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewModuleRepository(db *pgxpool.Pool) *ModuleRepository {
	return &ModuleRepository{
		db:      db,
		queries: database.New(db),
	}
}

// Create implements course.ModuleRepository
func (r *ModuleRepository) Create(ctx context.Context, m *module.Module) error {
	params := database.CreateModuleParams{
		ID:         m.ID(),
		CourseID:   m.CourseID(),
		Title:      m.Title(),
		OrderIndex: int32(m.Order()),
	}

	if err := r.queries.CreateModule(ctx, params); err != nil {
		return errors.Wrap(err, "failed to create module")
	}

	return nil
}

// Update implements course.ModuleRepository
func (r *ModuleRepository) Update(ctx context.Context, m *module.Module) error {
	params := database.UpdateModuleParams{
		ID:         m.ID(),
		Title:      m.Title(),
		OrderIndex: int32(m.Order()),
	}

	if err := r.queries.UpdateModule(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update module")
	}

	return nil
}

// Delete implements course.ModuleRepository
func (r *ModuleRepository) Delete(ctx context.Context, id string) error {
	if err := r.queries.DeleteModule(ctx, id); err != nil {
		return errors.Wrap(err, "failed to delete module")
	}

	return nil
}

// Get implements course.ModuleRepository
func (r *ModuleRepository) Get(ctx context.Context, id string) (*module.Module, error) {
	dbModule, err := r.queries.GetModuleByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get module")
	}

	return r.toDomainModule(dbModule)
}

// GetByCourseID implements course.ModuleRepository
func (r *ModuleRepository) GetByCourseID(ctx context.Context, courseID string) ([]*module.Module, error) {
	dbModules, err := r.queries.GetModulesByCourseID(ctx, courseID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get modules by course")
	}

	modules := make([]*module.Module, 0, len(dbModules))
	for _, dbModule := range dbModules {
		domainModule, err := r.toDomainModule(dbModule)
		if err != nil {
			return nil, err
		}
		modules = append(modules, domainModule)
	}

	return modules, nil
}

// Exists implements course.ModuleRepository
func (r *ModuleRepository) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := r.queries.ModuleExists(ctx, id)
	if err != nil {
		return false, errors.Wrap(err, "failed to check module existence")
	}
	return exists, nil
}

// UpdateOrder implements course.ModuleRepository
func (r *ModuleRepository) UpdateOrder(ctx context.Context, id string, order int) error {
	params := database.UpdateModuleOrderParams{
		ID:         id,
		OrderIndex: int32(order),
	}

	if err := r.queries.UpdateModuleOrder(ctx, params); err != nil {
		return errors.Wrap(err, "failed to update module order")
	}

	return nil
}

// ReorderModules implements course.ModuleRepository
func (r *ModuleRepository) ReorderModules(ctx context.Context, moduleOrders map[string]int) error {
	// Start a transaction for atomic update
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Update each module's order
	for moduleID, order := range moduleOrders {
		params := database.UpdateModuleOrderParams{
			ID:         moduleID,
			OrderIndex: int32(order),
		}
		if err := qtx.UpdateModuleOrder(ctx, params); err != nil {
			return errors.Wrap(err, "failed to update module order")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// Helper methods

func (r *ModuleRepository) toDomainModule(dbModule database.Module) (*module.Module, error) {
	return module.NewModule(
		dbModule.ID,
		dbModule.CourseID,
		dbModule.Title,
		int(dbModule.OrderIndex),
	)
}
