import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "Agent Sandbox",
  description: "为 AI Agent 提供安全的代码执行和文件操作能力的沙箱服务",
  themeConfig: {
    nav: [
      { text: '首页', link: '/' },
      { text: '开发指南', link: '/development' },
      { text: 'GitHub', link: 'https://github.com/HJH0924/agent-sandbox' }
    ],

    sidebar: [
      {
        text: '介绍',
        items: [
          { text: '快速开始', link: '/index' },
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
      copyright: 'Copyright © 2024'
    }
  }
})
