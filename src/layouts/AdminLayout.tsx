import { useEffect, type ReactNode } from 'react'
import { Link, Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  BarChartOutlined,
  DashboardOutlined,
  FileTextOutlined,
  HomeOutlined,
  LogoutOutlined,
  ReloadOutlined,
  SettingOutlined,
  TagsOutlined,
  TeamOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { Avatar, Button, Layout, Menu, Space, Typography, message } from 'antd'
import type { ItemType } from 'antd/es/menu/interface'
import { adminApi } from '../api/client'
import { useAuthStore } from '../store/auth'
import type { Role } from '../types'

const { Header, Sider, Content } = Layout

const menuItems: Array<{ key: string; label: string; icon: ReactNode; roles: Role[] }> = [
  { key: '/dashboard', label: '仪表盘', icon: <DashboardOutlined />, roles: ['super_admin', 'admin', 'editor'] },
  { key: '/articles', label: '内容管理', icon: <FileTextOutlined />, roles: ['super_admin', 'admin', 'editor'] },
  { key: '/home-config', label: '页面配置', icon: <HomeOutlined />, roles: ['super_admin', 'admin'] },
  { key: '/about-config', label: '关于页配置', icon: <UserOutlined />, roles: ['super_admin', 'admin'] },
  { key: '/skills', label: '技术栈', icon: <TagsOutlined />, roles: ['super_admin', 'admin'] },
  { key: '/timeline', label: '工作经历', icon: <BarChartOutlined />, roles: ['super_admin', 'admin'] },
  { key: '/tags', label: '分类标签', icon: <TagsOutlined />, roles: ['super_admin', 'admin', 'editor'] },
  { key: '/site-config', label: '网站配置', icon: <SettingOutlined />, roles: ['super_admin', 'admin'] },
  { key: '/analytics', label: '数据统计', icon: <BarChartOutlined />, roles: ['super_admin', 'admin', 'editor'] },
  { key: '/users', label: '用户管理', icon: <TeamOutlined />, roles: ['super_admin'] },
  { key: '/profile', label: '个人资料', icon: <UserOutlined />, roles: ['super_admin', 'admin', 'editor'] },
]

export function AdminLayout() {
  const { token, user, logout, setUser } = useAuthStore()
  const location = useLocation()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) return
    adminApi.me().then(setUser).catch(() => undefined)
  }, [setUser, token])

  if (!token) {
    return <Navigate to="/login" replace />
  }

  const visibleItems: ItemType[] = menuItems
    .filter((item) => user?.role && item.roles.includes(user.role))
    .map((item) => ({
      key: item.key,
      icon: item.icon,
      label: <Link to={item.key}>{item.label}</Link>,
    }))

  return (
    <Layout className="admin-shell">
      <Sider breakpoint="lg" collapsedWidth={0} width={228} className="admin-sider">
        <Link to="/dashboard" className="admin-brand">
          <span className="admin-brand-mark">L</span>
          <span>
            <strong>Liang CMS</strong>
            <small>博客后台管理</small>
          </span>
        </Link>
        <Menu mode="inline" selectedKeys={[bestMenuKey(location.pathname)]} items={visibleItems} className="admin-menu" />
      </Sider>
      <Layout>
        <Header className="admin-header">
          <Typography.Text type="secondary">当前管理员</Typography.Text>
          <Space size={12}>
            <Avatar icon={<UserOutlined />} />
            <Typography.Text strong>{user?.nickname || user?.username}</Typography.Text>
            <Button icon={<ReloadOutlined />} onClick={() => window.location.reload()}>
              刷新
            </Button>
            <Button
              icon={<LogoutOutlined />}
              onClick={() => {
                logout()
                message.success('已退出登录')
                navigate('/login', { replace: true })
              }}
            >
              退出登录
            </Button>
          </Space>
        </Header>
        <Content className="admin-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

function bestMenuKey(pathname: string) {
  if (pathname.startsWith('/articles')) return '/articles'
  return menuItems.find((item) => pathname.startsWith(item.key))?.key ?? '/dashboard'
}
