## 1.  Summary
* **Project Name:** Bookmark Manager (Working Title)
* **Core Problem:** Users struggle to retain, categorize, and recall useful web links, often losing them in disorganized browser bookmarks or chaotic messaging apps.
* **Target Audience:** Content researchers, developers, students, and digital hoarders.
* **The Solution:** A lightweight, highly organized link repository where users can save URLs, assign them to custom categories, and append personal Markdown notes.

---

## 2. MVP Feature Matrix
To avoid scope creep, development is prioritized into core items required for launch versus future features.

| Feature Area | MVP (Scope) | Future Enhancements (Out of Scope) |
| :--- | :--- | :--- |
| **Authentication** | Email / Password sign-up and login | Google/GitHub OAuth, Magic Links |
| **Link Management** | CRUD (Create, Read, Update, Delete) URLs | Automatic page title & metadata scraping |
| **Organization** | Custom text categories / folders | Multi-tagging, Nested categories |
| **Notes** | Plain-text / Markdown notes per link | Rich-text WYSIWYG editor, PDF attachments |
| **Search & Filter**| Filter by category | Full-text search across note contents |
| **Sharing** | Private personal dashboard only | Public shareable boards, Team collaboration |

---

## 3. User Journey & UI Blueprint

### Core User Flows
1. **Onboarding:** User lands on landing page -> Clicks Register -> Fills form -> Redirected to Dashboard.
2. **Saving a Link:** Dashboard -> Click "Add Link" -> Paste URL, select Category, type Note -> Click Save.
3. **Filtering Content:** Dashboard -> Click Category Sidebar Item -> View filtered list of links.

### Planned App Routes
* `GET /` - Public marketing landing page explaining the utility.
* `GET /auth/signup` & `/auth/login` - Clean onboarding interfaces.
* `GET /dashboard` - Core view showing all saved links, category navigation, and add-link modal.
* `GET /categories` - Page to manage, edit, or delete existing categories.

---

## 4. Technical Stack Architecture

```
       [ Frontend UI: React ]
                     │  ▲
        REST API     │  │  JSON Response
     (JSON/Bearer)   ▼  │
       [ Backend API Framework: Go ]
                     │  ▲
        SQL Queries  │  │  Data Records
                     ▼  │
       [ Relational Database: PostgreSQL ]
```

* **Frontend:** **React** with **Tailwind CSS**
* **Backend:** **Golang (Go)**
* **Database:** **PostgreSQL**
* **Data Model Preview:**
  * **Users Table:** `id` (UUID), `email` (Unique), `password_hash`, `created_at`
  * **Categories Table:** `id` (UUID), `user_id` (FK), `name`, `created_at`
  * **Links Table:** `id` (UUID), `user_id` (FK), `category_id` (FK, Nullable), `url`, `notes` (Text), `created_at`

---

## 5. Production Timeline

### Phase 1: Foundation (Weeks 1 - 2)
* Setup local environments, GitHub repositories, and initialize DB schemas.
* Build backend authentication endpoints (`/signup`, `/login`) with secure JWT/session handling.
* Wire up basic UI frontend layout shell (Sidebar + Main panel).

### Phase 2: Core Engine Build (Weeks 3 - 5)
* Implement CRUD API logic for Categories and Links.
* Build the "Add Link" modal UI and integrate backend validation (checking for valid URL formats).
* Implement the interactive filtering mechanism on the frontend dashboard.

### Phase 3: Refinement & Testing (Weeks 6 - 7)
* Add local validation checks and unit testing scripts for API endpoints.
* Style the interface using Tailwind utility classes for a clean, modern aesthetic.
* Ensure responsive mobile views for saving links on the go.

### Phase 4: CI/CD & Launch (Week 8)
* Setup automated deployment tracking via GitHub Actions.
* Deploy production infrastructure onto target cloud environments (e.g., Vercel / Render).
* Execute final user verification checks.