import { useEffect, useState } from 'react'
import { Bar, Line } from '@ant-design/plots'
import { Card, Col, Row, Table, Typography, message } from 'antd'
import { adminApi } from '../../api/client'
import type { AnalyticsData, ArticleSummary } from '../../types'

export function AnalyticsPage() {
  const [data, setData] = useState<AnalyticsData | null>(null)
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    adminApi.analytics().then(setData).catch((error) => message.error(error.message)).finally(() => setLoading(false))
  }, [])
  return (
    <div className="admin-page">
      <Typography.Title level={2}>数据统计</Typography.Title>
      <Row gutter={[16, 16]}>
        <Col xs={24} xl={14}>
          <Card title="浏览趋势" loading={loading}>
            <Line height={320} data={data?.viewTrend ?? []} xField="date" yField="views" color="#2563EB" />
          </Card>
        </Col>
        <Col xs={24} xl={10}>
          <Card title="热门文章排行" loading={loading}>
            <Bar height={320} data={data?.hotArticles ?? []} xField="viewCount" yField="title" color="#2563EB" />
          </Card>
        </Col>
      </Row>
      <Card title="热门文章明细" className="mt-16">
        <Table<ArticleSummary>
          rowKey="id"
          pagination={false}
          dataSource={data?.hotArticles ?? []}
          columns={[
            { title: '标题', dataIndex: 'title' },
            { title: '分类', dataIndex: 'category' },
            { title: '浏览量', dataIndex: 'viewCount' },
          ]}
        />
      </Card>
    </div>
  )
}

