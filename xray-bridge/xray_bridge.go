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

// StartXray starts xray with provided JSON config
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

    // Starting watchdog
    ctx, cancel := context.WithCancel(context.Background())
    cancelWatch = cancel
    go watchdog(ctx)

    return nil
}

// StopXray stops xray
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

// IsRunning returns current status
func IsRunning() bool {
    mu.Lock()
    defer mu.Unlock()
    return instance != nil
}

// GetSocksPort returns SOCKS5 port
func GetSocksPort() int {
    return socksPort
}

// GetSocksUser returns SOCKS5 username
func GetSocksUser() string {
    return socksUser
}

// GetSocksPass returns SOCKS5 password
func GetSocksPass() string {
    return socksPass
}

// SetTunnelCallback sets tunnel status callback
func SetTunnelCallback(cb func(bool)) {
    tunnelUpCb = cb
}