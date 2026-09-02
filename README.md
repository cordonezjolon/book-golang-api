# Book Golang API

This repository contains a Go REST API for managing books. It uses PostgreSQL for persistence and can be built as a Docker image and deployed with Ansible through the Jenkins pipeline included in the project.

## Prerequisites

- Go 1.26+
- Docker / Podman
- PostgreSQL
- Ansible (for deployment)
- Docker Hub credentials for image publishing

## Configuration

This project requires a PostgreSQL database. Set the following environment variable in your `.env` file:

```env
DATABASE_URL=host=postgres port=5432 user= password= dbname=appdb sslmode=disable
```

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `host=localhost port=5432 user=postgres password= dbname=mydb sslmode=disable` |

## Build

### Run tests

```bash
make test
```

This executes:

```bash
go test ./... -v
```

### Build the Go binary locally

```bash
go build -o api-golang .
```

### Build the Docker image

The project includes a multi-stage Docker build in `Dockerfile`.

```bash
docker build -t books-go-api .
```

This produces a minimal runtime image based on `scratch` and runs the compiled binary from `/api-golang`.

## Deploy

### Deployment flow

The deployment process is defined in the Jenkins pipeline (`jenkinsfile`) and uses:

1. Go tests
2. Docker image build
3. Docker Hub push
4. Ansible deployment to the target host

The pipeline automatically selects the environment based on the Git branch:

- `develop` -> `dev`
- `staging` -> `stage`

### Ansible deploy configuration

The deploy playbook is:

```bash
ansible-playbook -i deploy/ansible/inventory/dev.yml deploy/ansible/deploy.yml \
  --extra-vars "image_tag=<tag> dockerhub_user=<dockerhub-user>"
```

The deployment renders the Docker Compose file from `deploy/ansible/templates/docker-compose.yml` and then runs:

```bash
podman-compose pull
podman-compose up -d
```

This deploys the application on port `8080` as configured in `deploy/ansible/group_vars/dev/vars.yml`.

### Manual deployment example

```bash
docker build -t <dockerhub-user>/books-go-api:<tag> .
docker push <dockerhub-user>/books-go-api:<tag>
ansible-playbook -i deploy/ansible/inventory/dev.yml deploy/ansible/deploy.yml \
  --extra-vars "image_tag=<tag> dockerhub_user=<dockerhub-user>"
```

### Runtime notes

- The app exposes port `8080`.
- The deployed container is recreated using the new image on every deploy.
- The `wait_for` step in the playbook waits until the app is listening before finishing.

## Project Structure

```text
.
├── Dockerfile
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

# add info deploy - prune removed
