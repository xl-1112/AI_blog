import { useCallback, useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { DeleteOutlined, EditOutlined, EyeOutlined, PlusOutlined, SendOutlined } from '@ant-design/icons'
import { Button, Card, DatePicker, Form, Image, Input, Popconfirm, Select, Space, Table, Tag, Typography, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { adminApi, getReadingMinutes, resolveAssetUrl } from '../../api/client'
import type { Article } from '../../types'

export function ArticlesAdminPage() {
  const [data, setData] = useState<Article[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [form] = Form.useForm()
  const navigate = useNavigate()

  const load = useCallback(async (nextPage = 1) => {
    setLoading(true)
    try {
      const values = form.getFieldsValue()
      const result = await adminApi.articles({
        page: nextPage,
        pageSize: 10,
        keyword: values.keyword,
        tag: values.tag,
        status: values.status,
        dateFrom: values.dateRange?.[0]?.format('YYYY-MM-DD'),
        dateTo: values.dateRange?.[1]?.format('YYYY-MM-DD'),
      })
      setData(result.list)
      setTotal(result.total)
      setPage(nextPage)
    } catch (error) {
      message.error(error instanceof Error ? error.message : '加载失败')
    } finally {
      setLoading(false)
    }
  }, [form])

  useEffect(() => {
    void load(1)
  }, [load])

  const columns: ColumnsType<Article> = [
    {
      title: '封面',
      dataIndex: 'coverUrl',
      width: 92,
      render: (value: string) => (value ? <Image width={64} height={42} src={resolveAssetUrl(value)} className="cover-thumb" /> : '-'),
    },
    { title: '标题', dataIndex: 'title', render: (value, record) => <Link to={`/articles/edit/${record.id}`}>{value}</Link> },
    { title: '摘要', dataIndex: 'summary', ellipsis: true },
    { title: '分类', dataIndex: 'category', width: 110 },
    { title: '阅读时长', width: 100, render: (_, record) => `${getReadingMinutes(record)} 分钟` },
    { title: '浏览量', dataIndex: 'viewCount', width: 90 },
    { title: '发布时间', dataIndex: 'date', width: 120 },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (value: string) => <Tag color={value === 'published' ? 'green' : 'default'}>{value === 'published' ? '已发布' : '草稿'}</Tag>,
    },
    {
      title: '操作',
      width: 210,
      render: (_, record) => (
        <Space>
          <Button icon={<EditOutlined />} size="small" onClick={() => navigate(`/articles/edit/${record.id}`)} />
          <Button icon={<EyeOutlined />} size="small" href={`#/site/articles/${record.slug || record.id}`} target="_blank" />
          {record.status !== 'published' ? (
            <Button
              icon={<SendOutlined />}
              size="small"
              onClick={async () => {
                await adminApi.publishArticle(record.id)
                message.success('文章已发布')
                void load(page)
              }}
            />
          ) : null}
          <Popconfirm title="确定删除这篇文章吗？" onConfirm={async () => {
            await adminApi.deleteArticle(record.id)
            message.success('已删除')
            void load(page)
          }}>
            <Button danger icon={<DeleteOutlined />} size="small" />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div className="admin-page">
      <div className="page-heading">
        <Typography.Title level={2}>内容管理</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/articles/create')}>
          新增文章
        </Button>
      </div>
      <Card className="mb-16">
        <Form form={form} layout="inline" onFinish={() => load(1)}>
          <Form.Item name="keyword">
            <Input.Search placeholder="标题搜索" allowClear onSearch={() => load(1)} />
          </Form.Item>
          <Form.Item name="tag">
            <Input placeholder="分类标签筛选" allowClear />
          </Form.Item>
          <Form.Item name="status">
            <Select
              placeholder="状态"
              allowClear
              style={{ width: 120 }}
              options={[
                { label: '草稿', value: 'draft' },
                { label: '已发布', value: 'published' },
              ]}
            />
          </Form.Item>
          <Form.Item name="dateRange">
            <DatePicker.RangePicker />
          </Form.Item>
          <Button htmlType="submit">查询</Button>
        </Form>
      </Card>
      <Card>
        <Table<Article>
          rowKey="id"
          loading={loading}
          columns={columns}
          dataSource={data}
          pagination={{ current: page, pageSize: 10, total, onChange: (next) => load(next) }}
        />
      </Card>
    </div>
  )
}
