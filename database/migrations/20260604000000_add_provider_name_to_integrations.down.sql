-- Reverse: remove provider_name column and restore original constraint.

ALTER TABLE integrations
    DROP CONSTRAINT integrations_workspace_channel_provider_unique;

ALTER TABLE integrations
    DROP COLUMN provider_name;

ALTER TABLE integrations
    ADD CONSTRAINT integrations_workspace_provider_unique
    UNIQUE (workspace_id, provider_id);
