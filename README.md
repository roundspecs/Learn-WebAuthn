# Project Name

> One line: what this project is.

**Status:** In progress / Done
**Started:** 2026-10-20
**Finished:** YYYY-MM-DD
**Stack:** e.g. Go, SQLite

---

## 1. What and why

WebAuthn is here to replace passwords and I think this is the future.

## 2. Core idea

<!-- Write this LAST, without looking at the code.
Explain the main idea in plain words, as if telling a friend.
If I can't write this, I haven't understood it yet. -->

## 3. How it's built

<!-- A small diagram or a file-by-file map. Example:

-->

```
main.go           -> backend
static/index.html -> frontend
static/main.js    -> frontend
```

## 4. Key decisions

<!-- Each decision: what I chose, why, and what I rejected. -->
- Database: Couchbase
  - Not in-memory maps because I want to test it across server restarts, and I also wanted to reuse this code in my couchbase projects
  - Not any RDBMS because the data being stored (e.g. webauthn.Credentials) are too complex

## 5. What broke

<!-- Bugs and wrong assumptions. For each: what happened, how I found the cause, the fix.
This is usually where the real learning is. -->

-

## 6. What I learned

<!-- 3-5 bullets. -->

-
-
-

## 7. Run it

<!-- Exact commands, copy-paste ready. Include setup. -->

```bash
go run .
```

## 8. Next

<!-- What I would extend or fix if I came back to this. -->

-

---

<!-- Terms I had to look up (keep a short glossary so future me doesn't have to look them up again) -->

## Glossary

| Term | Plain meaning |
|------|---------------|
|      |               |

## Links

<!-- Docs, papers, tutorials, or videos I used. -->

- [Web Dev Simplified - YouTube](https://www.youtube.com/watch?v=viZs1iVsLpA)
