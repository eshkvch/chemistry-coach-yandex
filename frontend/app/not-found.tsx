'use client'
import * as React from 'react'

export default function NotFound() {
  return (
    <div style={{ minHeight: '100vh', background: '#0E0E12', color: '#fff', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '16px', fontFamily: 'Inter, sans-serif' }}>
      <h1 style={{ fontSize: '1.5rem', fontWeight: 600 }}>404 — Страница не найдена</h1>
      <a href="/" style={{ color: '#6C63FF', textDecoration: 'underline' }}>
        На главную
      </a>
    </div>
  )
}
