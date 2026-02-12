package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	migrationCollectionName = "migrations"
)

// Migration represents a single migration
type Migration struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	AppliedAt   time.Time `bson:"applied_at"`
	Description string    `bson:"description,omitempty"`
}

// MongoMigrationFunc is a function that performs a migration
type MongoMigrationFunc func(ctx context.Context, db *mongo.Database) error

// MongoMigration represents a migration with its metadata
type MongoMigration struct {
	Name        string
	Description string
	Up          MongoMigrationFunc
	Down        MongoMigrationFunc // Optional: for rollback support
	UpSQL       string             // Optional: MongoDB JavaScript command as string
	DownSQL     string             // Optional: MongoDB JavaScript command as string
}

// MongoMigrations is the list of all MongoDB migrations
var MongoMigrations = []MongoMigration{
	{
		Name:        "001_create_devices_collection",
		Description: "Create devices collection with indexes",
		Up:          createDevicesCollection,
	},
	{
		Name:        "002_create_sensors_collection",
		Description: "Create sensors collection with indexes",
		Up:          createSensorsCollection,
	},
	{
		Name:        "003_create_sensor_readings_collection",
		Description: "Create sensor_readings collection with indexes",
		Up:          createSensorReadingsCollection,
	},
}

// RunMongoMigrations runs all pending MongoDB migrations
func RunMongoMigrations(client *mongo.Client) error {
	if client == nil {
		return fmt.Errorf("MongoDB client is nil")
	}

	dbName := viper.GetString("MONGO_DATABASE")
	if dbName == "" {
		return fmt.Errorf("MONGO_DATABASE is not set")
	}

	db := client.Database(dbName)
	ctx := context.Background()

	// Ensure migrations collection exists
	if err := ensureMigrationsCollection(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure migrations collection: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range MongoMigrations {
		if isApplied(appliedMigrations, migration.Name) {
			fmt.Printf("Migration %s already applied, skipping\n", migration.Name)
			continue
		}

		fmt.Printf("Running migration: %s - %s\n", migration.Name, migration.Description)

		// Execute migration - prefer SQL string if provided, otherwise use function
		var err error
		if migration.UpSQL != "" {
			err = executeMongoScript(ctx, db, migration.UpSQL)
		} else if migration.Up != nil {
			err = migration.Up(ctx, db)
		} else {
			return fmt.Errorf("migration %s has no Up function or UpSQL string", migration.Name)
		}

		if err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Name, err)
		}

		// Record migration as applied
		if err := recordMigration(ctx, db, migration.Name, migration.Description); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		fmt.Printf("Migration %s completed successfully\n", migration.Name)
	}

	return nil
}

// ensureMigrationsCollection creates the migrations tracking collection if it doesn't exist
func ensureMigrationsCollection(ctx context.Context, db *mongo.Database) error {
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": migrationCollectionName})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		// Create collection with unique index on name
		if err := db.CreateCollection(ctx, migrationCollectionName); err != nil {
			// Collection might already exist, ignore error
			if cmdErr, ok := err.(mongo.CommandError); !ok || cmdErr.Code != 48 {
				return err
			}
		}

		// Create unique index on migration name
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		_, err = db.Collection(migrationCollectionName).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return err
		}
	}

	return nil
}

