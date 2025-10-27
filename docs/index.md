---
layout: home

hero:
  name: "Agent Sandbox"
  text: "为 AI Agent 提供安全的沙箱服务"
  tagline: 在隔离环境中安全地执行 shell 命令和文件操作
  actions:
    - theme: brand
      text: 快速开始
      link: /development
    - theme: alt
      text: GitHub
      link: https://github.com/HJH0924/agent-sandbox

features:
  - title: 🔒 安全沙箱环境
    details: 每个沙箱实例拥有独立的工作空间，使用 API Key 进行安全访问控制
  - title: 📁 文件操作
    details: 支持文件读取、写入和编辑，安全地管理沙箱内的文件
  - title: 🖥️ Shell 命令执行
    details: 支持超时控制的命令执行，捕获 stdout 和 stderr 输出
  - title: 📊 结构化日志
    details: 使用 slog 进行 JSON 格式日志记录，便于观察和调试
  - title: 🔑 API Key 认证
    details: 使用 X-Sandbox-Api-Key 请求头进行安全的访问控制
  - title: ☁️ 云端部署
    details: 支持 E2B 云端沙箱部署，专为 AI Agent 设计
---
