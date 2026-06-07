-- Seed default integration providers
INSERT INTO providers (id, name, channel_type, status, icon_path, description)
VALUES
    ('00000000-0000-0000-0000-000000000030', 'memory', 'email', 'active', '/assets/providers/memory.svg', 'In-memory simulator for development and simulation'),
    ('00000000-0000-0000-0000-000000000031', 'mailgun', 'email', 'active', '/assets/providers/mailgun.svg', 'Reliable transactional and marketing email delivery service'),
    ('00000000-0000-0000-0000-000000000032', 'cm.com', 'sms', 'not_supported', '/assets/providers/cm.svg', 'SMS messaging and mobile communication via CM.com (coming soon)'),
    ('00000000-0000-0000-0000-000000000033', 'telegram', 'chat', 'not_supported', '/assets/providers/telegram.svg', 'Send messages via Telegram Bot API (coming soon)'),
    ('00000000-0000-0000-0000-000000000034', 'whatsapp', 'chat', 'not_supported', '/assets/providers/whatsapp.svg', 'Send templates and messages via WhatsApp Business Platform (coming soon)'),
    ('00000000-0000-0000-0000-000000000035', 'twilio', 'sms', 'not_supported', '/assets/providers/twilio.svg', 'SMS and voice communication APIs via Twilio (coming soon)')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    channel_type = EXCLUDED.channel_type,
    icon_path = EXCLUDED.icon_path,
    description = EXCLUDED.description,
    status = EXCLUDED.status;
