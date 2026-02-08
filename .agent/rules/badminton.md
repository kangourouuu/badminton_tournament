---
trigger: always_on
---

# GEMINI.md - Project: `badminton_tournament`

## 1. Context & Role

- **Project Name:** `badminton_tournament`
- **Role:** Senior Full-stack Engineer & DevOps.
- **Goal:** Build a Tournament Manager in 1 Day.
- **Stack:** \* **Backend:** Go (Gin, Bun, Postgres).
  - **Frontend:** Vue 3 (Vite, Tailwind, TypeScript).
  - **DB:** Neon.tech (Postgres) - Branched (Main/Dev).
  - **Deploy:** Docker (Render) & Vercel.

## 2. Infrastructure as Code (IaC) Requirements

Generate these files immediately upon request:

- **`backend/Dockerfile`:** Multi-stage (Go 1.22 -> Alpine). `CGO_ENABLED=0`.
- **`render.yaml`:** Define 2 services: `badminton-be-prod` (main branch) and `badminton-be-dev` (dev branch). Use `sync: false` for `DATABASE_URL`.
- **`frontend/vercel.json`:** SPA Rewrite rules.

## 3. Design System: "Tech-Flat Outfit" (Strict)

- **Font:** 'Outfit' (Google Font).
- **Shape:** `rounded-sm` (2px). Sharp edges.
- **Depth:** **NO SHADOWS**. Use `border-purple-200` to define structure.
- **Palette:** Background `#FDFBFF` (Pale Purple), Surface `#FFFFFF`, Accent `#7C3AED`.

## 4. Business Logic & Constraints

### A. Participants & Pools (CRITICAL)

- **Data Model:** `Participant` table must have a `pool` field (Enum/String: 'Mesoneer', 'Lab').
- **Input:** Webhook from Google Forms includes the Pool.

### B. Team Formation (The Wheel)

- **Constraint:** Teams must be formed **within the same Pool**.
  - _Logic:_ Admin selects "Spin for Mesoneer" -> System filters participants with `pool='Mesoneer'` -> Randomizes pairs -> Creates Teams.
  - _Logic:_ Same process for "Lab".
- **UI:** The Wheel component visually spins and outputs the pairs.

### C. Bracket & Scoring (GSL Format)

- **Generation:** Admin assigns Teams to Groups. System auto-generates 5 Matches per Group (GSL Structure).
- **Execution:**
  - **No Real-time:** Scores are entered only after the match ends (Final Result).
  - **Video:** Admin inputs `video_url` (YouTube link) along with the score.
- **Auto-Advance:** When Admin saves the Final Score -> System calculates Winner -> Updates the `team_id` in the next bracket slot automatically.

### D. Public Layer

- Read-only view.
- Shows the Bracket tree (GSL Layout).
- **Replay:** If `match.video_url` is not empty, show a "Play" icon (Outline style). Clicking it opens a Modal to watch.

## 5. Execution Roadmap

### Phase 1: Scaffolding (Infrastructure)

- Generate Monorepo structure: `backend/`, `frontend/`.
- Generate `Dockerfile`, `render.yaml`.
- Initialize Go Mod and Vue 3 project (with Tailwind/Outfit config).

### Phase 2: Backend Core (Go)

- **Models:** `Tournament`, `Participant` (Pool), `Team`, `Match` (VideoURL).
- **API:**
  - `POST /api/webhooks/form`: Ingest participants.
  - `POST /api/teams/generate`: Input `pool_name` -> Return random pairs.
  - `POST /api/matches/:id/result`: Input `{score_a, score_b, video_url}` -> Finalize match.

### Phase 3: Frontend UI (Vue)

- **Admin Dashboard:** \* Tabs for "Mesoneer" vs "Lab".
  - "Spin Wheel" button (filters by active tab).
  - Bracket Manager (Enter Score + Paste Video Link).
- **Public View:**
  - Flat Design Bracket.
  - Video Modal (`<iframe src="...">`).

### Phase 4: CI/CD

- Create `.github/workflows/pipeline.yml` (Test & Build Check).

## 6. Immediate Action

Start by generating the **Directory Structure**, **Dockerfile**, and **render.yaml** according to these specs so the user can initialize the repo.

# UI AUDIT & REFACTOR REQUEST

**Current Status:** I just checked the frontend, and it looks like raw HTML (Times New Roman, default buttons). The styling is completely missing.

**Task:** Refactor the `PublicView.vue` (and `App.vue` layout) IMMEDIATELY to match our **"Tech-Flat Outfit"** Design System defined in GEMINI.md.

**Specific Requirements to Implement Now:**

1.  **Global Font:** Ensure `font-family: 'Outfit', sans-serif` is applied to the `<body>` or root `div`.
2.  **Color Theme:**
    - Change Page Background to `#FDFBFF` (or Tailwind `bg-slate-50`).
    - Change Text Color to `text-slate-900` / `text-purple-900`.
3.  **The Header:**
    - Make it a proper Navbar: `bg-white`, `border-b border-purple-200`, `px-6 py-4`.
    - Title should be `font-bold text-xl text-purple-900`.
4.  **The Pool Buttons (Mesoneer/Lab):**
    - **DO NOT** use default buttons.
    - Style them as **Flat Tabs**:
      - **Active:** `bg-violet-600 text-white rounded-sm`.
      - **Inactive:** `bg-white border border-purple-200 text-purple-700 hover:bg-purple-50 rounded-sm`.
5.  **The Content:**
    - Wrap the "Matches are being arranged" text in a **Card Component**: `bg-white border border-purple-200 rounded-sm p-12 text-center`.
    - Center the card in a `max-w-5xl mx-auto` container.

**Action:** Rewrite the full code for the View component so I can copy-paste it and fix this ugly UI.
