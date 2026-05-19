import { Card, Descriptions, Typography } from 'antd'
import { useAuthStore } from '../../store/auth'

export function ProfilePage() {
  const user = useAuthStore((state) => state.user)

  return (
    <div className="admin-page">
      <Typography.Title level={2}>个人资料</Typography.Title>
      <Card>
        <Descriptions bordered column={1}>
          <Descriptions.Item label="用户名">{user?.username}</Descriptions.Item>
          <Descriptions.Item label="昵称">{user?.nickname}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user?.email || '-'}</Descriptions.Item>
          <Descriptions.Item label="角色">{user?.role}</Descriptions.Item>
          <Descriptions.Item label="状态">{user?.status}</Descriptions.Item>
          <Descriptions.Item label="最后登录">{user?.lastLoginAt || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}

