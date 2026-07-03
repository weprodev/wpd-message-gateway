CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ==============================================================================
-- SHARED DATABASE TRIGGERS
-- ==============================================================================

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==============================================================================
-- 1. USERS TABLE
-- ==============================================================================

CREATE TABLE users (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name     TEXT        NOT NULL,
    last_name      TEXT        NOT NULL,
    email          TEXT        NOT NULL,
    password_hash  TEXT        NOT NULL,
    email_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_unique      UNIQUE (email),
    CONSTRAINT users_email_lower_check CHECK  (email = lower(email))
);

CREATE TRIGGER trg_users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 2. WORKSPACES TABLE
-- ==============================================================================

CREATE TABLE workspaces (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    slug            TEXT         NOT NULL,
    owner_id        UUID         NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status          TEXT         NOT NULL DEFAULT 'active'
                                 CHECK (status IN ('active', 'suspended')),
    is_private      BOOLEAN      NOT NULL DEFAULT TRUE,
    hashed_pin_code TEXT,
    icon_key        VARCHAR(64),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT workspaces_slug_unique      UNIQUE (slug),
    CONSTRAINT workspaces_slug_lower_check CHECK  (slug = lower(slug))
);

COMMENT ON COLUMN workspaces.icon_key IS 'Optional icon identifier for workspace branding (e.g. marketing, product, support)';

CREATE TRIGGER trg_workspaces_set_updated_at
    BEFORE UPDATE ON workspaces
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_workspaces_owner_id ON workspaces(owner_id);

-- ==============================================================================
-- 3. ROLES TABLE (RBAC)
-- ==============================================================================

CREATE TABLE roles (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT        NOT NULL,
    guard_name   TEXT        NOT NULL DEFAULT 'web',
    display_name TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT roles_name_guard_unique UNIQUE (name, guard_name)
);

CREATE TRIGGER trg_roles_set_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 4. WORKSPACE_MEMBERS TABLE (Relation: users ↔ workspaces)
-- ==============================================================================

CREATE TABLE workspace_members (
    workspace_id UUID        NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id      UUID        NOT NULL REFERENCES users(id)      ON DELETE CASCADE,
    role_id      UUID        NOT NULL REFERENCES roles(id)      ON DELETE RESTRICT,
    joined_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (workspace_id, user_id)
);

CREATE INDEX idx_workspace_members_user_id ON workspace_members(user_id);
CREATE INDEX idx_workspace_members_role_id ON workspace_members(role_id);

-- ==============================================================================
-- 5. INVITATIONS TABLE
-- ==============================================================================

CREATE TABLE invitations (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    email        TEXT         NOT NULL,
    role_id      UUID         NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    token_hash   TEXT         NOT NULL,
    expires_at   TIMESTAMPTZ  NOT NULL,
    status       TEXT         NOT NULL DEFAULT 'pending'
                              CHECK (status IN ('pending', 'accepted', 'revoked')),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT invitations_email_lower_check CHECK (email = lower(email))
);

-- At most one active invite per (workspace, email) to prevent duplicates
CREATE UNIQUE INDEX idx_invitations_pending_per_workspace_email
    ON invitations (workspace_id, email)
    WHERE status = 'pending';

CREATE INDEX idx_invitations_workspace_id ON invitations(workspace_id);
CREATE INDEX idx_invitations_token_hash   ON invitations(token_hash);
CREATE INDEX idx_invitations_role_id       ON invitations(role_id);

-- ==============================================================================
-- 6. WORKSPACE_CHANNELS TABLE
-- ==============================================================================

CREATE TABLE workspace_channels (
    workspace_id UUID    NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    channel_type TEXT    NOT NULL CHECK (channel_type IN ('email', 'sms', 'push', 'chat')),
    enabled      BOOLEAN NOT NULL DEFAULT TRUE,

    PRIMARY KEY (workspace_id, channel_type)
);

-- ==============================================================================
-- 7. PROVIDERS TABLE (Global Platform Catalog)
-- ==============================================================================

CREATE TABLE providers (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT        NOT NULL,
    channel_type TEXT        NOT NULL CHECK (channel_type IN ('email', 'sms', 'push', 'chat')),
    status       TEXT        NOT NULL DEFAULT 'active'
                             CHECK (status IN ('active', 'not_supported', 'blocked')),
    icon_path    TEXT,
    description  TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT providers_name_unique             UNIQUE (name),
    CONSTRAINT providers_id_channel_type_unique UNIQUE (id, channel_type)
);

