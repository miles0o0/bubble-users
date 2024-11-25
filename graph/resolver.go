package graph

import "github.com/miles0o0/bubble-users/database"

// Resolver serves as the dependency injection container for the GraphQL resolvers.
type Resolver struct {
	Database *database.PostgresRepository
}
