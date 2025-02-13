# Contributing Guidelines

Thank you for considering contributing to this project! 🎉  
This project follows **Conventional Commits** and **Semantic Versioning** to maintain a clear and structured commit history.

---

## 📌 Commit Message Format
Each commit message must follow this structure:

    <type>(<optional scope>): <short description>

    [optional body]

    [optional footer(s)]

### ✅ Examples:
    feat(parser): add support for regex rules
    fix(evaluator): prevent nil pointer panic
    docs(readme): update installation instructions

### 🚀 Allowed Commit Types
| Type      | Purpose |
|-----------|---------|
| **feat**  | Introduces a new feature |
| **fix**   | Fixes a bug |
| **docs**  | Documentation updates (e.g., README) |
| **chore** | Maintenance tasks (e.g., CI/CD updates) |
| **refactor** | Code improvements without changing behavior |
| **test**  | Adding/updating tests |

---

## 🏗️ Development Workflow

1. **Fork** the repository and create a new branch:
    git checkout -b feat/new-feature

2. **Make changes & commit** following the **Conventional Commits** format.

3. **Push your branch** and **open a Pull Request (PR)**.

---

## 📜 Versioning Strategy

This project follows **Semantic Versioning (`MAJOR.MINOR.PATCH`)**:
- **MAJOR** → Breaking changes 🚨
- **MINOR** → New features, backward-compatible ✅
- **PATCH** → Bug fixes 🐛

Releases are automated using **GitHub Actions** and the `.semver` file.

---

## 🔧 Setting Up the Project

### 1️⃣ Install Dependencies
Ensure you have **Go installed** (version 1.21+ recommended).

    go mod tidy

### 2️⃣ Run Tests
    go test ./...

---

## ✅ Code Style & Formatting

- Follow **Go best practices** and use ` + "`gofmt`" + ` before committing.
- Keep code **clean and modular**.

    gofmt -w .

---

## 🎯 Need Help?
If you're unsure about something, feel free to open an **Issue** or **Pull Request**.  
I'm happy to help! 🚀