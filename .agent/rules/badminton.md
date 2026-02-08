---
trigger: always_on
---

## 1. Context & Constraints

- **Project:** `badminton_tournament`
- **Role:** Senior Full-stack Engineer (Go/Vue).
- **Deadline:** 1 Day.
- **Goal:** A fully functional tournament system with Google Form sync, GSL Brackets, and Video handling.

## 2. Design System: "Tech-Flat Outfit"

**Strictly adhere to these UI rules:**

- **Typography:** Global font family: `'Outfit', sans-serif`.
- **Visual Style:** Flat design. **NO shadows**.
- **Shapes:** `rounded-sm` (2px border-radius) for cards, buttons, inputs.
- **Colors:**
  - Bg: `#FDFBFF` (Ultra-light Purple).
  - Border: `#E9D5FF` (Purple-200).
  - Active/Brand: `#7C3AED` (Violet-600).
  - Text: `#4C1D95` (Purple-900).

## 3. Detailed Features & Logic

### A. Google Form Integration (Webhook)

- **Endpoint:** `POST /api/webhooks/form`
- **Expected Payload:**
  ```json
  {
    "email": "user@example.com",
    "name": "Nguyen Van A",
    "pool": "Mesoneer" // or "Lab"
  }
  ```
- **Backend Logic (Upsert):**
  - Check DB: `SELECT * FROM participants WHERE email = ?`
  - If exists: UPDATE pool/name.
  - If not exists: INSERT new record.

### B. Team Formation (Pool Constraint)

- **Logic:**
  - Input: `pool_name` (Enum: 'Mesoneer', 'Lab').
  - Process: Fetch participants by Pool -> Random Shuffle -> Pair (1 & 2, 3 & 4...).
  - Output: Create `Team` records in DB.
- **Validation:** A team CANNOT contain participants from different pools.

### C. GSL Bracket Logic (The Core)

- **Structure:** Each Group has exactly 5 Matches.
  1.  **M1:** T1 vs T2. (Winner->M3, Loser->M4).
  2.  **M2:** T3 vs T4. (Winner->M3, Loser->M4).
  3.  **M3 (Winners):** Win(M1) vs Win(M2). (Winner->Qualified #1, Loser->M5).
  4.  **M4 (Losers):** Lose(M1) vs Lose(M2). (Winner->M5, Loser->Eliminated).
  5.  **M5 (Decider):** Lose(M3) vs Win(M4). (Winner->Qualified #2).
- **Automation:**
  - Create `Service.AdvanceTeam(matchID, winnerID)`: Updates the `TeamA_ID` or `TeamB_ID` of the connected _next match_.

### D. Match Execution & Media

- **Data:** `Score` (String, e.g., "21-19, 21-20") and `VideoURL` (String).
- **UI:** Admin modal to input score and paste YouTube link.
- **Public:** If `VideoURL` is present, render a "Play" button on the match card.

## 4. Database Schema (Neon Postgres)

- `participants`: `id, email (unique), name, pool (varchar), created_at`
- `teams`: `id, player1_id, player2_id, pool`
- `tournaments`: `id, name, status`
- `groups`: `id, tournament_id, name (A, B...)`
- `matches`:
  - `id, group_id`
  - `team_a_id, team_b_id` (Nullable - waiting for advance)
  - `score, video_url`
  - `next_match_win_id` (FK to matches)
  - `next_match_lose_id` (FK to matches)

## 5. Execution Roadmap (Step-by-Step)

### Step 1: Infrastructure & Auth

- Generate `backend/Dockerfile` & `render.yaml`.
- Implement JWT Login (Admin Password from Env).
- **Middleware:** Protect `/api/teams/*` and `/api/matches/*`.

### Step 2: Backend Core (The Engine)

- Implement `ParticipantService` (Upsert logic).
- Implement `BracketService` (GSL Generation & Auto-Advance).
- **Unit Test:** Write a test for `GenerateBracket` to ensure connection logic (M1->M3) is correct.

### Step 3: Admin UI (Vue)

- **Login View:** Simple flat card.
- **Dashboard:**
  - Tabs: [Mesoneer] [Lab].
  - Action: [Sync Google Form] (Button triggers webhook test or refresh).
  - Action: [Generate Teams] -> Show list.
  - Bracket View: Admin can click any match to Update Result.

### Step 4: Public UI (Vue)

- **Layout:** Tech-Flat (Outfit Font).
- **Components:**
  - `BracketGrid`: Display the 5 matches layout.
  - `VideoModal`: Embed YouTube iframe.
  - `AutoRefresher`: Use `setInterval` to fetch bracket data.

## 6. Immediate Action

1.  Generate the **Folder Structure**.
2.  Generate the **Go Models** (structs) matching the GSL logic above.
3.  Create the **Google Form Webhook Handler** code.
