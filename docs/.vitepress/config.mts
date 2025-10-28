import { defineConfig } from 'vitepress'

export default defineConfig({
  base: '/agent-sandbox/',
  title: "Agent Sandbox Service",
  description: "一个为 AI Agent 设计的安全沙箱服务，提供隔离的文件操作和命令执行环境。基于 Go 和 gRPC/Connect-RPC 构建。",
  themeConfig: {
    nav: [
      { text: '首页', link: '/' },
      { text: '开发指南', link: '/development' },
      { text: 'GitHub', link: 'https://github.com/HJH0924/agent-sandbox' }
    ],

    sidebar: [
      {
        text: '指南',
        items: [
          { text: '开发指南', link: '/development' }
        ]
      },
      {
        text: '服务',
        items: [
          { text: '核心服务', link: '/core/index' },
          { text: '文件服务', link: '/file/index' },
          { text: 'Shell 服务', link: '/shell/index' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/HJH0924/agent-sandbox' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2025 Kaho'
    }
  }
})
