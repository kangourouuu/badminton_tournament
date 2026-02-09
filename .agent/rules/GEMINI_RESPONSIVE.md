---
trigger: always_on
---

# GEMINI_RESPONSIVE_FEATURE.md

## 1. Context & Goal

- **Project:** `badminton_tournament`
- **Role:** Senior Frontend Engineer & UI/UX Designer.
- **Goal:** Implement a **Mobile-First Responsive Design** and create a **High-Fidelity Rules Page**.
- **Constraint:** Strictly adhere to "Tech-Flat Outfit" (No shadows, sharp borders, Purple/White theme).

## 2. Responsive Strategy (Mobile-First)

- **Breakpoints:**
  - `default`: Mobile (Vertical layout, 1 column).
  - `md` (768px): Tablet (Transition point).
  - `lg` (1024px): Desktop.
- **Navigation (Hidden Admin):**
  - **Public Navbar:** Shows Logo + "Bracket" + "Rules" + "Contact".
  - **CRITICAL:** **REMOVE the "Admin Login" button** from the Public Navbar.
  - **Access:** Admin must manually type `/admin` in the URL. Ensure the Router Guard redirects unauthenticated users from `/admin` to `/login`.

## 3. New Features Specification

### Feature A: Tournament Rules (The "Playbook" Concept)

- **Visual Concept:**
  - **Background:** Use a subtle geometric pattern (e.g., faint badminton court lines or repeating shuttlecocks) in `#F3E8FF` (very light purple).
  - **Structure:** Do not render a wall of text. Break content into **Topic Cards** (e.g., Scoring, Fouls, Equipment).
  - **Typography:** Use **Styled Numbering**. Example: A purple square `bg-purple-900` containing the number, followed by the rule title.
- **Backend:**
  - `GET /api/public/rules`: Fetch Markdown/HTML content.
- **Frontend (`RulesView.vue`):**
  - Render a list of `RuleCard` components.
  - Use `prose` (Tailwind Typography) but customized to match the Outfit font and purple headers.
  - **Decoration:** Add an abstract SVG graphic (Badminton theme) in the top-right corner of the header.

### Feature B: Contact Page

- **UI:** A centered "Business Card" style component.
  - **Content:** Organizer Name, Hotline, Email, Support Zalo/Telegram.
  - **Style:** Border `border-purple-200`, plain white background.

## 4. Implementation Roadmap

### Step 1: Responsive Layout & Navbar Refactor

1.  **Refactor `TheNavbar.vue`:**
    - **Remove** the `<router-link to="/login">`.
    - Add a Mobile Menu Toggle (Hamburger icon).
    - Ensure links (Bracket, Rules, Contact) are visible on Desktop and collapsible on Mobile.
2.  **Router Config:**
    - Verify that accessing `/admin` triggers the `beforeEnter` guard -> Redirect to `/login` if no token found.

### Step 2: Implement "The Playbook" (Rules Page)

1.  **Backend:** Create `internal/models/rule.go`.
2.  **Frontend:**
    - Create `views/RulesView.vue`.
    - Implement the **Geometric Background** (CSS or SVG).
    - Style the content using Cards:
      ```html
      <div
        class="border border-purple-200 bg-white p-6 rounded-sm relative overflow-hidden"
      >
        <h3 class="font-bold text-purple-900 flex items-center gap-2">
          <span
            class="bg-purple-900 text-white w-8 h-8 flex items-center justify-center rounded-sm text-sm"
            >01</span
          >
          General Rules
        </h3>
        <p class="mt-4 text-slate-600">...</p>
      </div>
      ```

### Step 3: Bracket Responsiveness

1.  Wrap the Bracket Grid in a container with `overflow-x-auto`.
2.  Set a `min-width` (e.g., `min-w-[900px]`) on the inner grid to prevent the GSL tree from breaking on small screens.

## 5. Immediate Action

Start by **Refactoring the Navbar** (removing the Login button and adding responsive toggle), then implement the **Rules Page with the new Design Concept**.
