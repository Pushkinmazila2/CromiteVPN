# Stealth Browser (基于 Cromite)

一款基于 **Cromite** 的安卓私密浏览器，内置通过 **Xray-core (VLESS + Reality)** 实现的本地代理链。

本项目旨在让浏览器所有流量仅通过本地代理层传输，**无需调用 Android VpnService**。

---

## 🚀 核心特性

*   **原生 Xray 集成：** Xray 核心通过本地套接字直接集成到 Chromium 网络栈中。所有流量均在单一进程内封装。
*   **无需 VpnService：** 浏览器不创建系统级 VPN 接口。状态栏不会出现“钥匙”图标，系统或其他应用无法检测到代理的使用。
*   **VLESS + REALITY 隐匿传输：** 默认支持 REALITY 协议。对于外部观察者（运营商或 DPI）而言，您的流量看起来就像是去往信任资源（如 Microsoft 或 Google）的标准 TLS 会话。
*   **省电与高性能：** 避免了 TUN 接口层面的数据包拦截，降低了 CPU 负载，显著延长电池续航。
*   **防泄漏 DNS：** 所有 DNS 查询由 Xray 内置 DNS 模块处理并发送至加密隧道，彻底杜绝运营商侧的 DNS 泄漏。
*   **隐匿流量：** 设备端流量分析工具（如 PCAPdroid）只能看到指向您服务器 IP 且带有伪装 SNI 的单一连接。

## 🛠 技术实现

当前版本实现了完整的流量封装：
1. `Cronet` 网络栈将请求重定向至内置的 Xray 实例。
2. 使用 **全局模式 (Global Mode)**：浏览器所有流量默认通过配置的 VLESS 出站（Outbound）进行传输。
3. 通过严格的 `ice_candidate_policy` 策略防止 WebRTC 泄漏。

## 📥 配置说明

当前阶段可通过浏览器设置界面进行配置：
1. 打开 **Settings (设置)** -> **Privacy and Security (隐私与安全)**。
2. 选择 **Xray Configuration**。
3. 粘贴您的 JSON 配置（支持 VLESS + Reality）。
4. 重启应用以应用设置。

## 🗺 路线图 (Roadmap)

- [ ] **动态路由 (Dynamic Routing)：** 集成路由规则（直连 / 代理 / 屏蔽）。
- [ ] **GitHub 同步：** 自动从远程仓库更新分流列表。
- [ ] **DPI 分段 (Fragment)：** 支持数据包分段策略，无需代理即可直连访问特定资源（如 YouTube/Telegram）。
- [ ] **自定义 SNI 管理：** 便于更换伪装域名的用户界面。

## 致谢 (Credits)
* [Cromite Project](https://github.com/uazo/cromite) — 为私密浏览提供了最佳基础。
* [Project Xray (Xray-core)](https://github.com/xtls/xray-core) — 强大的网络引擎。

## ⚠️ 免责声明

本项目仅用于网络安全和隐私领域的科学研究。
它不使用系统级 VPN，不会影响设备上的其他应用。
用户需自行承担遵守当地法律法规的责任。

---

## 📜 许可证 (License)

待定。