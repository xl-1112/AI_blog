import { useEffect, useState } from 'react'
import { Line } from '@ant-design/plots'
import { Card, Col, Row, Statistic, Table, Typography, message } from 'antd'
import { adminApi } from '../../api/client'
import type { ArticleSummary, DashboardData } from '../../types'

export function DashboardPage() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    adminApi
      .dashboard()
      .then(setData)
      .catch((error) => message.error(error.message))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="admin-page">
      <Typography.Title level={2}>仪表盘</Typography.Title>
      <Row gutter={[16, 16]}>
        <StatCard title="总文章数" value={data?.totalArticles ?? 0} loading={loading} />
        <StatCard title="总标签数" value={data?.totalTags ?? 0} loading={loading} />
        <StatCard title="总浏览量" value={data?.totalViews ?? 0} loading={loading} />
        <StatCard title="今日访问量" value={data?.todayViews ?? 0} loading={loading} />
      </Row>
      <Row gutter={[16, 16]} className="mt-16">
        <Col xs={24} xl={14}>
          <Card title="近 30 天访问趋势" loading={loading}>
            <Line
              height={300}
              data={data?.viewTrend ?? []}
              xField="date"
              yField="views"
              point={{ size: 4 }}
              color="#2563EB"
            />
          </Card>
        </Col>
        <Col xs={24} xl={10}>
          <Card title="最近发布文章">
            <Table<ArticleSummary>
              rowKey="id"
              loading={loading}
              pagination={false}
              dataSource={data?.recentArticles ?? []}
              columns={[
                { title: '标题', dataIndex: 'title' },
                { title: '分类', dataIndex: 'category', width: 110 },
                { title: '浏览量', dataIndex: 'viewCount', width: 90 },
              ]}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

function StatCard({ title, value, loading }: { title: string; value: number; loading: boolean }) {
  return (
    <Col xs={12} lg={6}>
      <Card loading={loading}>
        <Statistic title={title} value={value} />
      </Card>
    </Col>
  )
}

