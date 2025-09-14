# Changelog


## [1.1.0] - 2025-09-14

### 🚀 Features

- **Org**
  - Refactored service/usecase to return standard errors instead of `ApiError` for consistency ([c910b75])
  - Updated repository delete method to return `record not found` error instead of `nil` when no row is affected ([279964f])
  - Delete workspace and unassign org when Clerk update fails in usecase ([b770d87])
- **Routes**
  - Moved org-related routes into the org module for better cohesion ([61eb714])
- **Docs**
  - Updated org API documentation ([38fe739])
- **Tests**
  - Added org integration test ([72c9b25])
  - Removed outdated org service test ([e36cc99])