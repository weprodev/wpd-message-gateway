package kavenegar

import (
    "github.com/weprodev/wpd-message-gateway/internal/app/registry"
    "github.com/weprodev/wpd-message-gateway/internal/core/port"
)



func init() {
    registry.RegisterSMSProvider("kavenegar", func(cfg registry.SMSConfig) (port.SMSSender, error) {
        return New(Config{
            APIKey:    cfg.APIKey,
            FromPhone:  cfg.FromPhone,
        })
    })
}
