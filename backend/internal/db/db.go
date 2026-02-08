package db

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"badminton_tournament/backend/internal/models"
)

var DB *bun.DB

// Connect initializes the database connection with logic to handle Neon.tech on Render
func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	// Auto-fix for corrupted DSNs (e.g. double copy-paste in Render Dashboard)
	// Example corruption: postgresql://user:pass://postgresql://real_user:real_pass@host...
	if strings.Count(dsn, "://") > 1 {
		log.Printf("WARNING: Detected malformed DATABASE_URL with multiple schemes. Attempting to auto-fix.")
		
		// Find the last occurrence of the scheme separator
		lastIndex := strings.LastIndex(dsn, "://")
		if lastIndex != -1 {
			// Backtrack to find the start of the scheme (e.g., "postgresql" or "postgres")
			// Simple heuristic: just look for the last "postgres" before the "://"
			prefixStart := strings.LastIndex(dsn[:lastIndex], "postgres")
			if prefixStart != -1 {
				dsn = dsn[prefixStart:]
				log.Printf("DEBUG: Auto-fixed DSN to start from index %d", prefixStart)
			}
		}
	}

	// 1. Parse the URL to extract host and fix options
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	// 2. Extract specific Neon endpoint ID for "options=endpoint=..."
	// Host format is usually: ep-xyz-123.region.aws.neon.tech
	host := parsedURL.Hostname()
	// Password is not needed here as it's part of the DSN
	// password, _ := parsedURL.User.Password()
	
	// Create pgdriver options
	pgOptions := []pgdriver.Option{
		pgdriver.WithDSN(dsn),
	}

	// Special handling for Neon.tech
	if strings.Contains(host, "neon.tech") {
		// Log masked DSN for debugging
		log.Printf("DEBUG: Detected Neon DB host: %s", host)


		// Critical: Set ServerName for SNI
		// Render's outgoing IP might not match the reverse lookup of the DB, 
		// so Go's strict TLS might fail or send wrong SNI without this.
		pgOptions = append(pgOptions, pgdriver.WithTLSConfig(&tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}))
	} else {
		// Standard TLS skip for other envs if needed (or prefer strict if possible)
		pgOptions = append(pgOptions, pgdriver.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}))
	}
	
	// Use the potentially modified DSN
	pgOptions[0] = pgdriver.WithDSN(parsedURL.String())

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgOptions...))

	DB = bun.NewDB(sqldb, pgdialect.New())

	// Basic check
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Println("Successfully connected to database")
	return nil
}

// CreateSchema creates all tables
func CreateSchema(ctx context.Context) error {
	modelsToRegister := []interface{}{
		(*models.Tournament)(nil),
		(*models.Participant)(nil),
		(*models.Team)(nil),
		(*models.Group)(nil),
		(*models.Match)(nil),
	}

	for _, model := range modelsToRegister {
		_, err := DB.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