CREATE TRIGGER trg_providers_set_updated_at
    BEFORE UPDATE ON providers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 8. PROVIDER_CONFIG_FIELDS TABLE (Metadata for provider setups)
-- ==============================================================================

CREATE TABLE provider_config_fields (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_id   UUID        NOT NULL REFERENCES providers(id) ON DELETE RESTRICT,
    key           TEXT        NOT NULL,
    label         TEXT        NOT NULL,
    description   TEXT,
    field_type    TEXT        NOT NULL
                              CHECK (field_type IN ('text', 'password', 'email', 'url', 'boolean', 'textarea', 'select')),
    required      BOOLEAN     NOT NULL DEFAULT FALSE,
    default_value TEXT,
    options       JSONB,      -- for field_type='select' options
    sort_order    SMALLINT    NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT provider_config_fields_provider_key_unique UNIQUE (provider_id, key)
);

CREATE TRIGGER trg_provider_config_fields_set_updated_at
    BEFORE UPDATE ON provider_config_fields
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 9. INTEGRATIONS TABLE (Workspace ↔ Provider Credentials)
-- ==============================================================================

CREATE TABLE integrations (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id     UUID        NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    provider_id      UUID        NOT NULL,
    channel_type     TEXT        NOT NULL CHECK (channel_type IN ('email', 'sms', 'push', 'chat')),
    encrypted_config BYTEA       NOT NULL,
    status           TEXT        NOT NULL DEFAULT 'connected'
                                 CHECK (status IN ('connected', 'error', 'disconnected')),
    is_default       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT integrations_workspace_provider_unique UNIQUE (workspace_id, provider_id),
    CONSTRAINT integrations_provider_channel_fk
        FOREIGN KEY (provider_id, channel_type)
        REFERENCES providers (id, channel_type)
        ON DELETE RESTRICT
);

-- At most one default provider per channel within a workspace
CREATE UNIQUE INDEX idx_integrations_default_per_channel
    ON integrations (workspace_id, channel_type)
    WHERE is_default = TRUE;

CREATE INDEX idx_integrations_provider_id ON integrations(provider_id);
CREATE INDEX idx_integrations_workspace_channel ON integrations(workspace_id, channel_type);

CREATE TRIGGER trg_integrations_set_updated_at
    BEFORE UPDATE ON integrations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 10. API_KEYS TABLE
-- ==============================================================================

CREATE TABLE api_keys (
    id                 UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id       UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name               VARCHAR(255) NOT NULL,
    client_id          VARCHAR(255) NOT NULL,
    client_secret_hash VARCHAR(255) NOT NULL,
    is_active          BOOLEAN      NOT NULL DEFAULT TRUE,
    expires_at         TIMESTAMPTZ,
    last_used_at       TIMESTAMPTZ,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT api_keys_client_id_unique UNIQUE (client_id)
);

CREATE INDEX idx_api_keys_workspace_id ON api_keys(workspace_id);

-- ==============================================================================
-- 11. MESSAGE_REQUEST_LOGS TABLE (Append-only)
-- ==============================================================================

