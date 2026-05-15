# Casino Proxy V2

Go microservice for the **Cometa Gaming** platform.
Generated from the [cometagaming/ms-template-go-v2](https://github.com/cometagaming/ms-template-go-v2) Go backend template.

---

## Local Development

No external services required. PostgreSQL and Kafka run locally via Docker Compose.

**Prerequisites:** Go 1.25.4+, Docker and Docker Compose, Node.js 18+.

**First-time setup** (run once after cloning):

```bash
npm run setup
```

| Mode | Command | When to use |
|------|---------|-------------|
| **Hybrid** — infra in Docker, app runs locally | `make up` then `make run` | Active development, debugger |
| **Full Docker** — everything in containers | `make up-full` | CI-equivalent smoke test |

Full setup guide: [application/LOCAL_DEVELOPMENT.md](application/LOCAL_DEVELOPMENT.md)

---

## Testing

```bash
cd application
make test-all  # unit -> integration -> acceptance (all via Testcontainers, no extra setup)
```

---

## References

- **Template:** [cometagaming/ms-template-go-v2](https://github.com/cometagaming/ms-template-go-v2) — the Go microservice template this project was generated from
