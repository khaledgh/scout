import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './styles/index.css'
import App from './App.tsx'

// Set RTL direction on the document root
document.documentElement.dir = import.meta.env.VITE_DEFAULT_LOCALE === 'en' ? 'ltr' : 'rtl'
document.documentElement.lang = import.meta.env.VITE_DEFAULT_LOCALE ?? 'ar'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
