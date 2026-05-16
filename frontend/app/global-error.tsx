'use client'
import * as React from 'react'

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <html>
      <body>
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '100vh', gap: '16px', fontFamily: 'sans-serif', background: '#0E0E12', color: '#fff' }}>
          <h2>Что-то пошло не так</h2>
          <button onClick={() => reset()} style={{ padding: '8px 16px', background: '#6C63FF', color: '#fff', border: 'none', borderRadius: '8px', cursor: 'pointer' }}>
            Попробовать снова
          </button>
        </div>
      </body>
    </html>
  )
}
