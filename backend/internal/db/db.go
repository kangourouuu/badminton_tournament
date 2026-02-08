package db

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	// Debug log: Print DSN masked
	maskedDSN := dsn
	if parts := strings.Split(dsn, "@"); len(parts) > 1 {
		maskedDSN = "postgres://****@" + parts[1]
	}
	fmt.Printf("DEBUG: Init DB connection with DSN: %s\n", maskedDSN)

	// Auto-fix for Neon.tech SNI if options are missing
	if strings.Contains(dsn, "neon.tech") && !strings.Contains(dsn, "endpoint=") {
		parts := strings.Split(dsn, "@")
		if len(parts) > 1 {
			hostParts := strings.Split(parts[1], ".")
			if len(hostParts) > 0 {
				endpointID := hostParts[0]
				endpointID = strings.TrimSuffix(endpointID, "-pooler")

				if strings.Contains(dsn, "?") {
					dsn += "&options=endpoint%3D" + endpointID
				} else {
					dsn += "?options=endpoint%3D" + endpointID
				}
				fmt.Printf("DEBUG: Auto-patched Neon DSN with endpoint ID: %s\n", endpointID)
			}
		}
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	))
	DB = bun.NewDB(sqldb, pgdialect.New())

	// Basic check
	if err := DB.Ping(); err != nil {
		return err
	}

	return nil
}
