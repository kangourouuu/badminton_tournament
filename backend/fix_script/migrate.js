const { Client } = require("pg");

const runMigration = async () => {
  const connectionString = process.env.DATABASE_URL;
  if (!connectionString) {
    console.error("DATABASE_URL is not set");
    process.exit(1);
  }

  // Auto-fix for corrupted DSNs
  let dsn = connectionString;
  if ((dsn.match(/:\/\//g) || []).length > 1) {
    console.log("Detecting double scheme in DSN, attempting fix...");
    const lastIndex = dsn.lastIndexOf("://");
    const prefixStart = dsn.lastIndexOf("postgres", lastIndex);
    if (prefixStart !== -1) {
      dsn = dsn.substring(prefixStart);
      console.log("Fixed DSN.");
    }
  }

  // Parse URL to modify params
  const url = new URL(dsn);

  // Clean options param if it conflicts with SNI
  if (url.searchParams.has("options")) {
    // The error "Inconsistent project name inferred from SNI" happens when
    // the endpoint ID in options doesn't match the hostname prefix.
    // Easiest fix: remove the options param and let PG infer or just rely on hostname.
    console.log("Removing options param to prevent SNI mismatch...");
    url.searchParams.delete("options");
    dsn = url.toString();
  }

  const client = new Client({
    connectionString: dsn,
    ssl: {
      rejectUnauthorized: false,
    },
  });

  try {
    await client.connect();
    console.log("Connected to database");

    const queries = [
      "ALTER TABLE teams ADD COLUMN IF NOT EXISTS pool TEXT DEFAULT '';",
      "ALTER TABLE groups ADD COLUMN IF NOT EXISTS pool TEXT DEFAULT '';",
      "ALTER TABLE matches ADD COLUMN IF NOT EXISTS sets_detail JSONB;",
    ];

    for (const q of queries) {
      console.log(`Running: ${q}`);
      await client.query(q);
    }

    console.log("Migration completed successfully.");
  } catch (err) {
    console.error("Migration failed:", err);
    process.exit(1);
  } finally {
    await client.end();
  }
};

runMigration();
