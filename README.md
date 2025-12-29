# Dropd
Dropd - Text Storage as a service

Skip to content
You said:
names for a pastebin like service
ChatGPT said:

Alright, pastebin names — I’ll give you good ones, not the cringe startup-generator shit. Grouped by vibe 👇
Clean / Minimal (devs love this)

    Pastey

    Snip

    Clip

    Binx

    Pastr

    Textbin

    DropTXT

    PasteUp

Dev / Hacker Vibe

    Hexbin

    BytePaste

    StackBin

    CodeDrop

    SnipBin

    RawPaste

    LogDump

    BitBin

Fast / Ephemeral (expires, throwaway pastes)

    BurnBin

    GhostPaste

    Vanish

    FlashPaste

    TempTXT

    GoneBin

    Blink

Slightly Edgy / Fun

    Paste.sh (classic, very OSS)

    YeetBin

    DumpIt

    Spill

    Leakpad

    DropIt

    PasteFck (lol, only if you’re bold)

Serious / OSS / Infra-grade

    OpenPaste

    Plaintext

    TextDepot

    PasteHub

    DataSlab

    ClipBoard.io (if domain free)

My top picks (no bullshit)

If you want:

    Minimal & respectable → Snip, Clip, Pastr

    Infra / OSS flex → Paste.sh, OpenPaste

    Ephemeral pastebin → BurnBin, GhostPaste

If you want, tell me:

    public vs private

    expiring-only or permanent

    dev-only vs normies

    domain constraints (.sh, .dev, .io, .in)

and I’ll narrow it to 3 perfect names. fck yea.
You said:
its in Go
ChatGPT said:

Cool, Go-native changes the vibe a bit — you can lean into Go / infra / boring-but-solid energy (which is good).

Here are names that fit a Go service and won’t feel out of place next to kubectl, etcd, minio, etc.
Go-coded / Gopher-friendly 🐹

Subtle, not cringe:

    Gopaste (obvious, but clean)

    GoBin

    GoSnip

    GoClip

    PasteGo (meh, but fine)

    GopherBin (only if you’re okay being a bit loud)

👉 My pick here: GoBin or GoSnip
Infra-grade (sounds like a CNCF project)

These feel correct in Go repos:

    Snipd

    Clipd

    Pasted

    Textd

    Bindump

    Logbin

    Dropd

    -d suffix = daemon energy. Very Go. Very based.

🔥 Snipd and Pasted go hard.
Minimal / Unix-y (perfect for curl | pipe)

These are chef’s kiss for Go CLIs & HTTP services:

    paste.sh

    snip.sh

    bin.sh

    pb (pastebin shorthand)

    txt

    drop

Example UX:

curl -F file=@log.txt https://paste.sh

This shit feels right in Go land.
Ephemeral-first (if you support TTLs)

Fits well with your S3 + PG + cleanup jobs:

    Burn

    Ash

    Ember

    Ghost

    Vanish

    Fade

🔥 Ash or Ember + Go backend = clean af.
Names I’d actually ship (top tier)

If I were building this in Go:

    Snipd – infra daemon, clean repo name

    GoBin – simple, discoverable

    paste.sh – OSS classic, curl-friendly

    Ash – perfect for expiring pastes

    Pasted – boring, solid, Go-core vibes

Quick sanity check (sparring partner mode)

    ❌ Over-Go-ifying (GoGoPaste) = cringe

    ❌ Cute names = age badly

    ✅ Short, CLI-friendly names = win

    ✅ Unix / daemon vibes = timeless

Tell me:

    repo name or product name?

    CLI planned or only HTTP?

    expiring by default?

I’ll lock one final name and even suggest:

    binary name

    module path

    package names

No fluff.
You said:
gimme a readme,
its in go
minio
pg
repo pattern
snofalke style key gen
cleanup cronjob for expired pastes
ChatGPT said:

Alright, here’s a clean, infra-grade README you can literally drop into the repo.
No marketing fluff, Go-first, boring-in-a-good-way.
📝 Paste Service

A minimal Pastebin-like service written in Go, backed by MinIO (S3) and PostgreSQL, designed with a clean repository pattern, Snowflake-style ID generation, and a cron-based cleanup job for expired pastes.

Built to be fast, boring, and production-friendly.
Features

    🚀 Go backend (net/http or chi/gin friendly)

    🪣 MinIO (S3-compatible) object storage for paste contents

    🗄️ PostgreSQL for metadata & indexing

    🧱 Repository pattern (storage + DB decoupled)

    🆔 Snowflake-style ID generator

        sortable

        unique

        URL-safe (Base62)

    ⏱️ TTL-based expiration

    🧹 Cron job to clean up expired pastes (S3 + PG)

    🔁 Idempotent paste creation (optional key support)

Architecture Overview

┌──────────┐     ┌────────────┐
│  Client  │ ──▶ │   HTTP API │
└──────────┘     └─────┬──────┘
                        │
            ┌───────────┴───────────┐
            │       Service Layer     │
            │  (business logic)       │
            └───────────┬───────────┘
                        │
        ┌───────────────┴───────────────┐
        │        Repository Layer        │
        │                               │
┌───────┴─────────┐          ┌──────────┴─────────┐
│ PostgreSQL Repo │          │   MinIO Repo        │
│ (metadata)     │          │ (paste content)     │
└────────────────┘          └─────────────────────┘

Tech Stack

    Language: Go

    Object Storage: MinIO (S3-compatible)

    Database: PostgreSQL

    ID Generation: Snowflake-inspired (timestamp + node + sequence)

    Cleanup: Cron-based background job

Data Model
paste_metadata (PostgreSQL)
Column	Type	Notes
id	TEXT	Base62 Snowflake ID
idempotency_key	TEXT	Optional
filename	TEXT	Original filename
created_at	TIMESTAMPTZ	Go time.Time
expires_at	TIMESTAMPTZ	Nullable
content_type	TEXT	Optional

Indexes:

    PRIMARY KEY (id)

    INDEX (expires_at)

    UNIQUE (idempotency_key) (optional)

ID Generation

IDs are generated using a Snowflake-style algorithm:

| timestamp | node_id | sequence |

Properties:

    Monotonically sortable

    Globally unique

    Encoded as Base62 for short URLs

Example:

https://paste.example.com/4fZk9aQ

Repository Pattern

Interfaces live in internal/repository:

type PasteRepository interface {
	Create(ctx context.Context, p *Paste) error
	Get(ctx context.Context, id string) (*Paste, error)
	Delete(ctx context.Context, id string) error
	ListExpired(ctx context.Context, before time.Time) ([]Paste, error)
}

Storage abstraction (MinIO):

type ObjectStorage interface {
	Put(ctx context.Context, key string, data []byte) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

Service layer orchestrates DB + S3 operations.
Cleanup Job (Expired Pastes)

A cron-triggered job runs periodically to:

    Query expired pastes from PostgreSQL

    Delete content from MinIO

    Delete metadata from PostgreSQL

┌──────── Cron ────────┐
│ SELECT expired IDs   │
│ DELETE from MinIO    │
│ DELETE from PG       │
└──────────────────────┘

Safe to retry. Failures won’t corrupt state.
Running Locally
Requirements

    Go 1.21+

    Docker / Docker Compose

Start dependencies

docker compose up -d

This starts:

    MinIO

    PostgreSQL

Run the service

go run ./cmd/api

(Optional cleanup worker)

go run ./cmd/cleanup

Example API Usage
Create Paste


ChatGPT can make mistakes. Check important info. See Cookie Preferences.
ChatGPT is still generating a response...
