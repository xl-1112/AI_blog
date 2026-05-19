import { useEffect, useState } from 'react'
import { Button, Card, Form, Input, Modal, Popconfirm, Select, Space, Table, Tabs, Tag, Typography, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { adminApi } from '../../api/client'
import type { LoginLog, Role, User, UserStatus } from '../../types'

export function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [logs, setLogs] = useState<LoginLog[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)
  const [open, setOpen] = useState(false)
  const [form] = Form.useForm()

  async function load() {
    setLoading(true)
    try {
      const [userResult, logResult] = await Promise.all([adminApi.users({ page: 1, pageSize: 50 }), adminApi.loginLogs({ page: 1, pageSize: 50 })])
      setUsers(userResult.list)
      setTotal(userResult.total)
      setLogs(logResult.list)
    } catch (error) {
      message.error(error instanceof Error ? error.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [])

  const columns: ColumnsType<User> = [
    { title: '用户名', dataIndex: 'username' },
    { title: '昵称', dataIndex: 'nickname' },
    { title: '邮箱', dataIndex: 'email' },
    { title: '角色', dataIndex: 'role', render: (value: Role) => roleText(value) },
    { title: '状态', dataIndex: 'status', render: (value: UserStatus) => <Tag color={value === 'active' ? 'green' : 'red'}>{value === 'active' ? '启用' : '禁用'}</Tag> },
    {
      title: '操作',
      render: (_, record) => (
        <Space>
          <Button size="small" onClick={() => {
            setEditing(record)
            form.setFieldsValue(record)
            setOpen(true)
          }}>编辑</Button>
          <Button size="small" onClick={() => resetPassword(record)}>重置密码</Button>
          <Button size="small" onClick={() => adminApi.setUserStatus(record.id, record.status === 'active' ? 'disabled' : 'active').then(load)}>
            {record.status === 'active' ? '禁用' : '启用'}
          </Button>
          <Popconfirm title="确定删除用户吗？" onConfirm={() => adminApi.deleteUser(record.id).then(load)}>
            <Button danger size="small">删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  async function saveUser() {
    const values = await form.validateFields()
    if (editing) {
      await adminApi.updateUser(editing.id, values)
    } else {
      await adminApi.createUser(values)
    }
    message.success('用户已保存')
    setOpen(false)
    setEditing(null)
    form.resetFields()
    await load()
  }

  function resetPassword(user: User) {
    Modal.confirm({
      title: `重置 ${user.username} 的密码`,
      content: <Input.Password id="reset-password" placeholder="请输入新密码" />,
      onOk: async () => {
        const input = document.getElementById('reset-password') as HTMLInputElement | null
        await adminApi.resetPassword(user.id, input?.value ?? '')
        message.success('密码已重置')
      },
    })
  }

  return (
    <div className="admin-page">
      <div className="page-heading">
        <Typography.Title level={2}>用户管理</Typography.Title>
        <Button type="primary" onClick={() => {
          setEditing(null)
          form.resetFields()
          form.setFieldsValue({ role: 'editor', status: 'active' })
          setOpen(true)
        }}>新增用户</Button>
      </div>
      <Tabs
        items={[
          {
            key: 'users',
            label: '用户列表',
            children: (
              <Card>
                <Table<User> rowKey="id" loading={loading} dataSource={users} columns={columns} pagination={{ total, pageSize: 50 }} />
              </Card>
            ),
          },
          {
            key: 'logs',
            label: '登录日志',
            children: (
              <Card>
                <Table<LoginLog>
                  rowKey="id"
                  dataSource={logs}
                  columns={[
                    { title: '用户名', dataIndex: 'username' },
                    { title: 'IP', dataIndex: 'ip' },
                    { title: '结果', dataIndex: 'success', render: (value: boolean) => value ? <Tag color="green">成功</Tag> : <Tag color="red">失败</Tag> },
                    { title: '原因', dataIndex: 'reason' },
                    { title: '时间', dataIndex: 'createdAt' },
                  ]}
                />
              </Card>
            ),
          },
        ]}
      />
      <Modal title={editing ? '编辑用户' : '新增用户'} open={open} onCancel={() => setOpen(false)} onOk={saveUser} destroyOnClose>
        <Form form={form} layout="vertical">
          {!editing ? (
            <>
              <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
                <Input />
              </Form.Item>
              <Form.Item name="password" label="初始密码" rules={[{ required: true }]}>
                <Input.Password />
              </Form.Item>
            </>
          ) : null}
          <Form.Item name="nickname" label="昵称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="role" label="角色" rules={[{ required: true }]}>
            <Select options={[
              { label: '超级管理员', value: 'super_admin' },
              { label: '管理员', value: 'admin' },
              { label: '编辑', value: 'editor' },
            ]} />
          </Form.Item>
          <Form.Item name="status" label="状态" rules={[{ required: true }]}>
            <Select options={[{ label: '启用', value: 'active' }, { label: '禁用', value: 'disabled' }]} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

function roleText(role: Role) {
  if (role === 'super_admin') return '超级管理员'
  if (role === 'admin') return '管理员'
  return '编辑'
}

