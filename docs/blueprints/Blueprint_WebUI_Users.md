# NAS.AI – WebUI Users Blueprint

## 1. Scope & Responsibilities
- User management (`/users`): CRUD, roles, MFA states.

## 2. UX & Layout
```
┌─────────────────────────────────────────────────────────────┐
│ Header: Users                    [Search] [Filters] [New]   │
├────────────┬───────────────────────────────────────────────┤
│ Sidebar    │ Table (sortable)                              │
│ [← Back]   │ ┌───────────────────────────────────────────┐ │
│ [Home]     │ │ username | email | role | MFA | actions   │ │
│ Filters    │ │ row highlight for admin/self              │ │
│  Role      │ └───────────────────────────────────────────┘ │
│  Status    │ Drawer (right) for user detail/edit.          │
├────────────┴───────────────────────────────────────────────┤
│ Footer: total users • pending invites • audit log link     │
└─────────────────────────────────────────────────────────────┘
```
## 3. Data Flow
- `/users` CRUD, `/users/:id/reset-password`, `/users/:id/mfa`.
- WebSocket `security:sessions` to show login alerts.

## 4. Components
- UsersTable, UserModal, InviteBanner, RoleBadge, MFAStatus.
- hooks: `useUsers`, `useInviteForm`.

## 5. Validation
- Unique username/email check.
- Role change restricted (Admin only).
- Force password reset toggle.

## 6. Tests/Telemetry
- Unit tests for table sorting, search.
- Playwright create/edit/delete user.
- Events: `user_created`, `user_role_changed`.

## 7. References
- `NAS_AI_SYSTEM.md`, `Blueprint_WebUI_Auth.md` (session context), `CVE_CHECKLIST.md`.