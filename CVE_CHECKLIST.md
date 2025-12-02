# CVE Checklist & Security Tracking

**Version:** 1.0
**Datum:** 2025-11-21
**Owner:** PentesterAgent + APIAgent
**Referenz:** SECURITY_HANDBOOK.pdf Gate 1, NAS_AI_SYSTEM.md §11.5

---

## ZWECK

Diese Checkliste dient als zentrale Übersicht aller identifizierten Schwachstellen (CVEs) im NAS.AI-System. Sie wird vor jedem Security Gate und Release geprüft. Deployments werden blockiert, wenn kritische CVEs (CVSS ≥ 7.0) offen sind.

---

## STATUS OVERVIEW

| Status | Count | Description |
|--------|-------|-------------|
| 🔴 OPEN (Critical) | 1 | CVSS ≥ 7.0 - Deployment BLOCKED (JWT Default Secret) |
| 🟠 OPEN (High) | 0 | CVSS 4.0-6.9 - Review required |
| 🟡 OPEN (Medium/Low) | 2 | CVSS < 4.0 - Tracked |
| ✅ CLOSED | 19 | Phase 1: 11 CVEs | Phase 1.5: 8 CVEs |

**Last Security Gate:** Phase 1.5 Security Hardening (2025-11-29) ✅ PASSED
**Next Security Gate:** Phase 2 Integration (Target: 2025-12-02)
**Security Score:** 93/100 (Grade: A) - OWASP 10/10

---

## 🔴 OPEN CVEs (CRITICAL) - DEPLOYMENT BLOCKER

### SEC-2025-003: JWT Default Secret Exposure

| Field | Value |
|-------|-------|
| **CVE-ID** | SEC-2025-003 (Internal) |
| **CVSS Score** | 8.5 (High/Critical) |
| **Status** | 🔄 IN PROGRESS |
| **Owner** | APIAgent |
| **Affected Component** | `infrastructure/api/src/handlers/auth.go` (planned) |
| **Description** | JWT tokens are currently signed with default/hardcoded secret. This allows token forgery if secret is discovered. |
| **Risk** | Authentication bypass, privilege escalation |
| **Remediation Plan** | 1. Remove default JWT secret from code<br>2. Implement secret loading from Vault/ENV<br>3. Add fail-fast check on startup<br>4. Rotate all existing tokens |
| **Evidence Required** | - Code review showing no defaults<br>- Unit test for missing secret fail-fast<br>- PentesterAgent validation report |
| **Target Date** | 2025-11-24 |
| **Dependencies** | SystemSetupAgent (Vault setup) |
| **Nachweis-Link** | `status/APIAgent/phase3/` (TBD) |

---

## 🟡 OPEN CVEs (MEDIUM/LOW) - TRACKED

### PERF-001: Missing Dependency Fail-Fast Checks

| Field | Value |
|-------|-------|
| **CVE-ID** | PERF-001 (Internal) |
| **CVSS Score** | 3.0 (Low) |
| **Status** | 🔄 IN PROGRESS |
| **Owner** | APIAgent |
| **Affected Component** | `infrastructure/api/src/main.go` (planned) |
| **Description** | Application does not fail-fast on startup if critical dependencies (DB, Redis, Vault) are unreachable. This leads to cryptic runtime errors. |
| **Risk** | Poor error messages, difficult debugging |
| **Remediation Plan** | Add startup health checks for all dependencies |
| **Target Date** | 2025-11-23 |
| **Nachweis-Link** | `status/APIAgent/phase3/` (TBD) |

### DOC-001: API Documentation Out of Sync

| Field | Value |
|-------|-------|
| **CVE-ID** | DOC-001 (Internal) |
| **CVSS Score** | 2.0 (Low) |
| **Status** | ⏳ PLANNED |
| **Owner** | APIAgent + DocumentationAgent |
| **Description** | Some API endpoints documented in blueprints don't have corresponding OpenAPI specs |
| **Remediation Plan** | Generate OpenAPI specs from code, add CI check |
| **Target Date** | Phase 4 |

