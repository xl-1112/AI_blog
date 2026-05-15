import { HashRouter, Route, Routes } from 'react-router-dom'
import { SiteHeader } from './components/SiteHeader'
import { ContentProvider } from './lib/content'
import { AboutPage } from './pages/AboutPage'
import { ArticleDetailPage } from './pages/ArticleDetailPage'
import { ArticlesPage } from './pages/ArticlesPage'
import { HomePage } from './pages/HomePage'

function App() {
  return (
    <ContentProvider>
      <HashRouter>
        <SiteHeader />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/articles" element={<ArticlesPage />} />
          <Route path="/articles/:id" element={<ArticleDetailPage />} />
          <Route path="/about" element={<AboutPage />} />
        </Routes>
      </HashRouter>
    </ContentProvider>
  )
}

export default App
