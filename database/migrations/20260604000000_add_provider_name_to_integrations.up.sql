-- Add provider_name to integrations (aligns Go domain with DB schema).
-- The Go code reads/writes provider_name; the original migration only had provider_id.

ALTER TABLE integrations
    ADD COLUMN provider_name TEXT NOT NULL DEFAULT '';

-- Backfill from the providers table for any existing rows.
UPDATE integrations i
SET provider_name = p.name
FROM providers p
WHERE i.provider_id = p.id
  AND i.provider_name = '';

-- Replace unique constraint: one integration per (workspace, channel, provider).
ALTER TABLE integrations
    DROP CONSTRAINT integrations_workspace_provider_unique;

ALTER TABLE integrations
    ADD CONSTRAINT integrations_workspace_channel_provider_unique
    UNIQUE (workspace_id, channel_type, provider_name);
