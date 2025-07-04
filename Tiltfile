LOCAL_PORT = 3000

# Load extensions
load('ext://restart_process', 'docker_build_with_restart')
load('ext://uibutton', 'cmd_button', 'text_input')

pod_exec_script = '''
set -eu
# get k8s pod name from tilt resource name
POD_NAME="$(tilt get kubernetesdiscovery "$resource" -ojsonpath='{.status.pods[0].name}')"
kubectl exec "$POD_NAME" -- $command
'''

# API Backend
docker_build_with_restart(
    'gradewise-api-backend',
    './api-backend',
    build_args={
        'DEV': 'true',
    },
    live_update=[
        # Sync all source files to the container
        sync('./api-backend/src', '/app/src'),
    ],
    entrypoint='cargo run --'
    
)

# Frontend
docker_build(
    'gradewise-frontend',
    './frontend',
    live_update=[
        sync('./frontend', '/app')
    ],
    build_args={
        'DEV': 'true',
        'BASE_URL': 'http://localhost:%d' % LOCAL_PORT
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
k8s_yaml('k8s/api-backend-deployment.yaml')
k8s_yaml('k8s/frontend-deployment.yaml')

## Traefik
k8s_yaml('k8s/traefik/role.yml')
k8s_yaml('k8s/traefik/account.yml')
k8s_yaml('k8s/traefik/role-binding.yml')
k8s_yaml('k8s/traefik/traefik.yml')
k8s_yaml('k8s/traefik/traefik-services.yml')

# Ingress
k8s_yaml('k8s/traefik/ingress/api-backend.yml')
k8s_yaml('k8s/traefik/ingress/frontend.yml')


# Expose Traefik
k8s_resource('traefik-deployment', port_forwards=['%d:80' % LOCAL_PORT, '8080:8080'])