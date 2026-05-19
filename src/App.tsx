import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { HashRouter, Navigate, Route, Routes, useLocation } from 'react-router-dom'
import { AdminLayout } from './layouts/AdminLayout'
import { SiteHeader } from './components/SiteHeader'
import { ContentProvider } from './lib/content'
import { AboutPage } from './pages/AboutPage'
import { ArticleDetailPage } from './pages/ArticleDetailPage'
import { ArticlesPage } from './pages/ArticlesPage'
import { HomePage } from './pages/HomePage'
import { AnalyticsPage } from './pages/admin/AnalyticsPage'
import { ArticleEditorPage } from './pages/admin/ArticleEditorPage'
import { ArticlesAdminPage } from './pages/admin/ArticlesAdminPage'
import { DashboardPage } from './pages/admin/DashboardPage'
import { AboutAdminPage, HomeConfigPage, SiteConfigPage, SkillsPage, TagsPage, TimelinePage } from './pages/admin/ConfigPages'
import { LoginPage } from './pages/admin/LoginPage'
import { ProfilePage } from './pages/admin/ProfilePage'
import { UsersPage } from './pages/admin/UsersPage'

function AppRoutes() {
  const location = useLocation()
  const isSite = location.pathname === '/site' || location.pathname.startsWith('/site/')

  return (
    <>
      {isSite ? <SiteHeader /> : null}
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AdminLayout />}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/articles" element={<ArticlesAdminPage />} />
          <Route path="/articles/create" element={<ArticleEditorPage />} />
          <Route path="/articles/edit/:id" element={<ArticleEditorPage />} />
          <Route path="/home-config" element={<HomeConfigPage />} />
          <Route path="/about-config" element={<AboutAdminPage />} />
          <Route path="/skills" element={<SkillsPage />} />
          <Route path="/timeline" element={<TimelinePage />} />
          <Route path="/tags" element={<TagsPage />} />
          <Route path="/site-config" element={<SiteConfigPage />} />
          <Route path="/analytics" element={<AnalyticsPage />} />
          <Route path="/users" element={<UsersPage />} />
          <Route path="/profile" element={<ProfilePage />} />
        </Route>
        <Route path="/site" element={<HomePage />} />
        <Route path="/site/articles" element={<ArticlesPage />} />
        <Route path="/site/articles/:id" element={<ArticleDetailPage />} />
        <Route path="/site/about" element={<AboutPage />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </>
  )
}

function App() {
  return (
    <ConfigProvider locale={zhCN} theme={{ token: { colorPrimary: '#2563EB', borderRadius: 12 } }}>
      <ContentProvider>
        <HashRouter>
          <AppRoutes />
        </HashRouter>
      </ContentProvider>
    </ConfigProvider>
  )
}

export default App
