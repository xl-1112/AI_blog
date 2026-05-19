import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { InboxOutlined } from '@ant-design/icons'
import { Button, Card, Col, Form, Input, Row, Select, Space, Switch, Typography, Upload, message } from 'antd'
import type { UploadProps } from 'antd'
import { Editor, Toolbar } from '@wangeditor/editor-for-react'
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'
import '@wangeditor/editor/dist/css/style.css'
import { adminApi, resolveAssetUrl } from '../../api/client'
import type { Article } from '../../types'

const fallbackCategories = ['产品复盘', '用户洞察', '信息架构', '产品设计', '技术笔记']
type ImageInsertFn = (url: string, alt: string, href: string) => void

function blankArticle(defaultCategory = fallbackCategories[0]): Article {
  const date = new Date().toISOString().slice(0, 10)
  return {
    id: '',
    title: '',
    slug: '',
    category: defaultCategory,
    date,
    summary: '',
    coverUrl: '',
    tags: [],
    featured: false,
    status: 'draft',
    wordCount: 0,
    viewCount: 0,
    content: '',
    seoTitle: '',
    seoKeywords: '',
    seoDescription: '',
  }
}

export function ArticleEditorPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [form] = Form.useForm<Article>()
  const [loading, setLoading] = useState(Boolean(id))
  const [saving, setSaving] = useState(false)
  const [content, setContent] = useState('')
  const [editor, setEditor] = useState<IDomEditor | null>(null)
  const [categories, setCategories] = useState(fallbackCategories)

  useEffect(() => {
    return () => {
      if (!editor) return
      editor.destroy()
      setEditor(null)
    }
  }, [editor])

  useEffect(() => {
    adminApi.tags().then((data) => {
      const nextCategories = data.list.map((item) => item.name)
      if (nextCategories.length > 0) {
        setCategories(nextCategories)
      }
      if (!id) {
        const initial = blankArticle(nextCategories[0] ?? fallbackCategories[0])
        form.setFieldsValue(initial)
        setContent(initial.content)
        return
      }
      adminApi.article(id).then((article) => {
        form.setFieldsValue({ ...article, tags: [article.category].filter(Boolean) })
        setContent(article.content)
      })
    }).catch((error) => message.error(error.message)).finally(() => setLoading(false))
  }, [form, id])

  const uploadProps = useMemo<UploadProps>(
    () => ({
      showUploadList: false,
      customRequest: async ({ file, onSuccess, onError }) => {
        try {
          const url = await adminApi.upload(file as File, 'cover')
          form.setFieldValue('coverUrl', url)
          onSuccess?.(url)
        } catch (error) {
          onError?.(error as Error)
        }
      },
    }),
    [form],
  )

  const toolbarConfig = useMemo<Partial<IToolbarConfig>>(
    () => ({
      excludeKeys: ['fullScreen', 'group-video'],
    }),
    [],
  )

  const editorConfig = useMemo<Partial<IEditorConfig>>(
    () => ({
      placeholder: '请输入文章正文，支持标题、引用、列表、表格、代码块和图片。',
      MENU_CONF: {
        uploadImage: {
          async customUpload(file: File, insertFn: ImageInsertFn) {
            try {
              const url = await adminApi.upload(file, 'article')
              const imageUrl = resolveAssetUrl(url)
              insertFn(imageUrl, file.name, imageUrl)
            } catch (error) {
              message.error(error instanceof Error ? error.message : '图片上传失败')
            }
          },
        },
      },
    }),
    [],
  )

  async function save(status: 'draft' | 'published') {
    const values = await form.validateFields()
    setSaving(true)
    try {
      const article = { ...blankArticle(categories[0]), ...values, tags: values.category ? [values.category] : [], content, status, draft: status === 'draft' }
      const saved = await adminApi.saveArticle(article, id)
      message.success(status === 'published' ? '文章已发布' : '草稿已保存')
      navigate(`/articles/edit/${saved.id}`, { replace: true })
    } catch (error) {
      message.error(error instanceof Error ? error.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="admin-page">
      <Typography.Title level={2}>{id ? '编辑文章' : '新增文章'}</Typography.Title>
      <Form form={form} layout="vertical" disabled={loading}>
        <Row gutter={[16, 16]}>
          <Col xs={24} xl={16}>
            <Card title="基础信息">
              <Row gutter={16}>
                <Col xs={24} md={12}>
                  <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
                    <Input />
                  </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                  <Form.Item name="slug" label="slug">
                    <Input placeholder="留空时使用文章 ID" />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="summary" label="摘要" rules={[{ required: true, message: '请输入摘要' }]}>
                <Input.TextArea rows={3} />
              </Form.Item>
              <Row gutter={16}>
                <Col xs={24} md={8}>
                  <Form.Item name="category" label="分类标签" rules={[{ required: true }]}>
                    <Select options={categories.map((value) => ({ label: value, value }))} />
                  </Form.Item>
                </Col>
                <Col xs={24} md={8}>
                  <Form.Item name="date" label="发布时间" rules={[{ required: true }]}>
                    <Input type="date" />
                  </Form.Item>
                </Col>
              </Row>
            </Card>
            <Card title="内容" className="mt-16">
              <div className="article-rich-editor">
                <Toolbar editor={editor} defaultConfig={toolbarConfig} mode="default" className="article-rich-toolbar" />
                <Editor
                  defaultConfig={editorConfig}
                  value={content}
                  onCreated={(nextEditor) => setEditor(nextEditor)}
                  onChange={(nextEditor) => setContent(nextEditor.getHtml())}
                  mode="default"
                  className="article-rich-body"
                />
              </div>
            </Card>
          </Col>
          <Col xs={24} xl={8}>
            <Card title="发布设置">
              <Form.Item name="coverUrl" label="封面">
                <Input />
              </Form.Item>
              <Upload.Dragger {...uploadProps}>
                <p className="ant-upload-drag-icon">
                  <InboxOutlined />
                </p>
                <p>点击或拖拽上传封面</p>
              </Upload.Dragger>
              <Form.Item shouldUpdate noStyle>
                {() => {
                  const cover = form.getFieldValue('coverUrl') as string
                  return cover ? <img className="cover-preview" src={resolveAssetUrl(cover)} alt="封面预览" /> : null
                }}
              </Form.Item>
              <Form.Item name="featured" label="设为精选" valuePropName="checked">
                <Switch />
              </Form.Item>
            </Card>
            <Card title="SEO" className="mt-16">
              <Form.Item name="seoTitle" label="SEO 标题">
                <Input />
              </Form.Item>
              <Form.Item name="seoKeywords" label="SEO Keywords">
                <Input />
              </Form.Item>
              <Form.Item name="seoDescription" label="SEO Description">
                <Input.TextArea rows={4} />
              </Form.Item>
            </Card>
            <Card className="mt-16">
              <Space>
                <Button loading={saving} onClick={() => save('draft')}>
                  保存草稿
                </Button>
                <Button type="primary" loading={saving} onClick={() => save('published')}>
                  发布
                </Button>
              </Space>
            </Card>
          </Col>
        </Row>
      </Form>
    </div>
  )
}
