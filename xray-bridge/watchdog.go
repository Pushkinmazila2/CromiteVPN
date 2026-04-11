package xraymobile

import (
    "context"
    "net"
    "time"
)

func watchdog(ctx context.Context) {
    consecutive := 0

    for {
        select {
        case <-ctx.Done():
            return
        case <-time.After(5 * time.Second):
            if pingServer() {
                consecutive = 0
                notifyTunnelStatus(true)
            } else {
                consecutive++
                // 3 consecutive failures = tunnel down
                if consecutive >= 3 {
                    notifyTunnelStatus(false)
                }
            }
        }
    }
}

func pingServer() bool {
    // Проверяем что unix socket живой
    conn, err := net.DialTimeout(
        "unix",
        UnixSocketPath(),
        2*time.Second,
    )
    if err != nil {
        return false
    }
    conn.Close()
    return true
}

func notifyTunnelStatus(up bool) {
    if tunnelUpCb != nil {
        tunnelUpCb(up)
    }
}