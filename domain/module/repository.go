package module

import "context"

// ModuleRepository manages Module aggregate persistence
type ModuleRepository interface {
	// Create saves a new module to the database
	Create(ctx context.Context, module *Module) error

	// Update modifies an existing module
	Update(ctx context.Context, module *Module) error

	// Delete removes a module and all associated lessons and exercises via CASCADE
	Delete(ctx context.Context, id string) error

	// Get retrieves a module by ID
	Get(ctx context.Context, id string) (*Module, error)

	// GetByCourseID retrieves all modules for a specific course, ordered by the order field
	GetByCourseID(ctx context.Context, courseID string) ([]*Module, error)

	// Exists checks if a module with the given ID exists
	Exists(ctx context.Context, id string) (bool, error)

	// ReorderModules updates the order of multiple modules in a single transaction
	ReorderModules(ctx context.Context, moduleOrders map[string]int) error
}