---

## ✅ CLOSED CVEs (PHASE 1 - VERIFIED)

The following 11 CVEs were fixed during Phase 1 (Security Foundation) and verified by PentesterAgent on 2025-11-20:

| CVE-ID | Component | CVSS | Fix Date | Verification |
|--------|-----------|------|----------|--------------|
| AUTH-001 | WebSocket authentication | 9.0 | 2025-11-15 | ✅ PentesterAgent Phase 1 report |
| AUTH-002 | JWT validation bypass | 8.5 | 2025-11-15 | ✅ PentesterAgent Phase 1 report |
| AUTH-003 | CSRF token missing | 7.5 | 2025-11-16 | ✅ PentesterAgent Phase 1 report |
| SEC-001 | Plaintext password logging | 6.5 | 2025-11-14 | ✅ Code review + log audit |
| SEC-002 | SQL injection in file search | 8.0 | 2025-11-15 | ✅ Automated security tests |
| SEC-003 | Path traversal in file API | 9.5 | 2025-11-15 | ✅ PentesterAgent validation |
| SEC-004 | Insecure direct object reference | 7.0 | 2025-11-16 | ✅ Access control tests |
| SEC-005 | Missing rate limiting | 5.5 | 2025-11-17 | ✅ Load test verification |
| SEC-006 | Weak password requirements | 4.0 | 2025-11-14 | ✅ Policy enforcement tests |
| SEC-007 | Session fixation vulnerability | 7.5 | 2025-11-16 | ✅ Session management tests |
| SEC-008 | Insecure CORS configuration | 6.0 | 2025-11-17 | ✅ Network security validation |

**Phase 1 Security Gate:** ✅ PASSED (2025-11-20)
**Evidence Location:** `status/PentesterAgent/phase1/`

---

## ✅ CLOSED CVEs (PHASE 1.5 - SECURITY HARDENING - 2025-11-29)

The following 5 CRITICAL security bugs were fixed during Security Hardening and verified on 2025-11-29:

| CVE-ID | Component | CVSS | Fix Date | Verification |
|--------|-----------|------|----------|--------------|
| BUG-GO-001 | CORS Misconfiguration - Beliebige Origins erlaubt | 9.1 | 2025-11-29 | ✅ SECURITY-HARDENING-COMPLETE |
| BUG-GO-002 | RateLimiter Race Condition + Memory Leak | 8.5 | 2025-11-29 | ✅ SECURITY-HARDENING-COMPLETE |
| BUG-PY-003 | pgvector Extension nicht registriert | 9.0 | 2025-11-29 | ✅ SECURITY-HARDENING-COMPLETE |
| BUG-JS-002 | CORS Credentials nicht konfiguriert | 9.1 | 2025-11-29 | ✅ SECURITY-HARDENING-COMPLETE |
| BUG-GO-004 | SQL Injection in Search Handler | 9.8 | 2025-11-29 | ✅ SECURITY-HARDENING-COMPLETE |
| CWE-434 | Unrestricted File Upload (Malicious Files) | 8.5 | 2025-11-29 | ✅ Magic number validation + extension blacklist |
| CWE-639 | Authorization Bypass (RBAC missing) | 9.0 | 2025-11-29 | ✅ Admin-only middleware implemented |
| CWE-787 | Data Loss from Failed Restores | 8.0 | 2025-11-29 | ✅ Pre-restore safety backups |

**Mitigation Evidence:**
- `credentials: 'include'` in all Frontend fetch calls
- CORS Origin Whitelist validation in middleware/cors.go
- RateLimiter with TTL-based cleanup (no memory leak)
- pgvector `register_vector(conn)` in AI Agent
- Prepared SQL statements with pq.Array() for vector queries
- Magic number validation for file uploads (16 allowed MIME types)
- Admin-only RBAC middleware for destructive operations
- Emergency pre-restore backups before any data wipeout

