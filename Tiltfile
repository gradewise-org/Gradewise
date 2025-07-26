LOCAL_PORT = 3000
DEV = 'true'

# Load extensions
load('ext://restart_process', 'docker_build_with_restart')
load('ext://uibutton', 'cmd_button', 'text_input')

# Prune settings
docker_prune_settings(
    disable=False,
    max_age_mins=10,
    num_builds=1,
    keep_recent=1
)

pod_exec_script = '''
set -eu
# get k8s pod name from tilt resource name
POD_NAME="$(tilt get kubernetesdiscovery "$resource" -ojsonpath='{.status.pods[0].name}')"
kubectl exec "$POD_NAME" -- $command
'''

# Set up a local resource to build the API server binary for Linux
local_resource(
    'api-server-compile',
    'cd ./backend && make all',
    deps='./backend',
    ignore=['./backend/bin', './backend/internal/api/gen.go']
)

# Build the Docker image, only including the built binary
load('ext://restart_process', 'docker_build_with_restart')

docker_build_with_restart(
    'gradewise-api-backend',
    './backend',
    dockerfile='./backend/deployments/Dockerfile.dev',
    build_args={
        'APP': 'server',
        'GIN_MODE': 'debug' if DEV == 'true' else 'release',
    },
    only=['bin/server'], # narrow context to only the relevant binary
    live_update=[
        sync('./backend/bin/server', '/app'),
    ],
    entrypoint='/app',
)

docker_build_with_restart(
    'gradewise-temporal-worker',
    './backend',
    dockerfile='./backend/deployments/Dockerfile.dev',
    build_args={
        'APP': 'worker',
        'GIN_MODE': 'debug' if DEV == 'true' else 'release',
    },
    only=['bin/worker'], # narrow context to only the relevant binary
    live_update=[
        sync('./backend/bin/worker', '/app'),
    ],
    entrypoint='/app',
)



# Frontend
docker_build(
    'gradewise-frontend',
    './frontend',
    live_update=[
        sync('./frontend', '/app')
    ],
    build_args={
        'BASE_URL': 'http://localhost:%d' % LOCAL_PORT,
        'DEV': DEV
    },
)

# -> Reinstall package dependencies
cmd_button(
    'gradewise-frontend:bun install',
    argv=['sh', '-c', pod_exec_script],
    resource='gradewise-frontend',
    icon_name='sync',
    text='Install dependencies for the frontend',
    inputs=[
        text_input(
            'resource',
            'Resource name',
            default='gradewise-frontend',
            placeholder='Enter resource name'
        ),
        text_input(
            'command',
            'Command',
            default='bun install',
            placeholder='Enter command to run in the pod'
        )
    ]
)

# Apply Kubernetes manifests
k8s_yaml('k8s/backend/postgres-deployment.yaml')
k8s_yaml('k8s/backend/temporal-server-deployment.yaml')
k8s_yaml('k8s/backend/api-backend-deployment.yaml')
k8s_yaml('k8s/backend/temporal-worker-deployment.yaml')
k8s_yaml('k8s/frontend/frontend-deployment.yaml')

# Traefik Ingress Controller
k8s_yaml('k8s/traefik/role.yml')
k8s_yaml('k8s/traefik/account.yml')
k8s_yaml('k8s/traefik/role-binding.yml')
k8s_yaml('k8s/traefik/traefik.yml')
k8s_yaml('k8s/traefik/traefik-services.yml')

# Ingress Routes
k8s_yaml('k8s/traefik/ingress/api-backend.yml')
k8s_yaml('k8s/traefik/ingress/frontend.yml')


# Resource Dependencies and Port Forwards
k8s_resource(
    'temporal-server',
    port_forwards=['8233:7233'],  # gRPC API only
    resource_deps=['postgres']
)

k8s_resource(
    'temporal-ui',
    port_forwards=['8081:8080'],  # Web UI
    resource_deps=['temporal-server']
)

k8s_resource(
    'gradewise-api-backend',
    resource_deps=['postgres']
)

k8s_resource(
    'gradewise-temporal-worker',
    resource_deps=['temporal-server']
)

# Main Application Gateway
k8s_resource(
    'traefik-deployment',
    port_forwards=['%d:80' % LOCAL_PORT, '8080:8080'],  # App traffic, Traefik dashboard
    resource_deps=['gradewise-frontend', 'gradewise-api-backend']
)
