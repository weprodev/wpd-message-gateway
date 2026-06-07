-- Seed config fields for Mailgun provider
INSERT INTO provider_config_fields (id, provider_id, key, label, description, field_type, required, default_value, sort_order)
VALUES
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000031', 'api_key', 'API Key', 'Your Mailgun Private API Key (starts with key-)', 'password', TRUE, NULL, 10),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000031', 'domain', 'Domain', 'Your Mailgun sending domain', 'text', TRUE, NULL, 20),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000031', 'base_url', 'Base URL', 'Mailgun API base URL (e.g. https://api.mailgun.net/v3 or https://api.eu.mailgun.net/v3)', 'url', FALSE, 'https://api.mailgun.net/v3', 30),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000031', 'from_email', 'From Email', 'Default sender email address (e.g. noreply@yourdomain.com)', 'email', TRUE, NULL, 40),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000031', 'from_name', 'From Name', 'Default sender display name (e.g. Support Team)', 'text', FALSE, NULL, 50)
ON CONFLICT (provider_id, key) DO NOTHING;
