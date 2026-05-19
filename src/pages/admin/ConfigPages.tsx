import { useEffect, useRef, useState } from 'react'
import { Button, Card, Empty, Form, Input, Modal, Select, Space, Switch, Table, Tag, Typography, message } from 'antd'
import { adminApi } from '../../api/client'
import type { AboutContent, ExperienceItem, HomeContent, SiteSettings, TagUsage, WorkStackGroup } from '../../types'

export function HomeConfigPage() {
  const [form] = Form.useForm<HomeContent>()
  const [articles, setArticles] = useState<{ label: string; value: string }[]>([])
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    Promise.all([adminApi.home(), adminApi.articles({ page: 1, pageSize: 1000 })])
      .then(([home, result]) => {
        form.setFieldsValue(home)
        setArticles(result.list.map((article) => ({ label: article.title, value: article.id })))
      })
      .catch((error) => message.error(error.message))
      .finally(() => setLoading(false))
  }, [form])
  return (
    <EditCard title="首页配置" loading={loading} onSave={() => saveForm(form, adminApi.saveHome, '首页配置已保存')}>
      <Form form={form} layout="vertical">
        <Form.Item name="introTitle" label="主标题" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="subtitle" label="副标题">
          <Input />
        </Form.Item>
        <Form.Item name="intro" label="介绍">
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item name="primaryCtaText" label="主 CTA 文案">
          <Input />
        </Form.Item>
        <Form.Item name="secondaryCtaText" label="次 CTA 文案">
          <Input />
        </Form.Item>
        <Form.Item name="featuredArticleIds" label="精选文章">
          <Select mode="multiple" options={articles} />
        </Form.Item>
      </Form>
    </EditCard>
  )
}

