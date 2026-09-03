# Book Golang API

This repository contains a Go REST API for managing books. It uses PostgreSQL for persistence and can be built as a Docker image.

## Prerequisites

- Docker with Docker Compose support
- VS Code with the Dev Containers extension
- Ansible, if using the Ansible deployment option

The repository uses the VS Code Dev Container configuration in [`.devcontainer`](.devcontainer/devcontainer.json). It starts an application container with Go 1.26 and a PostgreSQL 16 container.

## Configuration

The application reads `DATABASE_URL` from `.env`. The included Dev Container uses host networking for the app and publishes PostgreSQL on port `5432`, so use `localhost` for the database host:

```env
DATABASE_URL=host=localhost port=5432 user=postgres password=postgres dbname=appdb sslmode=disable
```

## Local Development

### Open the development container

Open the repository in VS Code and run **Dev Containers: Reopen in Container**. The container configuration starts both services and runs these checks after creation:

```bash
go version
psql --version
```

Run the following commands from the VS Code terminal inside the `app` container:

### Test

```bash
make test
```

This runs `go test ./... -v`.

For coverage:

```bash
make test-cover
```

### Run the API locally

Ensure `.env` contains `DATABASE_URL`, then start the application:

```bash
go run .
```

The API listens on `http://localhost:8080`. In another terminal, verify it with:

```bash
curl http://localhost:8080/books
```

The application creates the `books` table automatically when it starts.

### Build locally

```bash
go build -o api-golang .
docker build -t books-go-api .
```

The Docker image uses the multi-stage build in [`Dockerfile`](Dockerfile) and produces a minimal runtime image.
## Deployment

Deployment is controlled by [`jenkinsfile`](jenkinsfile). Configure the Jenkins SCM job to discover only the deployment branches. Feature branches are therefore not built or deployed by this job.

The current branch mapping is:

| Source branch | Image tag prefix | GitOps branch | Ansible inventory |
|---------------|------------------|---------------|-------------------|
| `develop` | `develop` | `develop` | `dev.yml` |
| `stage` | `stage` | `stage` | `stage.yml` |

Each build produces an image tag such as `develop-42`.

### Option 1: GitOps manifest update

This is the active deployment integration in the Jenkinsfile. Jenkins runs tests, builds and pushes the Docker image, then:

1. Clones `MANIFEST_REPO_URL` using the `manifest-repo-ssh` credential.
2. Checks out the normalized source branch in `MANIFEST_REPO_BRANCH`.
3. Updates `MANIFEST_FILE` with the new image tag.
4. Commits and pushes the manifest change.

The manifest must contain either an `image:` field containing `APP_IMAGE` or a `newTag:` field. A GitOps controller must watch the manifest repository and apply the committed change; this Jenkins stage does not restart a host directly.

Required Jenkins configuration includes the `docker-hub-credentials` and `manifest-repo-ssh` credentials, plus valid values for `MANIFEST_REPO_URL` and `MANIFEST_FILE`.

`MANIFEST_REPO_BRANCH` is replaced by the normalized source branch in the `Build our image` stage. For example, `origin/feature/login` becomes `feature-login`. Such branches should not be discovered by the SCM job when feature branches must not be built or deployed.

### Option 2: Ansible deployment

The Ansible `Deploy` stage is currently commented out in [`jenkinsfile`](jenkinsfile). If enabled, it runs on an agent labeled `ansible-deploy` after the image is built and pushed, using the `ansible-vault-password` file credential:

```bash
ansible-playbook -i deploy/ansible/inventory/<environment>.yml \
  deploy/ansible/deploy.yml \
  --vault-password-file <vault-password-file> \
  --extra-vars "image_tag=<branch>-<build-number> dockerhub_user=<dockerhub-user>"
```

The playbook renders [`deploy/ansible/templates/docker-compose.yml`](deploy/ansible/templates/docker-compose.yml), then runs `podman-compose pull` and `podman-compose up -d`. It requires the encrypted variables in `deploy/ansible/group_vars/app_servers/vault.yml` and the selected inventory configuration. Development defaults deploy to `/opt/books-go-api` on port `8080`.

For a manual deployment:

```bash
docker build -t <dockerhub-user>/books-go-api:<tag> .
docker push <dockerhub-user>/books-go-api:<tag>
ansible-playbook -i deploy/ansible/inventory/dev.yml deploy/ansible/deploy.yml \
  --vault-password-file <vault-password-file> \
  --extra-vars "image_tag=<tag> dockerhub_user=<dockerhub-user>"
```

Choose one deployment owner for a build: the GitOps controller or Ansible. Enabling both deploys the same image through two independent mechanisms.

## Project Structure

```text
.
├── Dockerfile
├── jenkinsfile
├── main.go
├── Makefile
├── deploy/
│   └── ansible/
│       ├── deploy.yml
│       ├── inventory/
│       ├── group_vars/
│       └── templates/
├── internal/
├── api.http
└── README.md
```

## Useful Commands

```bash
make test
make test-cover
go build -o api-golang .
docker build -t books-go-api .
```