// getAppliedMigrations returns a map of applied migration names
func getAppliedMigrations(ctx context.Context, db *mongo.Database) (map[string]bool, error) {
	cursor, err := db.Collection(migrationCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	applied := make(map[string]bool)
	for cursor.Next(ctx) {
		var migration Migration
		if err := cursor.Decode(&migration); err != nil {
			return nil, err
		}
		applied[migration.Name] = true
	}

	return applied, cursor.Err()
}

// isApplied checks if a migration has been applied
func isApplied(applied map[string]bool, name string) bool {
	return applied[name]
}

// recordMigration records a migration as applied
func recordMigration(ctx context.Context, db *mongo.Database, name, description string) error {
	migration := Migration{
		ID:          name,
		Name:        name,
		AppliedAt:   time.Now().UTC(),
		Description: description,
	}

	_, err := db.Collection(migrationCollectionName).InsertOne(ctx, migration)
	return err
}

// executeMongoScript executes a MongoDB JavaScript command string
// The script can contain multiple MongoDB operations separated by semicolons
func executeMongoScript(ctx context.Context, db *mongo.Database, script string) error {
	// Parse and execute MongoDB commands from the script string
	lines := splitScriptLines(script)
	for _, line := range lines {
		line = trimSpace(line)
		if line == "" || (len(line) > 0 && line[0] == '#') {
			continue // Skip empty lines and comments
		}

		if err := executeMongoCommand(ctx, db, line); err != nil {
			return fmt.Errorf("failed to execute command: %s, error: %w", line, err)
		}
	}

	return nil
}

// splitScriptLines splits a script into individual command lines
func splitScriptLines(script string) []string {
	var lines []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(script); i++ {
		char := script[i]

		if !inQuotes && (char == '"' || char == '\'') {
			inQuotes = true
			quoteChar = char
		} else if inQuotes && char == quoteChar && (i == 0 || script[i-1] != '\\') {
			inQuotes = false
		}

		if !inQuotes && char == ';' {
			line := current.String()
			if strings.TrimSpace(line) != "" {
				lines = append(lines, line)
			}
			current.Reset()
			continue
		}

		current.WriteByte(char)
	}

	// Add last line if not empty
	lastLine := current.String()
	if strings.TrimSpace(lastLine) != "" {
		lines = append(lines, lastLine)
	}

	return lines
}

// executeMongoCommand executes a single MongoDB command
// Supports declarative format: COMMAND:collection:jsonOptions
func executeMongoCommand(ctx context.Context, db *mongo.Database, command string) error {
	parts := strings.SplitN(command, ":", 3)
	if len(parts) < 2 {
		return fmt.Errorf("invalid command format: %s (expected COMMAND:collection:options)", command)
	}

	cmd := strings.TrimSpace(parts[0])
	collectionName := strings.TrimSpace(parts[1])
	var optionsJSON string
	if len(parts) > 2 {
		optionsJSON = strings.TrimSpace(parts[2])
	}

	switch cmd {
	case "CREATE_COLLECTION":
		return createCollectionFromString(ctx, db, collectionName, optionsJSON)
	case "CREATE_INDEX":
		return createIndexFromString(ctx, db, collectionName, optionsJSON)
	case "UPDATE_MANY":
		return updateManyFromString(ctx, db, collectionName, optionsJSON)
	case "INSERT_MANY":
		return insertManyFromString(ctx, db, collectionName, optionsJSON)
	case "DELETE_MANY":
		return deleteManyFromString(ctx, db, collectionName, optionsJSON)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// createCollectionFromString creates a collection from JSON options
func createCollectionFromString(ctx context.Context, db *mongo.Database, collectionName, optionsJSON string) error {
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	if len(collections) > 0 {
		return nil // Collection already exists
	}

	var opts options.CreateCollectionOptions
	if optionsJSON != "" {
		var optsMap map[string]interface{}
		if err := json.Unmarshal([]byte(optionsJSON), &optsMap); err == nil {
			if capped, ok := optsMap["capped"].(bool); ok {
				opts.SetCapped(capped)
			}
			if size, ok := optsMap["size"].(float64); ok {
				opts.SetSizeInBytes(int64(size))
			}
			if max, ok := optsMap["max"].(float64); ok {
				opts.SetMaxDocuments(int64(max))
			}
		}
	}

	return db.CreateCollection(ctx, collectionName, &opts)
}

// createIndexFromString creates an index from JSON options
func createIndexFromString(ctx context.Context, db *mongo.Database, collectionName, optionsJSON string) error {
	if optionsJSON == "" {
		return fmt.Errorf("index options required")
	}

	var indexSpec struct {
		Keys    bson.M                 `json:"keys"`
		Options map[string]interface{} `json:"options"`
	}

	if err := json.Unmarshal([]byte(optionsJSON), &indexSpec); err != nil {
		return fmt.Errorf("invalid index JSON: %w", err)
	}

	keys := bson.D{}
	for k, v := range indexSpec.Keys {
		keys = append(keys, bson.E{Key: k, Value: v})
	}

	indexOpts := options.Index()
	if indexSpec.Options != nil {
		if unique, ok := indexSpec.Options["unique"].(bool); ok && unique {
			indexOpts.SetUnique(true)
		}
		if name, ok := indexSpec.Options["name"].(string); ok {
			indexOpts.SetName(name)
		}
		if sparse, ok := indexSpec.Options["sparse"].(bool); ok {
			indexOpts.SetSparse(sparse)
		}
		if ttl, ok := indexSpec.Options["expireAfterSeconds"].(float64); ok {
			indexOpts.SetExpireAfterSeconds(int32(ttl))
		}
	}

	collection := db.Collection(collectionName)
	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: indexOpts,
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// If index already exists, that's okay
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 85 {
			return nil
		}
		return err
	}

	return nil
}

// updateManyFromString updates documents from JSON options
func updateManyFromString(ctx context.Context, db *mongo.Database, collectionName, optionsJSON string) error {
	if optionsJSON == "" {
		return fmt.Errorf("update options required")
	}

	var updateSpec struct {
		Filter bson.M `json:"filter"`
		Update bson.M `json:"update"`
	}

	if err := json.Unmarshal([]byte(optionsJSON), &updateSpec); err != nil {
		return fmt.Errorf("invalid update JSON: %w", err)
	}

	collection := db.Collection(collectionName)
	_, err := collection.UpdateMany(ctx, updateSpec.Filter, updateSpec.Update)
	return err
}

// insertManyFromString inserts documents from JSON options
func insertManyFromString(ctx context.Context, db *mongo.Database, collectionName, optionsJSON string) error {
	if optionsJSON == "" {
		return fmt.Errorf("insert documents required")
	}

	var documents []bson.M
	if err := json.Unmarshal([]byte(optionsJSON), &documents); err != nil {
		return fmt.Errorf("invalid insert JSON: %w", err)
	}

	// Convert []bson.M to []interface{}
	docs := make([]interface{}, len(documents))
	for i, doc := range documents {
		docs[i] = doc
	}

	collection := db.Collection(collectionName)
	_, err := collection.InsertMany(ctx, docs)
	return err
}

// deleteManyFromString deletes documents from JSON options
func deleteManyFromString(ctx context.Context, db *mongo.Database, collectionName, optionsJSON string) error {
	if optionsJSON == "" {
		return fmt.Errorf("delete filter required")
	}

	var filter bson.M
	if err := json.Unmarshal([]byte(optionsJSON), &filter); err != nil {
		return fmt.Errorf("invalid delete JSON: %w", err)
	}

	collection := db.Collection(collectionName)
	_, err := collection.DeleteMany(ctx, filter)
	return err
}

// Helper functions
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

// createDevicesCollection creates the devices collection with indexes
func createDevicesCollection(ctx context.Context, db *mongo.Database) error {
	collectionName := "devices"

	// Create collection if it doesn't exist
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			if cmdErr, ok := err.(mongo.CommandError); !ok || cmdErr.Code != 48 {
				return err
			}
		}
	}

	// Create indexes
	collection := db.Collection(collectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "device_code", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_device_code"),
		},
		{
			Keys:    bson.D{{Key: "is_active", Value: 1}},
			Options: options.Index().SetName("idx_is_active"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// If indexes already exist, that's okay
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 85 {
			return nil
		}
		return err
	}

	return nil
}

// createSensorsCollection creates the sensors collection with indexes
func createSensorsCollection(ctx context.Context, db *mongo.Database) error {
	collectionName := "sensors"

	// Create collection if it doesn't exist
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			if cmdErr, ok := err.(mongo.CommandError); !ok || cmdErr.Code != 48 {
				return err
			}
		}
	}

	// Create indexes
	collection := db.Collection(collectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "device_id", Value: 1}},
			Options: options.Index().SetName("idx_device_id"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "device_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_name_device_id"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// If indexes already exist, that's okay
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 85 {
			return nil
		}
		return err
	}

	return nil
}

// createSensorReadingsCollection creates the sensor_readings collection with indexes
func createSensorReadingsCollection(ctx context.Context, db *mongo.Database) error {
	collectionName := "sensor_readings"

	// Create collection if it doesn't exist
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return err
	}

	if len(collections) == 0 {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			if cmdErr, ok := err.(mongo.CommandError); !ok || cmdErr.Code != 48 {
				return err
			}
		}
	}

	// Create indexes
	collection := db.Collection(collectionName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "device_id", Value: 1}, {Key: "timestamp", Value: -1}},
			Options: options.Index().SetName("idx_device_timestamp"),
		},
		{
			Keys:    bson.D{{Key: "sensor_id", Value: 1}, {Key: "timestamp", Value: -1}},
			Options: options.Index().SetName("idx_sensor_timestamp"),
		},
		{
			Keys:    bson.D{{Key: "timestamp", Value: -1}},
			Options: options.Index().SetName("idx_timestamp"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at"),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// If indexes already exist, that's okay
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 85 {
			return nil
		}
		return err
	}

	return nil
}