export function AboutAdminPage() {
  const [form] = Form.useForm<AboutContent>()
  const [modalForm] = Form.useForm<AboutContent>()
  const [data, setData] = useState<AboutContent | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [activeBlock, setActiveBlock] = useState<AboutBlockKey | null>(null)

  useEffect(() => {
    adminApi
      .about()
      .then((nextData) => {
        const normalized = normalizeAboutConfig(nextData)
        setData(normalized)
        form.setFieldsValue(normalized)
      })
      .catch((error) => message.error(error.message))
      .finally(() => setLoading(false))
  }, [form])

  const blocks: AboutBlock[] = [
    { key: 'hero', name: '顶部介绍' },
    { key: 'contact', name: '联系卡片' },
    { key: 'skills', name: '技术栈区块' },
    { key: 'timeline', name: '工作经历区块' },
    { key: 'profile', name: '基础资料' },
  ]

  function openEditor(block: AboutBlockKey) {
    if (!data) return
    modalForm.setFieldsValue(data)
    setActiveBlock(block)
  }

  async function saveActiveBlock() {
    if (!data) return
    const values = await modalForm.validateFields()
    const nextData = normalizeAboutConfig({ ...data, ...values })
    setSaving(true)
    try {
      const saved = normalizeAboutConfig(await adminApi.saveAbout(nextData))
      setData(saved)
      form.setFieldsValue(saved)
      message.success('关于页配置已保存')
      setActiveBlock(null)
    } catch (error) {
      message.error(error instanceof Error ? error.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="admin-page">
      <Typography.Title level={2}>关于页配置</Typography.Title>
      <Card loading={loading}>
        <Table<AboutBlock>
          rowKey="key"
          dataSource={blocks}
          pagination={false}
          columns={[
            { title: '区块名称', dataIndex: 'name' },
            {
              title: '操作',
              width: 160,
              render: (_, record) => (
                <Button type="link" onClick={() => openEditor(record.key)}>
                  编辑
                </Button>
              ),
            },
          ]}
        />
      </Card>

      <Modal
        title={activeBlock ? `编辑${blocks.find((block) => block.key === activeBlock)?.name ?? ''}` : '编辑区块'}
        open={Boolean(activeBlock)}
        width={720}
        confirmLoading={saving}
        onOk={saveActiveBlock}
        onCancel={() => setActiveBlock(null)}
        destroyOnHidden
      >
        <Form form={modalForm} layout="vertical">
          {activeBlock ? renderAboutBlockForm(activeBlock) : null}
        </Form>
      </Modal>
    </div>
  )
}

type AboutBlockKey = 'hero' | 'contact' | 'skills' | 'timeline' | 'profile'
type AboutBlock = { key: AboutBlockKey; name: string }

function renderAboutBlockForm(block: AboutBlockKey) {
  switch (block) {
    case 'hero':
      return (
        <>
          <VisibilitySwitches
            items={[
              ['showHero', '显示区块'],
              ['showHeroBadge', '显示标签'],
              ['showHeroTitle', '显示主标题'],
              ['showHeroSubtitle', '显示副标题'],
              ['showHeroDescription', '显示介绍'],
            ]}
          />
          <Form.Item name="heroBadge" label="标签" rules={[{ required: true, message: '请输入标签' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="heroTitle" label="主标题" rules={[{ required: true, message: '请输入主标题' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="heroSubtitle" label="副标题">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="description" label="介绍">
            <Input.TextArea rows={3} />
          </Form.Item>
        </>
      )
    case 'contact':
      return (
        <>
          <VisibilitySwitches
            items={[
              ['showContact', '显示区块'],
              ['showContactBadge', '显示标签'],
              ['showContactTitle', '显示主标题'],
              ['showContactSubtitle', '显示副标题'],
              ['showContactDescription', '显示介绍'],
              ['showLocation', '显示位置'],
              ['showEmail', '显示邮箱'],
              ['showGithub', '显示 GitHub'],
            ]}
          />
          <Form.Item name="contactBadge" label="标签">
            <Input />
          </Form.Item>
          <Form.Item name="contactTitle" label="主标题">
            <Input />
          </Form.Item>
          <Form.Item name="contactSubtitle" label="副标题">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="contactDescription" label="介绍">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="city" label="位置">
            <Input />
          </Form.Item>
          <Form.Item name={['contact', 'email']} label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name={['contact', 'github']} label="GitHub">
            <Input />
          </Form.Item>
        </>
      )
    case 'skills':
      return (
        <>
          <VisibilitySwitches
            items={[
              ['showSkills', '显示区块'],
              ['showSkillsBadge', '显示标签'],
              ['showSkillsTitle', '显示主标题'],
              ['showSkillsSubtitle', '显示副标题'],
              ['showSkillsDescription', '显示介绍'],
            ]}
          />
          <Form.Item name="skillsBadge" label="标签" rules={[{ required: true, message: '请输入标签' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="skillsTitle" label="主标题" rules={[{ required: true, message: '请输入主标题' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="skillsSubtitle" label="副标题">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="skillsDescription" label="介绍">
            <Input.TextArea rows={3} />
          </Form.Item>
        </>
      )
    case 'timeline':
      return (
        <>
          <VisibilitySwitches
            items={[
              ['showTimeline', '显示区块'],
              ['showTimelineBadge', '显示标签'],
              ['showTimelineTitle', '显示主标题'],
              ['showTimelineSubtitle', '显示副标题'],
              ['showTimelineDescription', '显示介绍'],
            ]}
          />
          <Form.Item name="timelineBadge" label="标签" rules={[{ required: true, message: '请输入标签' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="timelineTitle" label="主标题" rules={[{ required: true, message: '请输入主标题' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="timelineSubtitle" label="副标题">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="timelineDescription" label="介绍">
            <Input.TextArea rows={3} />
          </Form.Item>
        </>
      )
    case 'profile':
      return (
        <>
          <Form.Item name="name" label="姓名">
            <Input />
          </Form.Item>
          <Form.Item name="title" label="资料标题">
            <Input />
          </Form.Item>
          <Form.Item name="avatarUrl" label="头像 URL">
            <Input />
          </Form.Item>
          <Form.Item name="richDescription" label="富文本简介">
            <Input.TextArea rows={6} />
          </Form.Item>
        </>
      )
    default:
      return null
  }
}

function VisibilitySwitches({ items }: { items: Array<[keyof AboutContent, string]> }) {
  return (
    <Space wrap size={[24, 12]} className="mb-16">
      {items.map(([name, label]) => (
        <Form.Item key={String(name)} name={name} label={label} valuePropName="checked" className="mb-0">
          <Switch />
        </Form.Item>
      ))}
    </Space>
  )
}

function normalizeAboutConfig(data: AboutContent): AboutContent {
  return {
    ...data,
    showHero: data.showHero ?? true,
    showHeroBadge: data.showHeroBadge ?? true,
    showHeroTitle: data.showHeroTitle ?? true,
    showHeroSubtitle: data.showHeroSubtitle ?? true,
    showHeroDescription: data.showHeroDescription ?? true,
    showContact: data.showContact ?? true,
    showLocation: data.showLocation ?? true,
    showEmail: data.showEmail ?? true,
    showGithub: data.showGithub ?? true,
    showContactBadge: data.showContactBadge ?? true,
    showContactTitle: data.showContactTitle ?? true,
    showContactSubtitle: data.showContactSubtitle ?? true,
    showContactDescription: data.showContactDescription ?? true,
    contactBadge: data.contactBadge || '联系我',
    contactTitle: data.contactTitle || '欢迎交流产品、增长和项目协作',
    contactSubtitle: data.contactSubtitle || '',
    contactDescription: data.contactDescription || '',
    showSkills: data.showSkills ?? true,
    showSkillsHeader: data.showSkillsHeader ?? true,
    showSkillsBadge: data.showSkillsBadge ?? true,
    showSkillsTitle: data.showSkillsTitle ?? true,
    showSkillsSubtitle: data.showSkillsSubtitle ?? true,
    showSkillsDescription: data.showSkillsDescription ?? true,
    skillsSubtitle: data.skillsSubtitle || '',
    showTimeline: data.showTimeline ?? true,
    showTimelineHeader: data.showTimelineHeader ?? true,
    showTimelineBadge: data.showTimelineBadge ?? true,
    showTimelineTitle: data.showTimelineTitle ?? true,
    showTimelineSubtitle: data.showTimelineSubtitle ?? true,
    showTimelineDescription: data.showTimelineDescription ?? true,
    timelineSubtitle: data.timelineSubtitle || '',
    timelineDescription: data.timelineDescription || '',
  }
}

export function SkillsPage() {
  const [items, setItems] = useState<WorkStackGroup[]>([])
  const [saving, setSaving] = useState(false)
  useEffect(() => {
    adminApi.skills().then(setItems).catch((error) => message.error(error.message))
  }, [])

  function addCategory() {
    setItems([...items, { title: `新分类 ${items.length + 1}`, items: [] }])
  }

  function updateCategory(index: number, next: WorkStackGroup) {
    setItems(items.map((item, itemIndex) => (itemIndex === index ? next : item)))
  }

  function moveCategory(index: number, direction: -1 | 1) {
    const target = index + direction
    if (target < 0 || target >= items.length) return
    const next = [...items]
    const current = next[index]
    next[index] = next[target]
    next[target] = current
    setItems(next)
  }

  async function saveSkills() {
    const cleaned = items
      .map((item) => ({
        ...item,
        title: item.title.trim(),
        items: Array.from(new Set(item.items.map((value) => value.trim()).filter(Boolean))),
      }))
      .filter((item) => item.title)

    if (cleaned.length === 0) {
      message.warning('请至少保留一个技术栈分类')
      return
    }

    setSaving(true)
    try {
      const saved = await adminApi.saveSkills(cleaned)
      setItems(saved)
      message.success('技术栈已保存')
    } catch (error) {
      message.error(error instanceof Error ? error.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="admin-page">
      <div className="page-heading">
        <div>
          <Typography.Title level={2}>技术栈</Typography.Title>
          <Typography.Text type="secondary">维护前台“技术栈/能力标签”内容，每个分类下可添加多个标签。</Typography.Text>
        </div>
        <Space>
          <Button onClick={addCategory}>新增分类</Button>
          <Button type="primary" loading={saving} onClick={saveSkills}>
            保存
          </Button>
        </Space>
      </div>

      {items.length === 0 ? (
        <Card>
          <Empty description="暂无技术栈分类">
            <Button type="primary" onClick={addCategory}>
              新增第一个分类
            </Button>
          </Empty>
        </Card>
      ) : (
        <Space direction="vertical" className="full-width" size={16}>
          {items.map((item, index) => (
            <Card
              key={`${item.id ?? 'new'}-${index}`}
              title={`分类 ${index + 1}`}
              extra={
                <Space>
                  <Button disabled={index === 0} onClick={() => moveCategory(index, -1)}>
                    上移
                  </Button>
                  <Button disabled={index === items.length - 1} onClick={() => moveCategory(index, 1)}>
                    下移
                  </Button>
                  <Button danger onClick={() => setItems(items.filter((_, itemIndex) => itemIndex !== index))}>
                    删除分类
                  </Button>
                </Space>
              }
            >
              <Space direction="vertical" className="full-width" size={12}>
                <Form layout="vertical">
                  <Form.Item label="分类名称" required>
                    <Input
                      value={item.title}
                      onChange={(event) => updateCategory(index, { ...item, title: event.target.value })}
                      placeholder="例如：产品方法"
                    />
                  </Form.Item>
                  <Form.Item label="分类标签">
                    <SkillTagEditor
                      value={item.items}
                      onChange={(values) => updateCategory(index, { ...item, items: values })}
                    />
                  </Form.Item>
                </Form>
              </Space>
            </Card>
          ))}
        </Space>
      )}
    </div>
  )
}

function SkillTagEditor({ value, onChange }: { value: string[]; onChange: (value: string[]) => void }) {
  const [inputValue, setInputValue] = useState('')

  function addTag() {
    const nextTags = inputValue
      .split(/,|，|\n/)
      .map((item) => item.trim())
      .filter(Boolean)
    if (nextTags.length === 0) return
    onChange(Array.from(new Set([...value, ...nextTags])))
    setInputValue('')
  }

  function removeTag(tag: string) {
    onChange(value.filter((item) => item !== tag))
  }

  return (
    <Space direction="vertical" className="full-width" size={10}>
      <Space.Compact className="full-width">
        <Input
          value={inputValue}
          onChange={(event) => setInputValue(event.target.value)}
          onPressEnter={addTag}
          placeholder="输入标签名称，例如：用户研究"
        />
        <Button type="primary" onClick={addTag}>
          添加标签
        </Button>
      </Space.Compact>
      <div className="skill-tag-list">
        {value.length > 0 ? (
          value.map((tag) => (
            <Tag key={tag} closable onClose={() => removeTag(tag)}>
              {tag}
            </Tag>
          ))
        ) : (
          <Typography.Text type="secondary">暂无标签，输入后点击“添加标签”。</Typography.Text>
        )}
      </div>
    </Space>
  )
}

export function TimelinePage() {
  const [items, setItems] = useState<ExperienceItem[]>([])
  useEffect(() => {
    adminApi.timeline().then(setItems).catch((error) => message.error(error.message))
  }, [])
  return (
    <ArrayEditor
      title="工作经历"
      items={items}
      setItems={setItems}
      newItem={{ period: '', title: '', body: '' }}
      render={(item, index, update) => (
        <Space direction="vertical" className="full-width">
          <Input value={item.period} onChange={(event) => update(index, { ...item, period: event.target.value })} placeholder="时间" />
          <Input value={item.title} onChange={(event) => update(index, { ...item, title: event.target.value })} placeholder="标题" />
          <RichTextEditor value={item.body} onChange={(body) => update(index, { ...item, body })} />
        </Space>
      )}
      onSave={async () => {
        await adminApi.saveTimeline(items)
        message.success('工作经历已保存')
      }}
    />
  )
}

function RichTextEditor({ value, onChange }: { value: string; onChange: (value: string) => void }) {
  const editorRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const editor = editorRef.current
    if (editor && editor.innerHTML !== value) {
      editor.innerHTML = value || ''
    }
  }, [value])

  function runCommand(command: string) {
    const editor = editorRef.current
    if (!editor) return
    editor.focus()
    document.execCommand(command)
    onChange(editor.innerHTML)
  }

  return (
    <div className="rich-text-editor">
      <div className="rich-text-toolbar">
        <Button size="small" onMouseDown={(event) => event.preventDefault()} onClick={() => runCommand('bold')}>
          加粗
        </Button>
        <Button size="small" onMouseDown={(event) => event.preventDefault()} onClick={() => runCommand('italic')}>
          斜体
        </Button>
        <Button size="small" onMouseDown={(event) => event.preventDefault()} onClick={() => runCommand('insertUnorderedList')}>
          无序列表
        </Button>
        <Button size="small" onMouseDown={(event) => event.preventDefault()} onClick={() => runCommand('insertOrderedList')}>
          有序列表
        </Button>
        <Button size="small" onMouseDown={(event) => event.preventDefault()} onClick={() => runCommand('removeFormat')}>
          清除格式
        </Button>
      </div>
      <div
        ref={editorRef}
        className="rich-text-body"
        contentEditable
        suppressContentEditableWarning
        data-placeholder="请输入工作经历描述，可使用加粗、列表等富文本格式"
        onInput={(event) => onChange(event.currentTarget.innerHTML)}
      />
    </div>
  )
}

export function TagsPage() {
  const [tags, setTags] = useState<TagUsage[]>([])
  const [text, setText] = useState('')
  useEffect(() => {
    adminApi.tags().then((data) => {
      setTags(data.list)
      setText(data.list.map((tag) => tag.name).join('\n'))
    })
  }, [])
  return (
    <div className="admin-page">
      <Typography.Title level={2}>分类标签</Typography.Title>
      <Card className="mb-16">
        <Input.TextArea value={text} onChange={(event) => setText(event.target.value)} rows={8} placeholder="每行一个分类标签，例如：产品复盘" />
        <Button type="primary" className="mt-16" onClick={() => adminApi.saveTags(text.split(/\n|,|，/).map((item) => item.trim()).filter(Boolean)).then(() => message.success('标签已保存'))}>
          保存分类标签
        </Button>
      </Card>
      <Card>
        <Table<TagUsage> rowKey="name" dataSource={tags} pagination={false} columns={[{ title: '分类标签名称', dataIndex: 'name' }, { title: '文章使用次数', dataIndex: 'useCount' }, { title: '创建时间', dataIndex: 'createdAt' }]} />
      </Card>
    </div>
  )
}

export function SiteConfigPage() {
  const [form] = Form.useForm<SiteSettings>()
  const [loading, setLoading] = useState(true)
  useEffect(() => {
    adminApi.site().then((data) => form.setFieldsValue(data)).catch((error) => message.error(error.message)).finally(() => setLoading(false))
  }, [form])
  return (
    <EditCard title="网站配置" loading={loading} onSave={() => saveForm(form, adminApi.saveSite, '网站配置已保存')}>
      <Form form={form} layout="vertical">
        <Form.Item name="name" label="站点名称" rules={[{ required: true }]}>
          <Input />
        </Form.Item>
        <Form.Item name="siteTitle" label="网站标题">
          <Input />
        </Form.Item>
        <Form.Item name="description" label="Description">
          <Input.TextArea rows={3} />
        </Form.Item>
        <Form.Item name="keywords" label="Keywords">
          <Input />
        </Form.Item>
        <Form.Item name="logoUrl" label="Logo URL">
          <Input />
        </Form.Item>
        <Form.Item name="faviconUrl" label="Favicon URL">
          <Input />
        </Form.Item>
        <Form.Item name="icp" label="ICP备案号">
          <Input />
        </Form.Item>
        <Form.Item name="analyticsCode" label="Google Analytics 统计代码">
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item name="role" label="角色定位">
          <Input />
        </Form.Item>
        <Form.Item name="location" label="城市">
          <Input />
        </Form.Item>
        <Form.Item name={['contact', 'email']} label="邮箱">
          <Input />
        </Form.Item>
        <Form.Item name={['contact', 'github']} label="GitHub">
          <Input />
        </Form.Item>
      </Form>
    </EditCard>
  )
}

function EditCard<T extends object>({ title, loading, children, onSave }: { title: string; loading: boolean; children: React.ReactNode; onSave: () => Promise<void> }) {
  void (null as T | null)
  return (
    <div className="admin-page">
      <div className="page-heading">
        <Typography.Title level={2}>{title}</Typography.Title>
        <Button type="primary" onClick={onSave}>
          保存
        </Button>
      </div>
      <Card loading={loading}>{children}</Card>
    </div>
  )
}

async function saveForm<T extends object>(form: ReturnType<typeof Form.useForm<T>>[0], save: (value: T) => Promise<T>, success: string) {
  const values = await form.validateFields()
  await save(values)
  message.success(success)
}

function ArrayEditor<T extends object>({
  title,
  items,
  setItems,
  newItem,
  render,
  onSave,
}: {
  title: string
  items: T[]
  setItems: (items: T[]) => void
  newItem: T
  render: (item: T, index: number, update: (index: number, item: T) => void) => React.ReactNode
  onSave: () => Promise<void>
}) {
  const update = (index: number, item: T) => setItems(items.map((current, currentIndex) => (currentIndex === index ? item : current)))
  return (
    <div className="admin-page">
      <div className="page-heading">
        <Typography.Title level={2}>{title}</Typography.Title>
        <Space>
          <Button onClick={() => setItems([...items, newItem])}>新增</Button>
          <Button type="primary" onClick={onSave}>保存</Button>
        </Space>
      </div>
      <Space direction="vertical" className="full-width" size={16}>
        {items.map((item, index) => (
          <Card key={index} extra={<Button danger onClick={() => setItems(items.filter((_, itemIndex) => itemIndex !== index))}>删除</Button>}>
            {render(item, index, update)}
          </Card>
        ))}
      </Space>
    </div>
  )
}
