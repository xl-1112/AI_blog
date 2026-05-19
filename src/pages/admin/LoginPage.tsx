import { Navigate, useNavigate } from 'react-router-dom'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
import { Button, Card, Form, Input, Typography, message } from 'antd'
import { adminApi } from '../../api/client'
import { useAuthStore } from '../../store/auth'

export function LoginPage() {
  const navigate = useNavigate()
  const { token, setSession } = useAuthStore()

  if (token) return <Navigate to="/dashboard" replace />

  return (
    <main className="login-page">
      <Card className="login-card">
        <Typography.Title level={2}>Liang 博客后台</Typography.Title>
        <Typography.Paragraph type="secondary">使用数据库中的管理员账号登录。</Typography.Paragraph>
        <Form
          layout="vertical"
          initialValues={{ username: 'admin' }}
          onFinish={async (values: { username: string; password: string }) => {
            const hide = message.loading('正在登录...')
            try {
              const result = await adminApi.login(values)
              setSession(result.token, result.userInfo)
              message.success('登录成功')
              navigate('/dashboard', { replace: true })
            } catch (error) {
              message.error(error instanceof Error ? error.message : '登录失败')
            } finally {
              hide()
            }
          }}
        >
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input size="large" prefix={<UserOutlined />} autoComplete="username" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password size="large" prefix={<LockOutlined />} autoComplete="current-password" />
          </Form.Item>
          <Button type="primary" size="large" htmlType="submit" block>
            登录
          </Button>
        </Form>
      </Card>
    </main>
  )
}

