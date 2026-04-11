package xraymobile

import (
    "crypto/rand"
    "encoding/json"
    "math/big"
)



// Profile - connection profile
type Profile struct {
    Type      string // vless, vmess, trojan, shadowsocks
    Address   string
    Port      int
    UUID      string // for vless/vmess/trojan
    Password  string // for shadowsocks/trojan
    Transport string // tcp, ws, grpc, reality
    TLS       TLSConfig
    Name      string
}

type TLSConfig struct {
    Enabled    bool
    SNI        string
    Fingerprint string
    // Reality specific
    PublicKey  string
    ShortID    string
}


var appPackageName = "com.example.browser"

func AppPackageName() string {
    return appPackageName
}

// BuildConfig generates xray JSON config
func BuildConfig(profile Profile) (string, error) {
    port, err := randomPort()
    if err != nil {
        return "", err
    }
    socksPort = port
    socksUser = randomString(16)
    socksPass = randomString(32)

    config := map[string]interface{}{
        "log": map[string]interface{}{
            "loglevel": "warning",
        },
        "inbounds": []interface{}{
            map[string]interface{}{
                "tag":      "socks-in",
                "protocol": "socks",
                // Unix socket вместо TCP порта
                "listen": "unix:" + UnixSocketPath(),
                "settings": map[string]interface{}{
                    "auth": "password",
                    "accounts": []interface{}{
                        map[string]interface{}{
                            "user": socksUser,
                            "pass": socksPass,
                        },
                    },
                    "udp": false,
                },
            },
        },
        "outbounds": []interface{}{
            buildOutbound(profile),
            map[string]interface{}{
                "tag":      "direct",
                "protocol": "freedom",
            },
            map[string]interface{}{
                "tag":      "block",
                "protocol": "blackhole",
            },
        },
        "routing": map[string]interface{}{
            "domainStrategy": "IPIfNonMatch",
            "rules": []interface{}{
                map[string]interface{}{
                    "type":        "field",
                    "outboundTag": "proxy",
                    "network":     "tcp,udp",
                },
            },
        },
        "stats":  map[string]interface{}{},
        "policy": buildPolicy(),
    }

    data, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func buildOutbound(p Profile) map[string]interface{} {
    switch p.Type {
    case "vless":
        return buildVless(p)
    case "vmess":
        return buildVmess(p)
    case "trojan":
        return buildTrojan(p)
    case "shadowsocks":
        return buildShadowsocks(p)
    default:
        return buildVless(p)
    }
}

func buildVless(p Profile) map[string]interface{} {
    streamSettings := buildStreamSettings(p)

    return map[string]interface{}{
        "tag":      "proxy",
        "protocol": "vless",
        "settings": map[string]interface{}{
            "vnext": []interface{}{
                map[string]interface{}{
                    "address": p.Address,
                    "port":    p.Port,
                    "users": []interface{}{
                        map[string]interface{}{
                            "id":         p.UUID,
                            "encryption": "none",
                            "flow":       "xtls-rprx-vision",
                        },
                    },
                },
            },
        },
        "streamSettings": streamSettings,
    }
}

func buildStreamSettings(p Profile) map[string]interface{} {
    settings := map[string]interface{}{
        "network": p.Transport,
    }

    if p.TLS.Enabled {
        if p.Transport == "reality" {
            settings["security"] = "reality"
            settings["realitySettings"] = map[string]interface{}{
                "serverName":  p.TLS.SNI,
                "fingerprint": p.TLS.Fingerprint,
                "publicKey":   p.TLS.PublicKey,
                "shortId":     p.TLS.ShortID,
            }
            settings["network"] = "tcp"
        } else {
            settings["security"] = "tls"
            settings["tlsSettings"] = map[string]interface{}{
                "serverName":         p.TLS.SNI,
                "fingerprint":        p.TLS.Fingerprint,
                "allowInsecure":      false,
            }
        }
    }

    return settings
}

func buildPolicy() map[string]interface{} {
    return map[string]interface{}{
        "system": map[string]interface{}{
            "statsInboundUplink":   true,
            "statsInboundDownlink": true,
        },
    }
}

// UnixSocketPath returns unix socket path
func UnixSocketPath() string {
    return "/data/data/" + AppPackageName() + "/files/xray.sock"
}

func randomPort() (int, error) {
    n, err := rand.Int(rand.Reader, big.NewInt(16383))
    if err != nil {
        return 0, err
    }
    return int(n.Int64()) + 49152, nil
}

func randomString(n int) string {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    result := make([]byte, n)
    for i := range result {
        num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
        result[i] = chars[num.Int64()]
    }
    return string(result)
}

func buildVmess(p Profile) map[string]interface{} {
    return map[string]interface{}{
        "tag":      "proxy",
        "protocol": "vmess",
        "settings": map[string]interface{}{
            "vnext": []interface{}{
                map[string]interface{}{
                    "address": p.Address,
                    "port":    p.Port,
                    "users": []interface{}{
                        map[string]interface{}{
                            "id":       p.UUID,
                            "alterId":  0,
                            "security": "auto",
                        },
                    },
                },
            },
        },
        "streamSettings": buildStreamSettings(p),
    }
}

func buildTrojan(p Profile) map[string]interface{} {
    return map[string]interface{}{
        "tag":      "proxy",
        "protocol": "trojan",
        "settings": map[string]interface{}{
            "servers": []interface{}{
                map[string]interface{}{
                    "address":  p.Address,
                    "port":     p.Port,
                    "password": p.Password,
                },
            },
        },
        "streamSettings": buildStreamSettings(p),
    }
}

func buildShadowsocks(p Profile) map[string]interface{} {
    return map[string]interface{}{
        "tag":      "proxy",
        "protocol": "shadowsocks",
        "settings": map[string]interface{}{
            "servers": []interface{}{
                map[string]interface{}{
                    "address":  p.Address,
                    "port":     p.Port,
                    "method":   p.TLS.SNI,
                    "password": p.Password,
                },
            },
        },
    }
}