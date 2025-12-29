# DropD - Minimal Text Storage as a Service

## Features

🚀 Go backend (net/http)

🪣 MinIO (S3-compatible) object storage for paste contents

🗄️ PostgreSQL for metadata & indexing

🧱 Repository pattern (storage + DB decoupled)

🆔 Snowflake-style ID generator

- URL-safe (Base62)

⏱️ TTL-based expiration

🧹 Cron job to clean up expired pastes (S3 + PG)

🔁 Idempotent paste creation (optional key support)
