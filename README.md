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

Registration ceremony:
``` text
client                                              server
------                                              ------
POST /register/begin \
--body {username: "zarif"}------------------------->
                                                    if !username: error
                                                    if username exists: error
                                                    create user with a random 64-byte ID
                                                    {id, user} -> db
                                                    webauthn.BeginReg(user) -> options, session, err
                                                    // session has UserId inside it. We will use it later.
                                                    if err: status 500 (internal server error)
                                                    generate random token for session cookie
                                                    {"webauthn_ceremony", token} -> cookie
                                                    {token, session} -> db
<---------------------------------------------------options // it has the challenge inside it
navigator.credentials.create(options) -> data
// data has the signed answer to challenge
POST /register/finish \
--body options
--cookie "webauthn_ceremony=..."------------------->
                                                    cookie("webauthn_ceremony") -> token
                                                    db(token) -> session
                                                    deleteCookie("webauthn_ceremony") // we don't need it any more
                                                    db(session.UserID) -> user
                                                    webauthn.FinishReg(user, session, data) -> credentials, err
                                                    if err: status 401 (unauthorized)
                                                    user.credentials.append(credentials)
                                                    {id, user} -> db // upsert user. now it has new credentials registered
<---------------------------------------------------status 201 (created)
```

Login ceremony:
``` text
client                                              server
------                                              ------
POST /login/begin \
--body {username: "zarif"}------------------------->
                                                    if !username: error
                                                    if username !exists: status 404 (not found)
                                                    db(username) -> user
                                                    webauthn.BeginReg(user) -> assertion, session, err
                                                    // session has UserId inside it. We will use it later.
                                                    if err: status 500 (internal server error)
                                                    generate random token for session cookie
                                                    {"webauthn_ceremony", token} -> cookie
                                                    {token, session} -> db
<---------------------------------------------------assertion // it has the challenge inside it
navigator.credentials.get(assertion) -> data
// data has the signed answer to challenge
POST /register/finish \
--body options
--cookie "webauthn_ceremony=..."------------------->
                                                    cookie("webauthn_ceremony") -> token
                                                    db(token) -> session
                                                    deleteCookie("webauthn_ceremony") // we don't need it any more
                                                    db(session.UserID) -> user
                                                    webauthn.FinishReg(user, session, data) -> credentials, err
                                                    if err: status 401 (unauthorized)
                                                    user.credentials.append(credentials)
                                                    {id, user} -> db // upsert user. now it has new credentials registered
<---------------------------------------------------status 201 (created)
```
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
