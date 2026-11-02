# Security Documentation & Vulnerability Assessment

## Overview

Medix Backend implements security best practices across authentication, authorization, data validation, and asset management. This document outlines security architecture and addresses key vulnerability considerations.

---

## 🔒 Top 3 Security Considerations & Mitigations

### 1. JWT Secret Strength & Environment Exposure
- **Risk / Vulnerability:** If `JWT_SECRET` uses a weak string or hardcoded default in configuration (`config/config.yaml`), tokens can be forged by unauthorized actors to gain administrator access.
- **System Mitigation:**
  - The authentication middleware (`internal/middleware/AuthMiddleware.go`) reads `JWT_SECRET` from environment variables (`os.Getenv("JWT_SECRET")`) or secure configuration.
  - Tokens have strict expiration times and support session revocation tracking.
  - **Production Recommendation:** Use a random secret key of at least 32 characters and store it in production environment variables / secret managers (never commit to git).

### 2. File Upload Validation & Storage (Path Traversal / Arbitrary File Upload)
- **Risk / Vulnerability:** Profile photo upload endpoints (`POST /users/profile/photo`) accept multipart files. Without strict extension and filename validation, attackers could upload malicious scripts or perform path traversal.
- **System Mitigation:**
  - Strict file size limit enforcement (max 2MB, `file.Size > 2<<20`).
  - Strict file extension whitelist: `.jpg`, `.jpeg`, `.png`, `.webp`.
  - Secure server-side filename generation using the format `user_{id}_{timestamp}{ext}` (ignoring the user-supplied filename).
  - Files are stored in an isolated directory (`uploads/profiles`).

### 3. Insecure Direct Object Reference (IDOR) & Role Authorization
- **Risk / Vulnerability:** Users attempting to access or manipulate transaction data or profiles belonging to other users by guessing resource IDs (`/transactions/{id}`).
- **System Mitigation:**
  - Implementation of RBAC (`RequireRoles`) on sensitive routes (admin, cashier, owner).
  - Resource ownership checks in transaction handlers (`internal/transaction/handler/TransactionHandler.go`), validating that `transaction.IDUser == userID` before returning sensitive data or price histories.
  - Separation of public and protected routes to prevent access without a valid token.

---

## 🛡️ Best Practices & Standard Security Controls

- **Password Hashing:** Uses **Bcrypt** with secure cost factors for storing password hashes in the database.
- **SQL Injection Prevention:** Uses **GORM** ORM with automatic query parameterization to prevent SQL injection vulnerabilities.
- **CORS Policy:** Explicitly restricted (default: `http://localhost:5173` for Vite frontend).
- **Database Constraints:** Data integrity enforcement at the database level using `CHECK` constraints, foreign keys, and unique indexes.

## Reporting Security Issues
If you discover any security vulnerability in this project, please report it to the internal Medix development team.
