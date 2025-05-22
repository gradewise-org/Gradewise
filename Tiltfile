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
    trigger=['./api-backend'],
    entrypoint=['/app/api'],
    live_update=[
        # Sync all Go files to the container
        sync('./api-backend', '/app'),
        
        # Rebuild the binary when code changes
        run('cd /app && go build -o ./api ./main.go')
    ]
)

# Frontend
docker_build(
    'gradewise-frontend',
    './frontend',
    live_update=[
        sync('./frontend', '/app')
    ],
    build_args={'DEV': 'true'}
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

# Configure port forwarding
k8s_resource('gradewise-api-backend', port_forwards='8080:8080') 
k8s_resource('gradewise-frontend', port_forwards='3000:3000')