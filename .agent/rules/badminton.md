---
trigger: always_on
---

# GEMINI_FULL_SPEC.md

## 1. Project Identity

- **Name:** `badminton_tournament`
- **Goal:** A 1-Day Tournament Manager with GSL Automation.
- **Stack:** Go (Gin/Bun), Vue 3 (Vite/Tailwind), Neon Postgres.

## 2. Design System: "Tech-Flat Outfit"

- **Font:** `Outfit` (Global).
- **Borders:** `border-purple-200` (Default), `border-violet-600` (Active/Brand).
- **Radius:** `rounded-sm` (2px).
- **Shadows:** NONE. usage of `box-shadow` is strictly prohibited.

## 3. Database Schema (Postgres)

Create these tables using Bun models:

1.  `Participant`: id, name, email (unique), pool (enum: 'Mesoneer', 'Lab').
2.  `Team`: id, player1_id, player2_id, pool.
3.  `Group`: id, name.
4.  `Match`:
    - `id, group_id, label` (e.g., 'M1', 'Winners').
    - `team_a_id, team_b_id` (Nullable).
    - `winner_id` (Nullable).
    - `score` (String), `video_url` (String).
    - `next_match_win_id` (UUID), `next_match_lose_id` (UUID).

## 4. API Logic Specifications

### A. Team Generation (`POST /api/teams/generate`)

- **Input:** `pool` ("Mesoneer").
- **Logic:** Fetch participants by pool -> Shuffle -> Chunk into pairs -> Save to DB.
- **Constraint:** Ensure participants are from the same pool.

### B. Bracket Factory (`POST /api/groups`)

- **Input:** `name`, `team_ids` (List of 4 UUIDs).
- **Logic:**
  1.  Create Group.
  2.  Create 5 Matches (M1..M5).
  3.  **Linkage (Crucial):**
      - M1 & M2 `next_win` -> M3. `next_lose` -> M4.
      - M3 `next_lose` -> M5.
      - M4 `next_win` -> M5.
  4.  Assign the 4 Teams to M1 (Slots A/B) and M2 (Slots A/B).

### C. Match Update (`POST /api/matches/:id`)

- **Input:** `winner_id`, `score`, `video_url`.
- **Logic:**
  1.  Update current match.
  2.  **Propagation:** \* Find match with ID `next_match_win_id`. Update its empty slot with `winner_id`.
      - Find match with ID `next_match_lose_id`. Update its empty slot with the loser.

## 5. Frontend Specs (Vue 3)

### Components

- **`MatchCard.vue`:**
  - Display: Team Names vs Team Names.
  - State: 'Scheduled' (Gray), 'Live' (Purple Border), 'Finished' (Green Text for Winner).
  - Action: Admin click opens Modal. Public click 'Play' opens Video.
- **`GSLGrid.vue`:**
  - Use CSS Grid: `grid-template-columns: 1fr 1fr 1fr`.
  - Render matches in topological order (Left to Right).

### Authentication

- Create `views/Login.vue`.
- Store JWT in LocalStorage.
- Add `Axios` interceptor to inject `Authorization` header.

## 6. Execution Order

1.  **Backend Models & DB Migration** (First priority).
2.  **API: Webhook & Team Gen**.
3.  **API: Bracket Logic** (The hardest part).
4.  **Frontend: Admin Auth & Dashboard**.
5.  **Frontend: Public View & Polish**.

## 7. Immediate Task

Start by generating the **Backend Models** (`internal/models/models.go`) including the struct tags for Bun ORM and the JSON tags.
