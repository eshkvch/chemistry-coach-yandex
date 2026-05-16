CREATE TABLE IF NOT EXISTS users (
    id         TEXT PRIMARY KEY,
    birth_year INTEGER NOT NULL,
    initials   TEXT NOT NULL DEFAULT 'ИВ',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    id            TEXT PRIMARY KEY,
    user_id       TEXT REFERENCES users(id) ON DELETE CASCADE,
    goal_id       TEXT NOT NULL,
    persona_id    TEXT NOT NULL,
    status        TEXT DEFAULT 'active',
    system_prompt TEXT,
    started_at    TIMESTAMP DEFAULT NOW(),
    finished_at   TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id           TEXT PRIMARY KEY,
    session_id   TEXT REFERENCES sessions(id) ON DELETE CASCADE,
    sender       TEXT NOT NULL,
    text         TEXT NOT NULL,
    status       TEXT,
    status_label TEXT,
    clarity      INTEGER,
    confidence   INTEGER,
    respect      INTEGER,
    balance      INTEGER,
    consent_risk BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS debriefs (
    id                    TEXT PRIMARY KEY,
    session_id            TEXT REFERENCES sessions(id) ON DELETE CASCADE UNIQUE,
    scores_json           JSONB NOT NULL,
    strengths_json        JSONB NOT NULL DEFAULT '[]',
    weaknesses_json       JSONB NOT NULL DEFAULT '[]',
    risk_flags_json       JSONB NOT NULL DEFAULT '[]',
    improved_replies_json JSONB NOT NULL DEFAULT '[]',
    tip_for_next          TEXT,
    has_risk              BOOLEAN DEFAULT FALSE,
    created_at            TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
