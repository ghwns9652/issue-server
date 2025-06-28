# Issue Tracker API

Simple MVC‑structured REST API server built with [Gin](https://github.com/gin-gonic/gin). It manages **Users** and **Issues** following the business rules outlined in the specification.

---

## Prerequisites
- **Go 1.24** or newer
- **Newman** (CLI runner for Postman collections)
  ```bash
  npm install -g newman
  ```

---

## Run locally
```bash
git clone https://github.com/ghwns9652/issue-server.git
cd issue-server
go mod tidy
go run main.go
```

The server starts at **<http://localhost:8080>**.

On start‑up three default users are pre‑loaded:

```json
[
  { "id": 1, "name": "김개발" },
  { "id": 2, "name": "이디자인" },
  { "id": 3, "name": "박기획" }
]
```

---

## API Endpoints

| Method | Path        | Description                   |
| ------ | ----------- | ----------------------------- |
| POST   | /issue      | Create a new issue            |
| GET    | /issues     | List issues (optional status) |
| GET    | /issue/:id  | Get issue by ID               |
| PATCH  | /issue/:id  | Update issue                  |

All error responses follow:

```json
{ "error": "에러 메시지", "code": 400 }
```

---

## Example Newman Commands
테스트 컬렉션 전체를 돌려 보고 싶다면:
```bash
newman run postman_collection.json -r cli,htmlextra --reporter-htmlextra-export newman-report.html
```
생성된 `newman-report.html`에서 상세 리포트를 확인할 수 있습니다.

---

## Notes

- **Storage** is in‑memory for simplicity. Restarting the server resets data. Swap with a DB (e.g. SQLite + GORM) by replacing the repository layer.
- All timestamps are UTC (ISO‑8601).
- Business rules (status transitions, user validation, immutability after completion/cancellation) are enforced in `controllers/issue_controller.go`.
