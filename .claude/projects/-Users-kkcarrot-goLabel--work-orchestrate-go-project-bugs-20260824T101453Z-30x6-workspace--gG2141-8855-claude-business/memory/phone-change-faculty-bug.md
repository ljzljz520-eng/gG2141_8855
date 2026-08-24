---
name: phone-change-faculty-bug
description: Member phone change cross-contaminated same-faculty peers; root cause and fix location
metadata:
  type: project
---

In the `clubmembers` module, the bug "修改手机号后把同学院另一名成员的号码也替换了" was rooted in `internal/members/snapshot.go`'s `ReplacePhone`. It had an inner loop that copied the new phone onto every other member sharing the same `Faculty`, corrupting unrelated records (the workflow path `UpdatePhone→Change→ApplyChange→SaveMember` was already correct; `ReplacePhone` is the list-level entry point with the defect).

**Why:** The function was written to broadcast phone changes faculty-wide, which contradicts the requirement that only the targeted member's profile changes while peers stay untouched.

**How to apply:** Fix keeps `ReplacePhone` scoped to the single targeted member (bump its Revision, leave peers alone). Tests live in `internal/members/snapshot_test.go` (ReplacePhone unit) and `internal/workflows/repro_test.go` (end-to-end `UpdatePhone`). When touching member mutation helpers, grep for `.Phone =` and faculty-grouping loops to ensure no other path spreads field changes across members.
