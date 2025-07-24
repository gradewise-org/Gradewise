# Gradewise

## Prerequisites

### 1. Install Tilt

**Recommended: Homebrew (macOS/Linux)**

```bash
brew install tilt-dev/tap/tilt
```

**Other platforms:** See [Tilt installation docs](https://docs.tilt.dev/install.html)

### 2. Install Kind (Kubernetes in Docker)

Kind is recommended for local development - it's fast (~20s startup), robust, and supports local image registries.

**Recommended: Homebrew + ctlptl**

```bash
brew install kind ctlptl
ctlptl create cluster kind --registry=ctlptl-registry
```

**Other platforms:** See [Tilt's cluster setup guide](https://docs.tilt.dev/choosing_clusters.html)

### 3. Install Traefik CRDs

Install [Traefik CRDs](https://doc.traefik.io/traefik/providers/kubernetes-crd/) for advanced routing:

```bash
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v3.5/docs/content/reference/dynamic-configuration/kubernetes-crd-definition-v1.yml
```

### 4. Key Commands

- `tilt up` - Start your development environment
- `tilt down` - Stop all services
- `kubectl get pods` - Check service status

## Architecture

- **Frontend**: SvelteKit application
- **Backend**: Go API with Temporal workflows
- **Database**: PostgreSQL with separate databases for API and Temporal
- **Infrastructure**: Kubernetes with Traefik ingress

## Quick Start

1. **Start the development environment**:

   ```bash
   tilt up
   ```

   ⏱️ **Note**: Initial startup may take ~60 seconds as health probes need time to succeed and services start in dependency order (PostgreSQL → Temporal → Worker).

2. **Access the application**:

   - Frontend: http://localhost:3000
   - API: http://localhost:3000/api/health
   - Temporal Web UI: http://localhost:8081

3. **Test the workflow**:
   ```bash
   curl "http://localhost:3000/api/greet?name=jeremiah"
   ```

## Development

### Database Management

- **Clean database**: `./backend/scripts/clean-dev-data.sh` (stops Tilt, cleans data)
- **Quick clean**: Use "clean-dev-data" button in Tilt UI (while Tilt running)

### Building & Debugging

- **Build backend**: The backend must be built for Linux containers. Tilt handles this automatically with `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make all`
- **View logs**: `kubectl logs -l app=gradewise-api-backend`
- **Check status**: `kubectl get pods`

### Development Workflow

1. Start: `tilt up`
2. Code → Tilt auto-rebuilds
3. Test endpoints at http://localhost:3000 (main app), http://localhost:8081 (Temporal UI)
4. Clean data when needed with script or Tilt button
5. Debug with `kubectl logs <pod-name>` or Traefik dashboard at http://localhost:8080
