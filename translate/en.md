# Stealth Browser (Cromite-based)

A private Android browser based on **Cromite** with a built-in local proxy chain powered by **Xray-core (VLESS + Reality)**.

This project is designed to route all browser traffic through an internal proxy layer **without utilizing the Android VpnService**.

---

## 🚀 Key Features

*   **Native Xray Integration:** The Xray core is integrated directly into the Chromium network stack via local sockets. All traffic is encapsulated within a single process.
*   **No VpnService Required:** The browser does not create a system-wide VPN interface. There is no "key" icon in the status bar, and the OS or other apps cannot detect proxy usage.
*   **VLESS + REALITY Stealth:** Optimized for REALITY protocol. To external observers (ISPs and DPI), your traffic appears as a standard TLS session to a trusted resource (e.g., Microsoft or Google).
*   **Battery & Performance:** By avoiding packet interception at the TUN interface level, CPU overhead is minimized, significantly extending battery life.
*   **Anti-Leak DNS:** All DNS queries are handled by Xray's built-in DNS module and sent through the encrypted tunnel, preventing ISP-side leaks.
*   **Hidden Traffic:** On-device traffic analyzers (e.g., PCAPdroid) will only show a single connection to your server's IP with a spoofed SNI.

## 🛠 Technical Implementation

This version implements full traffic encapsulation:
1. The `Cronet` network stack redirects requests to the internal Xray instance.
2. **Global Mode** is used: all browser traffic is routed through the configured VLESS outbound by default.
3. Enhanced protection against WebRTC leaks via strict `ice_candidate_policy`.

## 📥 Configuration

Currently, configuration is managed through the browser settings interface:
1. Go to **Settings** -> **Privacy and Security**.
2. Select **Xray Configuration**.
3. Paste your JSON config (VLESS + Reality supported).
4. Restart the app to apply changes.

## 🗺 Roadmap

- [ ] **Dynamic Routing:** Integration of routing rules (Direct / Proxy / Block).
- [ ] **GitHub-Sync:** Automatic update of unlock lists from a remote repository.
- [ ] **DPI Fragment:** Support for packet fragmentation strategies to access resources (YouTube/Telegram) directly without a proxy.
- [ ] **Custom SNI Management:** User-friendly interface for managing decoy domains.

## Credits
* [Cromite Project](https://github.com/uazo/cromite) — The best foundation for private browsing.
* [Project Xray (Xray-core)](https://github.com/xtls/xray-core) — A powerful network engine.

## ⚠️ Disclaimer

This project is intended for research in network security and privacy.
It does not use a system-wide VPN and does not affect other applications on the device.
Users are responsible for complying with their local laws and regulations.

---

## 📜 License

TBD.
