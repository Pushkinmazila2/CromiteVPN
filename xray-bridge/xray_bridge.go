package xraymobile

import (
    "context"
    "fmt"
    "sync"
    "strings"

    core "github.com/xtls/xray-core/core"
    _ "github.com/xtls/xray-core/main/distro/all"
)

var (
    instance   *core.Instance
    mu         sync.Mutex
    cancelWatch context.CancelFunc
    socksPort  int
    socksUser  string
    socksPass  string
    tunnelUpCb func(bool)
)

// StartXray запускает xray с переданным JSON конфигом
func StartXray(configJSON string) error {
    mu.Lock()
    defer mu.Unlock()

    if instance != nil {
        return fmt.Errorf("already running")
    }

    config, err := core.LoadConfig("json", strings.NewReader(configJSON))
    if err != nil {
        return fmt.Errorf("load config: %w", err)
    }

    inst, err := core.New(config)
    if err != nil {
        return fmt.Errorf("create instance: %w", err)
    }

    if err := inst.Start(); err != nil {
        return fmt.Errorf("start: %w", err)
    }

    instance = inst

    // Запускаем watchdog
    ctx, cancel := context.WithCancel(context.Background())
    cancelWatch = cancel
    go watchdog(ctx)

    return nil
}

// StopXray останавливает xray
func StopXray() error {
    mu.Lock()
    defer mu.Unlock()

    if cancelWatch != nil {
        cancelWatch()
        cancelWatch = nil
    }

    if instance == nil {
        return nil
    }

    err := instance.Close()
    instance = nil
    return err
}

// IsRunning возвращает статус
func IsRunning() bool {
    mu.Lock()
    defer mu.Unlock()
    return instance != nil
}

// GetSocksPort возвращает порт SOCKS5
func GetSocksPort() int {
    return socksPort
}

// GetSocksUser возвращает логин для SOCKS5
func GetSocksUser() string {
    return socksUser
}

// GetSocksPass возвращает пароль для SOCKS5
func GetSocksPass() string {
    return socksPass
}

// SetTunnelCallback устанавливает колбэк на изменение статуса туннеля
func SetTunnelCallback(cb func(bool)) {
    tunnelUpCb = cb
}