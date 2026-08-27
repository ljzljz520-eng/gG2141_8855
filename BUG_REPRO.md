# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	clubmembers/cmd/member-service	[no test files]
?   	clubmembers/internal/audit	[no test files]
?   	clubmembers/internal/fixtures	[no test files]
?   	clubmembers/internal/validation	[no test files]
ok  	clubmembers/internal/catalog	0.008s
--- FAIL: TestMemberPhoneChangeIsIsolated (0.00s)
    snapshot_test.go:14: peer phone changed to "13900139001"
FAIL
FAIL	clubmembers/internal/members	0.012s
ok  	clubmembers/internal/reports	0.010s
ok  	clubmembers/internal/store	0.091s
ok  	clubmembers/internal/workflows	0.033s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/member-service): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/member-service): exit `0`
