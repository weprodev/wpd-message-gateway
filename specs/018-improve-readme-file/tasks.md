# Improve Readme.md Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Improve the `Readme.md` file to make it more welcoming, informative, and clearer for new users by reorganizing badges, adding a Key Features section, and simplifying the quick start.

**Architecture:** We will surgically update `Readme.md` sections, prioritizing the recommended Docker setup for the HTTP server, structuring the features cleanly, and clarifying the configuration.

**Tech Stack:** Markdown

---

### Task 1: Reorganize Introduction and Badges

**Files:**
- Modify: `Readme.md`

- [x] **Step 1: Group badges and update intro**
Update the top of `Readme.md` so the badges are organized better. Update the introduction text to explicitly separate "Embedded Go SDK" and "Standalone HTTP Gateway".

- [x] **Step 2: Add Key Features section**
Add a `## Key Features` section right after the introduction with a bulleted list of main capabilities (e.g., Unified interface, Provider abstraction, DB-first config, Workspace isolation).

- [x] **Step 3: Enhance Visuals**
Enhance or replace the "Two Ways to Use" diagram in `Readme.md` to make it more visually engaging and clear about the two usage modes (SDK vs. Server).

- [x] **Step 4: Commit changes**
```bash
git add Readme.md
git commit -m "docs: reorganize badges, update intro, add key features and enhance visuals"
```

### Task 2: Update Quick Start and Configuration

**Files:**
- Modify: `Readme.md`

- [x] **Step 1: Simplify HTTP Server Quick Start**
Under `### Option B: HTTP Server`, make the Docker (`make dev`) path the primary getting-started method. Move the manual setup (`make install`, `make start`) to an "Alternative: Manual Setup" subsection.

- [x] **Step 2: Validate Quick Start Instructions**
**Verify and validate all quick start instructions** (including `make dev` and manual setup) to ensure they are accurate and work as described.

- [x] **Step 3: Clarify Configuration**
In the `## Configuration` section, clearly distinguish between `configs/local.yml` (used only for server/DB configs) and the Portal UI (used for messaging provider configurations/API keys).

- [x] **Step 4: Validate Links**
Validate that all links within the `Readme.md` are valid after modifications.

- [x] **Step 5: Commit changes**
```bash
git add Readme.md
git commit -m "docs: simplify quick start, clarify configuration and validate instructions"
```