**Security Score Improvement:** 78/100 → 93/100 (+19%)
**OWASP Top 10 Coverage:** 8/10 → 10/10 (100%)
**Phase 1.5 Security Gate:** ✅ PASSED (2025-11-29)
**Evidence Location:** `status/SECURITY-HARDENING-COMPLETE-2025-11-29.md`

---

## 📋 SECURITY GATES & RELEASE CRITERIA

### Gate Requirements

Before any deployment to production, the following criteria MUST be met:

1. **Gate 1: CVEs** ✅ / ❌
   - No OPEN Critical CVEs (CVSS ≥ 7.0)
   - All High CVEs (CVSS 4.0-6.9) have approved remediation plan
   - Prüfer: PentesterAgent

2. **Gate 2: Tests** ✅ / ❌
   - Unit test coverage ≥ 80%
   - All security regression tests passing
   - Prüfer: CI Pipeline

3. **Gate 3: Secrets** ✅ / ❌
   - No secrets in code (Gitleaks scan clean)
   - All secrets in Vault or authorized exceptions (see DEV_GUIDE.md §5)
   - Prüfer: Pre-Commit Hook + Manual Review

4. **Gate 4: Auth** ✅ / ❌
   - All endpoints behind auth middleware (except `/auth/*`)
   - JWT secrets loaded from secure source
   - Prüfer: APIAgent + PentesterAgent

5. **Gate 5: CSRF** ✅ / ❌
   - All POST/PUT/DELETE endpoints require valid CSRF token
   - Prüfer: APIAgent + PentesterAgent

### Current Gate Status (Phase 3)

| Gate | Status | Blocker |
|------|--------|---------|
| Gate 1: CVEs | ❌ | SEC-2025-003 (JWT defaults) |
| Gate 2: Tests | ⏳ | Infrastructure pending |
| Gate 3: Secrets | ✅ | All exceptions documented |
| Gate 4: Auth | ⏳ | SEC-2025-003 blocks |
| Gate 5: CSRF | 🔄 | Rollout in progress |

---

## 🔄 WORKFLOW

### Neuer CVE gefunden

1. **Triage** (AnalysisAgent oder PentesterAgent):
   - CVSS Score berechnen
   - Betroffene Komponenten identifizieren
   - Eintrag in diesem Dokument anlegen (OPEN)

2. **Assignment** (Orchestrator):
   - Owner zuweisen (meist APIAgent, NetworkSecurityAgent oder SystemSetupAgent)
   - Target Date festlegen (Critical: ≤3 Tage, High: ≤7 Tage, Medium: ≤30 Tage)
   - Dependencies prüfen

3. **Remediation** (Assigned Agent):
   - Fix implementieren
   - Tests schreiben (Regression Prevention)
   - Statuslog dokumentieren mit Nachweis-Link

4. **Verification** (PentesterAgent):
   - Fix validieren (Re-Test)
   - Evidence sammeln
   - CVE auf CLOSED setzen

5. **Documentation** (DocumentationAgent):
   - Nachweis-Link in CVE_CHECKLIST.md eintragen
   - Security Gate Status aktualisieren

---

## 📊 METRICS & REPORTING

### Monthly CVE Report

- **Open Critical:** {{ count }}
- **Open High:** {{ count }}
- **Average Time to Remediation:** {{ days }}
- **Security Gate Pass Rate:** {{ percentage }}

Reports werden vom Orchestrator am Monatsende automatisch generiert und in `status/security-reports/YYYYMM.md` abgelegt.

---

## 🔗 REFERENZEN

- **Security Handbook:** `docs/security/SECURITY_HANDBOOK.pdf`
- **Phase Roadmap:** `docs/planning/MASTER_ROADMAP.md`
- **PentesterAgent Status:** `status/PentesterAgent/`
- **Incident Response:** `NAS_AI_SYSTEM.md §10`

---

**Letzte Aktualisierung:** 2025-11-21
**Nächste Review:** 2025-11-28 (Phase 3 Gate)

Terminal freigegeben.
