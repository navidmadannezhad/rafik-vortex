# Rafik Vortex

Reviews pull requests with an AI model, and copies Cursor rules and skills from a library repo into a project when you ask for them.

It listens on port `8080` and exposes two webhooks:

| Endpoint | When to call it |
| --- | --- |
| `POST /api/v1/webhook` | A pull request is opened or updated |
| `POST /api/v1/request-ai-instructions` | A push should install or refresh AI instructions |

The git API calls are Gitea-compatible (`/api/v1/repos/...`). Reviews are sent through [OpenRouter](https://openrouter.ai).

## Setup

Copy the sample env file and fill it in:

```bash
cp .env.sample .env
```

```env
# Base URL of your git host, for example https://felan.gitea.ir
GIT_URL=

# Token with read access to the library repo and write access to the repos you review or update
GIT_TOKEN=

# https://openrouter.ai/keys
OPENROUTER_API_KEY=

# Repo that stores rules and skills, organized by stack
LIBRARY_REPO_OWNER=
LIBRARY_REPO_NAME=
```

## Run

Locally, from the repo root:

```bash
go run main.go
```

With Docker, `start.sh` builds the image and starts it. `--network` controls which address port `8080` is published on: `local` (`127.0.0.1`), `full` (`0.0.0.0`), or `netbird` (the host's NetBird IPv4).

```bash
# development
./start.sh --profile dev --network local --action up

# production
./start.sh --profile prod --network full --action up

# stop
./start.sh --profile prod --action down
```

Point your git host's webhooks at `http://<host>:8080` plus the path below. The request body is the webhook payload as the host sends it.

## Review a pull request

Add a webhook on the target repo for pull request events, aimed at:

```text
POST /api/v1/webhook
```

On a pull request event the service:

1. Fetches the PR diff.
2. Asks the model for review comments.
3. Posts those comments back on the pull request.

A successful call returns `{"data":"Review completed"}`.

## Request AI instructions

This is how a project gets Cursor rules and skills from the library.

### 1. Layout of the library repo

`LIBRARY_REPO_OWNER` / `LIBRARY_REPO_NAME` must contain one folder per stack. Each stack has a `rules` directory and a `skills` directory. Every skill is its own folder with a `SKILL.md` file.

```text
react/
  rules/
    review.mdc
  skills/
    component-patterns/
      SKILL.md
vue/
django/
fastapi/
nest/
gin/
```

Supported stacks are `react`, `vue`, `django`, `fastapi`, `nest`, and `gin`. The folder name must match the stack exactly.

### 2. Add the webhook

On the **target** repo (the project that should receive the files), add a webhook for **push** events:

```text
POST /api/v1/request-ai-instructions
Content-Type: application/json
```

The body is the push payload from your git host. The handler reads `repository.owner.login`, `repository.name`, `ref`, and `commits`.

### 3. Trigger it with a commit message

The handler runs only when at least one commit in the push contains:

```text
--update-ai-instructs
```

Example:

```bash
git commit --allow-empty -m "chore: sync cursor files --update-ai-instructs"
git push
```

Pushes that do not include that flag are ignored.

### 4. What the service does

1. Lists the root of the target repo and uses the first of `package.json`, `requirements.txt`, or `go.mod` that it finds.
2. Reads that file and picks a stack:

   | File | Matched when the file contains | Stack |
   | --- | --- | --- |
   | `package.json` | `react` or `next` | `react` |
   | `package.json` | `vue` | `vue` |
   | `package.json` | `nest` | `nest` |
   | `requirements.txt` | `django` or `djangorestframework` | `django` |
   | `requirements.txt` | `fastapi` | `fastapi` |
   | `go.mod` | `go` | `gin` |

3. Loads `{stack}/rules/*` and `{stack}/skills/*/SKILL.md` from the library repo (default branch `main`).
4. Writes them into the target repo on the same branch that was pushed:

   ```text
   .cursor/rules/<rule file>
   .cursor/skills/<skill name>/SKILL.md
   ```

   Existing files are updated. Missing files are created. The commit message is `feat: add some ai instructions for you comrade!`.

A successful call returns `{"data":"Instructions requested"}`.

If the root has no stack file, or the detected stack has no folder in the library, the call fails and nothing is written.

## Contributing

Contributions are welcome.
