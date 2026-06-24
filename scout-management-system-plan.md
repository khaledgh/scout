# Kashfi — Scout Troop Management System
### Full Architecture & Build Blueprint (نظام إدارة الفوج الكشفي)

> A production-grade, multi-role Scout Management System for a Lebanese scout troop (الفوج).
> Not just a database — a complete operations platform: members, units, activities, badges,
> attendance, training, gamification, internal comms, and analytics.
>
> **Stack:** Go (Echo + GORM) + MySQL backend · React + TypeScript + Tailwind v3 + React Query frontend
> **Conventions:** `models / handlers / services` backend layering · Zod + React Hook Form on every form · all config in `.env`

---

## Table of Contents

1. [Goals & Scope](#1-goals--scope)
2. [Tech Stack](#2-tech-stack)
3. [System Architecture](#3-system-architecture)
4. [Roles & Permissions (RBAC)](#4-roles--permissions-rbac)
5. [Domain Model & Database Schema](#5-domain-model--database-schema)
6. [Backend Structure (Go / Echo / GORM)](#6-backend-structure-go--echo--gorm)
7. [API Surface](#7-api-surface)
8. [Frontend Structure (React / TS / Tailwind)](#8-frontend-structure-react--ts--tailwind)
9. [Design System & UI Direction](#9-design-system--ui-direction)
10. [Feature Modules (Detailed)](#10-feature-modules-detailed)
11. [Environment Variables (.env)](#11-environment-variables-env)
12. [Security](#12-security)
13. [Implementation Order (Phased)](#13-implementation-order-phased)
14. [Future / WOW Features](#14-future--wow-features)

---

## 1. Goals & Scope

Build a single platform that lets scout leaders (القادة) run the entire troop digitally:

- **Members** — full profile per scout, including medical info, parent contacts, skills, attendance, and leader evaluation.
- **Units / Patrols** — create patrols (طليعة / زمرة), assign leaders, track unit-level scores and rankings.
- **Activities** — schedule camps, hikes, trainings; track attendance; collect photos & post-activity feedback.
- **Badges & Progress** — gamified skill/badge unlock system per scout.
- **Attendance** — fast QR-code check-in (with optional GPS check-in).
- **Training Hub** — lessons (knots, first aid, scout skills) with media and quizzes.
- **Gamification** — XP, levels, leaderboards, awards.
- **Communication** — internal announcements, unit chat, push notifications.
- **Reports & Dashboard** — per-member, per-unit, monthly troop reports; leader dashboard with alerts.

### Non-goals (v1)
- Public-facing website (this is an internal management tool).
- Online payments (no commerce in v1 — could be added later for camp fees).
- Native mobile app (v1 is a responsive web app / PWA; React Native can follow later).

---

## 2. Tech Stack

### Backend
| Concern | Choice |
|---|---|
| Language | Go (1.22+) |
| Web framework | Echo v4 |
| ORM | GORM |
| Database | MySQL 8 |
| Auth | JWT (access + refresh), bcrypt password hashing |
| Validation | `go-playground/validator` |
| Config | `.env` via `github.com/joho/godotenv` + typed config struct |
| File storage | Local disk (`/uploads`) with served static route (S3-compatible later) |
| Realtime | WebSocket (`gorilla/websocket` or Echo's built-in upgrade) for chat & live attendance |
| Migrations | GORM `AutoMigrate` for dev + versioned SQL migrations for prod (`golang-migrate`) |
| Background jobs | Lightweight in-process worker (goroutine + ticker) for reminders / report generation |

### Frontend
| Concern | Choice |
|---|---|
| Framework | React 18 + Vite |
| Language | TypeScript (strict mode) |
| Styling | Tailwind CSS **v3** |
| Server state | TanStack React Query v5 |
| Forms | React Hook Form + **Zod** (validation on **every** form) |
| Routing | React Router v6 |
| HTTP | Axios instance with interceptors (auth header + refresh) |
| Icons | lucide-react |
| Charts | Recharts |
| Notifications (UI) | sonner / react-hot-toast |
| i18n / RTL | Full Arabic (RTL) + optional French/English, `dir="rtl"` aware layout |
| PWA | `vite-plugin-pwa` (installable, offline shell) |

---

## 3. System Architecture

```
┌──────────────────────────────────────────────────────────┐
│                      Client (Browser / PWA)                │
│   React + TS + Tailwind v3 + React Query + RHF/Zod         │
└───────────────┬──────────────────────────┬────────────────┘
                │ REST (JSON)               │ WebSocket
                ▼                           ▼
┌──────────────────────────────────────────────────────────┐
│                   Go API (Echo)                            │
│   Router → Middleware (Auth, RBAC, CORS, Logger, Recover)  │
│      → Handlers → Services → Models (GORM) → MySQL         │
│   WS Hub (chat / live attendance) · Job worker (reminders) │
└───────────────┬──────────────────────────────────────────┘
                ▼
         ┌─────────────┐        ┌─────────────────┐
         │   MySQL 8   │        │  /uploads (disk)│
         └─────────────┘        └─────────────────┘
```

**Layering rule:** Handlers never touch GORM directly. Handlers parse/validate input and shape responses; **Services** hold business logic and own all DB access through models; **Models** are GORM structs + query helpers.

**Deployment target:** single VPS (aaPanel/BaoTa-friendly). Go compiled to one binary; React built to static assets served either by the Go binary (embedded via `embed.FS`) or by Nginx. Recommended: Go binary serves the compiled SPA so the whole app runs on one port.

---

## 4. Roles & Permissions (RBAC)

| Role | Description | Key capabilities |
|---|---|---|
| `super_admin` | Troop chief / system owner | Everything, manage leaders, system settings |
| `leader` | Unit leader (قائد) | Manage own units' members, activities, attendance, evaluations |
| `assistant` | Assistant leader (مساعد) | Same as leader but scoped, no member deletion |
| `member` | Scout (عنصر) | View own profile, badges, activities, training, quizzes, unit chat |
| `parent` | Parent (الأهل) — optional v1.1 | View their child's attendance, activities, announcements |

- Permissions enforced by **middleware** (`RequireRole(...)`) + **service-level scoping** (a leader only sees members in units they lead).
- Store role on the user; store unit-leadership in a join table (`unit_leaders`) so a leader can lead multiple units.

---

## 5. Domain Model & Database Schema

All tables use `id BIGINT UNSIGNED AUTO_INCREMENT`, `created_at`, `updated_at`, soft-delete `deleted_at`.

### Core entities

**users**
- `id`, `full_name`, `email` (nullable), `phone`, `password_hash`, `role` (enum), `avatar_url`, `is_active`, `last_login_at`

**members** (scout profile — linked to a user when the scout has a login)
- `id`, `user_id` (nullable FK), `full_name`, `birth_date`, `gender`, `section` (enum: `ashbal`/أشبال, `kashaf`/كشاف, `jawala`/جوالة, etc.), `rank_stage`, `join_date`, `photo_url`
- `parent_name`, `parent_phone`, `secondary_phone`, `address`
- `xp_total`, `level`, `status` (active/inactive)

**member_medical** (1:1 with member — sensitive, separate table)
- `member_id` (FK), `blood_type`, `allergies`, `chronic_conditions`, `medications`, `emergency_notes`

**units** (patrols / طلائع)
- `id`, `name`, `section` (enum), `motto`, `emblem_url`, `level`, `score_total`, `is_active`

**unit_leaders** (M:N users↔units, role-scoped leadership)
- `unit_id` (FK), `user_id` (FK), `role_in_unit` (leader/assistant)

**unit_members** (M:N members↔units; usually a member belongs to one unit but model M:N for flexibility/history)
- `unit_id` (FK), `member_id` (FK), `is_primary`, `joined_at`, `left_at`

**activities**
- `id`, `title`, `description`, `type` (enum: camp/مخيم, hike/مسير, training/تدريب, meeting/اجتماع, service/خدمة), `location`, `location_lat`, `location_lng`, `starts_at`, `ends_at`, `responsible_user_id` (FK), `unit_id` (nullable — null = whole troop), `status` (enum: planned/ongoing/completed/cancelled), `cover_image_url`

**activity_attendance**
- `id`, `activity_id` (FK), `member_id` (FK), `status` (enum: present/absent/excused/late), `check_in_at`, `check_in_method` (qr/gps/manual), `lat`, `lng`, `recorded_by` (FK user)

**activity_feedback**
- `id`, `activity_id` (FK), `member_id` (FK), `rating` (1–5), `what_went_well`, `what_to_improve`, `comment`

**activity_media**
- `id`, `activity_id` (FK), `url`, `media_type` (image/video), `uploaded_by` (FK), `caption`

**badges** (catalog)
- `id`, `name`, `description`, `category`, `icon_url`, `xp_reward`, `criteria_json` (rules for auto-unlock), `is_active`

**member_badges** (earned)
- `id`, `member_id` (FK), `badge_id` (FK), `awarded_at`, `awarded_by` (FK user, null = auto), `progress` (0–100)

**skills** (catalog: leadership, first aid, sports…)
- `id`, `name`, `category`, `description`, `max_level`

**member_skills**
- `id`, `member_id` (FK), `skill_id` (FK), `level`, `assessed_by` (FK user), `assessed_at`

**evaluations** (leader's evaluation of a member)
- `id`, `member_id` (FK), `evaluator_id` (FK user), `period` (e.g. `2026-Q2`), `discipline`, `participation`, `leadership`, `skill`, `overall`, `notes`, `created_at`

**training_lessons**
- `id`, `title`, `category` (knots/first-aid/scout-skills…), `content` (rich text/markdown), `cover_url`, `media_json` (videos/images), `order_index`, `is_published`

**quizzes** / **quiz_questions** / **quiz_attempts**
- quiz: `id`, `lesson_id` (FK), `title`, `pass_score`, `xp_reward`
- question: `id`, `quiz_id`, `text`, `options_json`, `correct_index`, `points`
- attempt: `id`, `quiz_id`, `member_id`, `score`, `passed`, `answers_json`, `attempted_at`

**announcements**
- `id`, `title`, `body`, `audience` (enum: all/unit/leaders), `unit_id` (nullable), `author_id` (FK), `pinned`, `published_at`

**messages** (unit / direct chat)
- `id`, `channel_id` (FK), `sender_id` (FK), `body`, `attachment_url`, `created_at`
- **channels**: `id`, `type` (unit/group/direct), `unit_id` (nullable), `name`

**notifications**
- `id`, `user_id` (FK), `title`, `body`, `type`, `data_json`, `read_at`, `created_at`

**equipment** (gear: tents, ropes…)
- `id`, `name`, `category`, `quantity_total`, `quantity_available`, `condition`, `notes`

**equipment_loans**
- `id`, `equipment_id` (FK), `borrowed_by` (FK), `activity_id` (nullable FK), `quantity`, `due_date`, `returned_at`

**xp_events** (audit trail for gamification)
- `id`, `member_id` (FK), `unit_id` (nullable), `source` (enum: attendance/badge/quiz/leadership/manual), `points`, `ref_id`, `note`, `created_at`

> **Indexing notes:** index FKs, `activities(starts_at)`, `activity_attendance(activity_id, member_id)` unique, `member_badges(member_id, badge_id)` unique, `xp_events(member_id, created_at)`.

---

## 6. Backend Structure (Go / Echo / GORM)

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # entrypoint: load config, db, router, server
├── internal/
│   ├── config/
│   │   └── config.go            # typed struct loaded from .env (godotenv)
│   ├── db/
│   │   ├── db.go                # GORM connection + pool settings
│   │   └── migrate.go           # AutoMigrate (dev) / migration runner
│   ├── models/                  # GORM structs + query helpers ONLY
│   │   ├── user.go
│   │   ├── member.go
│   │   ├── unit.go
│   │   ├── activity.go
│   │   ├── badge.go
│   │   ├── training.go
│   │   ├── message.go
│   │   ├── equipment.go
│   │   └── ...
│   ├── services/                # business logic, all DB access lives here
│   │   ├── auth_service.go
│   │   ├── member_service.go
│   │   ├── unit_service.go
│   │   ├── activity_service.go
│   │   ├── attendance_service.go
│   │   ├── badge_service.go
│   │   ├── gamification_service.go   # XP, levels, leaderboard
│   │   ├── training_service.go
│   │   ├── report_service.go
│   │   ├── notification_service.go
│   │   └── chat_service.go
│   ├── handlers/                # Echo handlers: parse, validate, call service, respond
│   │   ├── auth_handler.go
│   │   ├── member_handler.go
│   │   ├── unit_handler.go
│   │   ├── activity_handler.go
│   │   ├── attendance_handler.go
│   │   ├── badge_handler.go
│   │   ├── training_handler.go
│   │   ├── report_handler.go
│   │   ├── announcement_handler.go
│   │   └── chat_handler.go      # WS upgrade + REST history
│   ├── middleware/
│   │   ├── auth.go              # JWT parse → context
│   │   ├── rbac.go              # RequireRole, scope helpers
│   │   ├── cors.go
│   │   └── ratelimit.go
│   ├── dto/                     # request/response structs + validate tags
│   │   ├── auth_dto.go
│   │   ├── member_dto.go
│   │   └── ...
│   ├── ws/
│   │   └── hub.go               # WebSocket hub (chat + live attendance)
│   ├── jobs/
│   │   └── scheduler.go         # activity reminders, monthly reports
│   └── utils/
│       ├── jwt.go
│       ├── hash.go
│       ├── response.go          # standard JSON envelope + error helpers
│       ├── pagination.go
│       └── upload.go            # file save + validation
├── uploads/                     # served static (gitignored)
├── migrations/                  # versioned .sql for prod
├── .env
├── .env.example
├── go.mod
└── go.sum
```

### Conventions
- **Standard response envelope:** `{ "success": bool, "data": ..., "error": { "code", "message" }, "meta": { pagination } }`.
- **Handlers**: bind DTO → `validator.Validate(dto)` → call service → map to response. No GORM in handlers.
- **Services**: receive `*gorm.DB` (or repository) + context; return domain results & typed errors. Cross-entity logic (e.g. "award XP on attendance") lives here, wrapped in a transaction where needed.
- **Errors**: define sentinel errors (`ErrNotFound`, `ErrForbidden`, `ErrConflict`) in services; map to HTTP codes in a central error handler.
- **Transactions**: keep them short — only DB writes inside. No external/network calls inside a transaction (avoids lock-wait issues).

---

## 7. API Surface

> Prefix all with `/api/v1`. Auth via `Authorization: Bearer <token>`.

### Auth
```
POST   /auth/login
POST   /auth/refresh
POST   /auth/logout
GET    /auth/me
POST   /auth/change-password
```

### Members
```
GET    /members                 # list (filter: unit, section, status, search; paginated)
POST   /members
GET    /members/:id             # full profile (badges, skills, attendance summary)
PUT    /members/:id
DELETE /members/:id
GET    /members/:id/medical     # restricted (leader/admin)
PUT    /members/:id/medical
GET    /members/:id/timeline    # activities + badges + xp events
POST   /members/:id/evaluate    # create evaluation
```

### Units
```
GET    /units
POST   /units
GET    /units/:id
PUT    /units/:id
DELETE /units/:id
POST   /units/:id/members       # add member(s)
DELETE /units/:id/members/:mid
POST   /units/:id/leaders       # assign leader/assistant
GET    /units/leaderboard       # ranked by score_total
```

### Activities
```
GET    /activities              # filter by type, date range, unit, status
POST   /activities
GET    /activities/:id
PUT    /activities/:id
DELETE /activities/:id
POST   /activities/:id/media    # upload photos/video
GET    /activities/:id/attendance
POST   /activities/:id/attendance        # bulk record / manual
POST   /activities/:id/checkin           # QR/GPS self-check-in
POST   /activities/:id/feedback
GET    /activities/:id/feedback/summary
```

### Badges & Skills
```
GET    /badges                  # catalog
POST   /badges                  # admin
POST   /members/:id/badges      # award badge
GET    /members/:id/badges
GET    /skills
POST   /members/:id/skills      # assess skill level
```

### Training
```
GET    /training/lessons
POST   /training/lessons
GET    /training/lessons/:id
GET    /training/lessons/:id/quiz
POST   /training/quizzes/:id/attempt
GET    /training/me/progress
```

### Gamification
```
GET    /leaderboard/members     # ranked by xp_total (filter by section/unit)
GET    /leaderboard/units
GET    /members/:id/xp          # xp event history
```

### Communication
```
GET    /announcements
POST   /announcements
GET    /channels                # chat channels user can access
GET    /channels/:id/messages   # paginated history
WS     /ws/chat                 # realtime messaging
GET    /notifications
PUT    /notifications/:id/read
```

### Reports & Dashboard
```
GET    /dashboard               # leader KPIs: member count, attendance %, top unit, at-risk members
GET    /reports/member/:id
GET    /reports/unit/:id
GET    /reports/monthly?month=2026-06
GET    /reports/export?type=...&format=csv|pdf
```

### Equipment (v1.1)
```
GET    /equipment
POST   /equipment
POST   /equipment/:id/loan
PUT    /equipment/loans/:id/return
```

---

## 8. Frontend Structure (React / TS / Tailwind)

```
frontend/
├── src/
│   ├── main.tsx
│   ├── App.tsx                  # router + providers (QueryClient, Auth, i18n, Theme)
│   ├── lib/
│   │   ├── api.ts               # axios instance + interceptors
│   │   ├── queryClient.ts
│   │   └── i18n.ts              # ar (RTL) / fr / en
│   ├── types/                   # shared TS types (mirror backend DTOs)
│   ├── schemas/                 # Zod schemas (one per form)
│   │   ├── auth.schema.ts
│   │   ├── member.schema.ts
│   │   ├── activity.schema.ts
│   │   └── ...
│   ├── hooks/                   # React Query hooks per domain
│   │   ├── useAuth.ts
│   │   ├── useMembers.ts
│   │   ├── useUnits.ts
│   │   ├── useActivities.ts
│   │   ├── useBadges.ts
│   │   ├── useTraining.ts
│   │   ├── useLeaderboard.ts
│   │   └── useDashboard.ts
│   ├── components/
│   │   ├── ui/                  # design-system primitives (Button, Card, Input, Modal, Badge, Table, Tabs, Avatar…)
│   │   ├── layout/              # Sidebar, Topbar, PageHeader, RTLProvider
│   │   ├── forms/               # Form fields wired to RHF + Zod
│   │   └── charts/              # Recharts wrappers
│   ├── features/
│   │   ├── auth/
│   │   ├── dashboard/
│   │   ├── members/            # list, profile, create/edit, medical, evaluation
│   │   ├── units/              # list, detail, leaderboard
│   │   ├── activities/         # calendar, list, detail, attendance, feedback, media
│   │   ├── badges/
│   │   ├── training/           # lessons, lesson view, quiz
│   │   ├── gamification/       # leaderboards, XP timeline
│   │   ├── communication/      # announcements, chat
│   │   └── reports/
│   ├── routes/
│   │   ├── ProtectedRoute.tsx
│   │   └── RoleRoute.tsx
│   └── styles/
│       └── index.css           # Tailwind directives + theme tokens
├── public/
├── .env                        # VITE_ vars only
├── .env.example
├── tailwind.config.ts          # Tailwind v3 + theme tokens + RTL plugin
├── vite.config.ts              # + vite-plugin-pwa
├── tsconfig.json
└── package.json
```

### Form rule (every form)
Every form uses **React Hook Form + Zod**:
```ts
const form = useForm<MemberInput>({ resolver: zodResolver(memberSchema) });
```
- One Zod schema per entity in `src/schemas`.
- Reuse schema-inferred types (`z.infer<typeof memberSchema>`) as the form type **and** the API payload type.
- Inline field errors + a top-level error summary; submit disabled while `isSubmitting` / mutation pending.

### React Query rule
- One hook file per domain; query keys are arrays (`['members', filters]`).
- Mutations invalidate relevant keys + optimistic updates for attendance toggles and chat.
- Global `onError` interceptor surfaces toast + handles 401 → refresh → retry.

---

## 9. Design System & UI Direction

A **modern, energetic, scout-themed** template — clean and "fancy" (فاخر) without being heavy.

### Brand feel
- Outdoor / adventure tone: think compass, terrain, badges, campfire — but expressed through a refined modern UI, not clip-art.
- Card-based dashboards, soft shadows, rounded-2xl corners, generous whitespace, clear hierarchy.
- Gamified accents: progress rings, XP bars, badge tiles, leaderboard medals.

### Color tokens (Tailwind theme — tune to troop colors)
```
primary    #2F6B3C   // scout green (forest)
secondary  #C8932A   // gold / badge accent
accent     #1E5A8A   // trail blue
neutral    slate scale
success / warning / danger  standard
surface    near-white in light, deep slate in dark
```
- Support **light & dark** mode via CSS variables + Tailwind `dark:`.
- **RTL-first**: layout works in Arabic by default; logical properties / `dir`-aware spacing. Mirror the sidebar and icons where appropriate.

### Typography
- Arabic-friendly font (e.g. **Cairo** / **Tajawal**) + a clean Latin font (Inter) for numbers/English.
- Strong heading scale; readable body; tabular numerals for stats.

### Key screens
1. **Login** — branded split screen (illustration/photo + form).
2. **Leader Dashboard** — KPI cards (members, attendance %, top unit), at-risk members list, upcoming activities, recent activity feed, mini leaderboard.
3. **Member profile** — "scout LinkedIn": header with avatar/level/XP ring, tabs (Overview · Badges · Skills · Activities · Evaluations · Medical[restricted]).
4. **Units** — grid of unit cards with emblem, score, member count; detail page with roster + leaderboard.
5. **Activity calendar + detail** — month/agenda view; detail with attendance grid, media gallery, feedback summary.
6. **Attendance check-in** — leader scans member QR (camera) or members self-check-in; live count updates via WS.
7. **Badges** — catalog grid + per-member earned wall with locked/unlocked states.
8. **Training Hub** — lesson cards by category → lesson reader → quiz.
9. **Leaderboards** — members and units, with medals and XP bars, filter by section.
10. **Chat / Announcements** — channel list + thread; pinned announcements at top.

> Before building UI, consult the `frontend-design` skill for design tokens and styling constraints. Keep components in `components/ui` reusable and theme-driven (no hardcoded colors — use tokens).

---

## 10. Feature Modules (Detailed)

### 10.1 Members
- CRUD with filters (unit, section, status, search by name/phone).
- Separate **medical** sub-resource, visible only to `leader`/`admin`.
- Profile aggregates: total XP, level, badges earned, attendance rate, latest evaluation.
- Phone numbers stored normalized (Lebanese format handling: accept `03/70/71/76/78/79/81`, `+961`, etc.; normalize to E.164 on save).

### 10.2 Units / Patrols
- Assign one leader + assistant(s) via `unit_leaders`.
- `score_total` accumulates from member XP + unit challenges; drives `/units/leaderboard`.
- Weekly ranking snapshot (job) for "unit of the week".

### 10.3 Activities
- Types: camp / hike / training / meeting / service.
- Lifecycle: planned → ongoing → completed → cancelled.
- Media gallery (images/video) per activity.
- Post-activity **feedback** (rating + what-went-well / what-to-improve) with an aggregated summary for leaders.

### 10.4 Attendance (Smart)
- Each member has a **QR code** (encodes member id + signed token).
- Leader opens an activity's check-in screen → scans → marks present; or **GPS self-check-in** within a geofence radius of the activity location.
- Records method (qr/gps/manual), timestamp, and (for GPS) coordinates.
- Awards XP automatically on `present` via the gamification service.
- Live attendance counter over WebSocket.

### 10.5 Badges & Skills
- Badge catalog with `criteria_json` → gamification service can **auto-unlock** when criteria met (e.g. "attend 5 camps").
- Manual award by leaders also supported (logged in `member_badges.awarded_by`).
- Skills tracked with levels, assessed by leaders, feeding the member profile.

### 10.6 Training Hub
- Lessons with rich content + media, grouped by category and ordered.
- Optional quiz per lesson; passing awards XP + can unlock a badge.
- Per-member progress tracking (`/training/me/progress`).

### 10.7 Gamification
- Central `gamification_service`:
  - `AwardXP(memberID, source, points, refID)` → writes `xp_events`, updates `member.xp_total`, recomputes `level`, propagates to `unit.score_total`.
  - Level curve defined in config (e.g. level = f(xp)).
  - Example rules: attend activity `+10`, lead a team `+50`, pass quiz `+pass_score`, earn badge `+badge.xp_reward`.
- Leaderboards (members + units), filterable, with weekly/all-time views.
- Awards: "Future Leader" (وسام قائد المستقبل), "Most Improved", "Best Attendance" — computed from data.

### 10.8 Communication
- **Announcements**: targeted to all / a unit / leaders; pinnable.
- **Chat**: WebSocket channels (unit channels + groups + direct). REST endpoint for history (paginated).
- **Notifications**: in-app + Web Push (PWA). Triggered by: new activity, time change, announcement, reminder (job: "tomorrow's activity at 4 PM").

### 10.9 Reports & Dashboard
- **Dashboard KPIs**: member count, attendance %, top unit, at-risk members (high absence), trend deltas.
- **Reports**: per-member, per-unit, monthly troop report; export CSV/PDF.
- **Alerts**: member absent too often; member improving (positive trend) — surfaced on dashboard.

### 10.10 Equipment (v1.1)
- Inventory of gear with available/total counts; loan & return tracking tied to activities.

---

## 11. Environment Variables (.env)

> **Everything configurable lives in `.env`** — no hardcoded secrets, URLs, or tunables. Commit a `.env.example` with placeholder values; never commit the real `.env`.

### Backend `.env`
```env
# --- App ---
APP_NAME=Kashfi
APP_ENV=production            # development | production
APP_PORT=8080
APP_URL=https://scout.example.com
APP_TIMEZONE=Asia/Beirut

# --- Database (MySQL) ---
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=kashfi
DB_USER=kashfi_user
DB_PASSWORD=change_me
DB_PARAMS=charset=utf8mb4&parseTime=True&loc=Local
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=300      # seconds

# --- Auth / JWT ---
JWT_SECRET=change_me_long_random
JWT_ACCESS_TTL=900            # seconds (15 min)
JWT_REFRESH_TTL=604800        # seconds (7 days)
BCRYPT_COST=12

# --- CORS ---
CORS_ALLOWED_ORIGINS=https://scout.example.com

# --- Uploads ---
UPLOAD_DIR=./uploads
UPLOAD_MAX_SIZE_MB=15
UPLOAD_ALLOWED_TYPES=image/jpeg,image/png,image/webp,video/mp4
PUBLIC_UPLOAD_PATH=/uploads

# --- Attendance / Geofence ---
GEOFENCE_RADIUS_METERS=150
QR_SIGNING_SECRET=change_me

# --- Gamification ---
XP_PER_ATTENDANCE=10
XP_PER_QUIZ_PASS=20
XP_PER_LEADERSHIP=50
LEVEL_BASE_XP=100            # xp needed for level 1→2, scaled per curve

# --- Notifications (Web Push - VAPID) ---
VAPID_PUBLIC_KEY=
VAPID_PRIVATE_KEY=
VAPID_SUBJECT=mailto:admin@example.com

# --- Email (optional, reminders/reports) ---
SMTP_HOST=
SMTP_PORT=587
SMTP_USER=
SMTP_PASSWORD=
SMTP_FROM=no-reply@example.com

# --- Jobs ---
REMINDER_LOOKAHEAD_HOURS=24
MONTHLY_REPORT_DAY=1
```

### Frontend `.env` (Vite — only `VITE_` vars are exposed to the client)
```env
VITE_API_BASE_URL=https://scout.example.com/api/v1
VITE_WS_URL=wss://scout.example.com/ws
VITE_APP_NAME=Kashfi
VITE_DEFAULT_LOCALE=ar
VITE_VAPID_PUBLIC_KEY=
VITE_UPLOADS_BASE_URL=https://scout.example.com/uploads
```

> Load backend env via a typed `config.Config` struct (godotenv + env parsing) at startup; **fail fast** if a required var is missing.

---

## 12. Security

- **Passwords**: bcrypt (cost from env). Never log secrets or medical data.
- **JWT**: short-lived access + rotating refresh; store refresh server-side or use rotation + reuse detection.
- **RBAC**: enforce in middleware **and** service scoping (defense in depth). A leader can only read/write members in their units.
- **Medical data**: separate table, restricted endpoints, audit reads if possible (sensitive: this is minors' data — handle with care).
- **Validation**: every input validated server-side (validator) regardless of client-side Zod.
- **Uploads**: validate MIME + size + extension; store with random filenames; never trust client filename; serve from a non-executable path.
- **QR check-in**: signed tokens with expiry so QR codes can't be forged or replayed.
- **Geofence**: server verifies check-in coordinates against activity location + radius.
- **Rate limiting** on auth + check-in endpoints.
- **CORS** locked to known origins (from env).
- **SQL**: GORM parameterization only — no string-built queries.
- **Headers**: secure headers middleware (HSTS, X-Content-Type-Options, etc.).
- Since this involves **minors' personal and medical data**, restrict access tightly, minimize what's exposed to `member`/`parent` roles, and keep backups encrypted.

---

## 13. Implementation Order (Phased)

Built to be handed to an AI coding assistant **phase by phase**.

**Phase 0 — Foundations**
1. Repo scaffolding (backend + frontend folders as above).
2. Backend: config loader from `.env`, GORM connection, base Echo server, middleware (CORS, logger, recover), standard response envelope.
3. Frontend: Vite + TS + Tailwind v3 + React Query + Router + axios instance + i18n/RTL shell + design tokens.

**Phase 1 — Auth & RBAC**
4. User model, auth service (login/refresh/logout/me), JWT + bcrypt.
5. Auth middleware + `RequireRole`. Frontend: login page (RHF+Zod), auth context, ProtectedRoute/RoleRoute.

**Phase 2 — Members & Units**
6. Member + medical models/services/handlers; Lebanese phone normalization.
7. Unit + unit_leaders + unit_members.
8. Frontend: members list/profile/create-edit, medical (restricted), units list/detail.

**Phase 3 — Activities & Attendance**
9. Activity CRUD + media upload + feedback.
10. Attendance: manual + QR signing + GPS geofence + WS live counter.
11. Frontend: activity calendar/detail, attendance check-in screen (camera QR), feedback.

**Phase 4 — Gamification, Badges, Skills**
12. `gamification_service` (AwardXP, levels), `xp_events`.
13. Badges (catalog + auto-unlock + manual award), skills.
14. Frontend: badge wall, skill levels, member/unit leaderboards, XP rings.

**Phase 5 — Training Hub**
15. Lessons + quizzes + attempts + progress; XP/badge on pass.
16. Frontend: lesson list/reader, quiz flow, progress.

**Phase 6 — Communication**
17. Announcements; chat (WS hub + channels + history); notifications + Web Push.
18. Frontend: announcements, chat UI, notification center, PWA push opt-in.

**Phase 7 — Reports & Dashboard**
19. Dashboard KPIs + at-risk/improving alerts.
20. Reports (member/unit/monthly) + CSV/PDF export + scheduled monthly report job + activity reminders.

**Phase 8 — Hardening & Deploy**
21. Validation pass, rate limiting, secure headers, audit on sensitive reads.
22. Versioned migrations for prod, seed data, build SPA (embed in Go binary or Nginx), deploy on VPS (aaPanel), backups.

**Phase 9 (optional) — Equipment**
23. Equipment inventory + loans.

---

## 14. Future / WOW Features

- **AI activity suggestions** — recommend activities tailored to a unit's age/section and recent history.
- **Auto leader-evaluation drafts** — generate a draft evaluation per member from attendance, XP trend, quiz results, and badges (leader reviews/edits before saving).
- **Weekly missions** — auto-generated weekly tasks per unit with bonus XP.
- **"Future Leader" award** — algorithmic detection of leadership potential from participation + skill + peer signals.
- **Parent portal** — read-only view for parents (attendance, announcements, medical-on-file confirmation).
- **Camp map** — interactive map of camp sites + GPS tracking during hikes.
- **Native mobile app** — React Native (Expo) reusing the same API, with offline check-in sync.
- **Offline-first attendance** — queue check-ins offline (PWA), sync when back online.

---

*End of blueprint. Build phase by phase; keep handlers thin, services smart, and every secret in `.env`.*
