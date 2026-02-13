# .env
set dotenv-load
# Only if name !== .env
# set dotenv-filename := ".env"

# Aliases
alias mc := migrate-create
alias mr := migrate-run
alias mf := migrate-force
alias mv := migrate-version

# Run Air - Live reload
[group("dev")]
air:
  air -c .air.toml

# Create migration name = string
[group("migration")]
migrate-create name:
  migrate create -ext sql -dir ./migrations -seq {{name}}
  @echo "{{GREEN}}Migration {{name}} created"

# Run migration pos = up | down | down 1
[group("migration")]
migrate-run pos="up":
  @tmpfile=$(mktemp); \
    migrate -path ./migrations -database $DATABASE_URL {{pos}} 2>&1 | tee "$tmpfile"; \
    exit_code=${PIPESTATUS[0]}; \
    if grep -qi "no change" "$tmpfile"; then \
      :; \
    elif [ $exit_code -eq 0 ]; then \
      echo "{{GREEN}}Migration {{pos}} completed"; \
    fi; \
    rm "$tmpfile"

[group("migration")]
migrate-force version:
  migrate -path ./migrations -database $DATABASE_URL force {{version}}

[group("migration")]
migrate-version:
  @version=$(migrate -path ./migrations -database $DATABASE_URL version 2>&1 | head -1) && echo "Current migration version: $version"