CREATE TABLE message_request_logs (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id  UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    api_key_id    UUID         REFERENCES api_keys(id) ON DELETE SET NULL,
    channel_type  TEXT         NOT NULL CHECK (channel_type IN ('email', 'sms', 'push', 'chat')),
    provider_name VARCHAR(64),
    http_method   VARCHAR(16)  NOT NULL,
    status_code   SMALLINT     NOT NULL,
    endpoint      VARCHAR(512) NOT NULL,
    request_id    VARCHAR(64),
    duration_ms   INT          CHECK (duration_ms >= 0),
    error_message TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_message_request_logs_workspace_created
    ON message_request_logs (workspace_id, created_at DESC);

CREATE INDEX idx_message_request_logs_api_key
    ON message_request_logs (api_key_id)
    WHERE api_key_id IS NOT NULL;

CREATE INDEX idx_message_request_logs_workspace_api_created
    ON message_request_logs (workspace_id, api_key_id, created_at DESC)
    WHERE api_key_id IS NOT NULL;

-- ==============================================================================
-- 11.5. MESSAGE_REQUEST_PAYLOADS TABLE (Append-only body storage)
-- ==============================================================================

CREATE TABLE message_request_payloads (
    log_id        UUID        PRIMARY KEY REFERENCES message_request_logs(id) ON DELETE CASCADE,
    request_body  TEXT,
    response_body TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==============================================================================
-- 12. TEMPLATES TABLE
-- ==============================================================================

CREATE TABLE templates (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    unique_key   VARCHAR(255) NOT NULL,
    channel_type TEXT         NOT NULL DEFAULT 'email'
                              CHECK (channel_type IN ('email', 'sms', 'push', 'chat')),
    subject      VARCHAR(512),
    content      TEXT         NOT NULL,
    category     VARCHAR(64),
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    is_default   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT templates_workspace_unique_key UNIQUE (workspace_id, unique_key)
);

COMMENT ON COLUMN templates.category IS 'e.g. transactional, marketing';

CREATE TRIGGER trg_templates_set_updated_at
    BEFORE UPDATE ON templates
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 13. TEMPLATE_COMPONENTS TABLE
-- ==============================================================================

CREATE TABLE template_components (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    kind         TEXT         NOT NULL CHECK (kind IN ('header', 'footer', 'snippet')),
    name         VARCHAR(255) NOT NULL,
    content      TEXT         NOT NULL,
    metadata     JSONB,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT template_components_workspace_kind_name UNIQUE (workspace_id, kind, name)
);

COMMENT ON COLUMN template_components.metadata IS 'AI prompts, asset references, extra JSON';

CREATE TRIGGER trg_template_components_set_updated_at
    BEFORE UPDATE ON template_components
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 14. WORKSPACE_SETTINGS TABLE (Key-value store)
-- ==============================================================================

CREATE TABLE workspace_settings (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    key          VARCHAR(255) NOT NULL,
    value        TEXT         NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT workspace_settings_workspace_key_unique UNIQUE (workspace_id, key)
);

CREATE TRIGGER trg_workspace_settings_set_updated_at
    BEFORE UPDATE ON workspace_settings
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 15. WORKSPACE_ACCESS_AUDITS TABLE (Append-only)
-- ==============================================================================

CREATE TABLE workspace_access_audits (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id  UUID         NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    actor_user_id UUID         REFERENCES users(id) ON DELETE SET NULL,
    event         VARCHAR(64)  NOT NULL,
    ip_address    INET,
    metadata      JSONB,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workspace_access_audits_workspace_created
    ON workspace_access_audits (workspace_id, created_at DESC);

CREATE INDEX idx_workspace_access_audits_actor_user_id
    ON workspace_access_audits(actor_user_id)
    WHERE actor_user_id IS NOT NULL;

-- ==============================================================================
-- 16. PERMISSIONS TABLE (RBAC Engine)
-- ==============================================================================

CREATE TABLE permissions (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name         TEXT        NOT NULL,
    guard_name   TEXT        NOT NULL DEFAULT 'web',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT permissions_name_guard_unique UNIQUE (name, guard_name)
);

CREATE TRIGGER trg_permissions_set_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==============================================================================
-- 17. ROLE_HAS_PERMISSIONS JOIN TABLE (RBAC Engine)
-- ==============================================================================

CREATE TABLE role_has_permissions (
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,

    PRIMARY KEY (permission_id, role_id)
);

-- ==============================================================================
-- 18. MODEL_HAS_ROLES TABLE (Polymorphic RBAC mappings)
-- ==============================================================================

CREATE TABLE model_has_roles (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id    UUID        NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    model_type TEXT        NOT NULL,
    model_id   UUID        NOT NULL,
    team_id    UUID        REFERENCES workspaces(id) ON DELETE CASCADE,

    CONSTRAINT model_has_roles_unique UNIQUE NULLS NOT DISTINCT (role_id, model_id, model_type, team_id)
);

CREATE INDEX idx_model_has_roles_lookup ON model_has_roles (model_id, model_type, team_id);

-- ==============================================================================
-- 19. MODEL_HAS_PERMISSIONS TABLE (Polymorphic RBAC mappings)
-- ==============================================================================

CREATE TABLE model_has_permissions (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    permission_id UUID        NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    model_type    TEXT        NOT NULL,
    model_id      UUID        NOT NULL,
    team_id       UUID        REFERENCES workspaces(id) ON DELETE CASCADE,

    CONSTRAINT model_has_permissions_unique UNIQUE NULLS NOT DISTINCT (permission_id, model_id, model_type, team_id)
);

CREATE INDEX idx_model_has_permissions_lookup ON model_has_permissions (model_id, model_type, team_id);


