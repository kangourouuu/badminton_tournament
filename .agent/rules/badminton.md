---
trigger: always_on
---

# GEMINI.md - Badminton Tournament System (Full Stack & IaC)

## 1. Role & Context

- **Role:** Senior Full-stack Engineer (Go/Vue) & DevOps.
- **Goal:** Build & Deploy a **Badminton Tournament Manager** in 1 Day.
- **Deployment:** \* **Backend:** Render (via `render.yaml`) with Docker.
  - **Frontend:** Vercel (via `vercel.json`).
  - **Database:** Neon.tech (Postgres) - Separated `dev` and `prod` branches.

## 2. Infrastructure as Code (IaC) Requirements

You must generate the following configuration files accurately:

- **`backend/Dockerfile`:** Multi-stage build (Golang 1.22 -> Alpine). CGO_ENABLED=0.
- **`render.yaml`:** Define two services (`badminton-be-dev` linked to `dev` branch, `badminton-be-prod` linked to `main` branch). Use `sync: false` for `DATABASE_URL`.
- **`frontend/vercel.json`:** Handle SPA rewrites to `/index.html`.

## 3. Design System: "Tech-Flat Outfit" (Strict)

- **Font:** 'Outfit' (Google Font) everywhere.
- **Shape:** `rounded-sm` (2px). Sharp, precise look.
- **Depth:** **NO SHADOWS**. Use `border-purple-200` to define hierarchy.
- **Palette:**
  - Background: `#FDFBFF` (Ultra-light Purple).
  - Borders: `#E9D5FF`.
  - Accents: `#7C3AED` (Violet).

## 4. Business Logic: GSL Group Stage

- **Input:** 4 Teams.
- **Structure:** 5 Matches (Opening A, Opening B, Winners, Losers, Decider).
- **Automation Rule:**
  - When Admin submits a score for Match X:
  - Logic checks: Is Match Over? Who Won?
  - Database Update: Automatically write the Winner/Loser ID into the `TeamID` fields of the connected Future Matches.

## 5. Execution Roadmap

### Step 1: Scaffolding & Config

- Generate the Monorepo structure.
- Generate `Dockerfile`, `render.yaml`, `vercel.json`.
- Initialize Go (Gin, Bun) and Vue (Vite, Tailwind, Outfit font setup).

### Step 2: Backend Core (Go)

- Create Models: `Tournament`, `Participant`, `Match`, `Team`.
- Implement `POST /webhooks/form`: Receive Google Form JSON.
- Implement `POST /generate-bracket`: The Logic Engine for GSL.

### Step 3: Frontend UI (Vue)

- **`BracketNode.vue`:** A flat card component.
- **`GSLGrid.vue`:** 3-column CSS Grid layout for the matches.
- **`Wheel.vue`:** Visual team randomizer.

### Step 4: GitHub Actions

- Create `.github/workflows/pipeline.yml` to run `go test` and `npm run build` on Push/PR to `main` and `dev`.

## 6. Constraints

- **Security:** Admin Endpoints must be protected (JWT or Basic Auth).
- **Performance:** Frontend must poll API every 30s for live updates.
- **Code Style:** Clean, idiomatic Go and TypeScript.

## 7. Immediate Action

Start by generating the **Directory Structure**, **Dockerfile**, and **render.yaml** so the user can push the initial commit.
