kinds = [
    'pods',
    'deployments',
    'services',
    'namespaces'
]

output_formats = [
    'yaml',
    'json'
]

deployment_json = {
    'kind': 'Deployment',
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "name": '',
        "namespace": 'default',
        "labels": {
            "run": ''
        }
    },
    "spec": {
        "selector": {
            "matchLabels": {
                "run": ''
            }
        },
        "template": {
            "metadata": {
                "name": '',
                "labels": {
                    "run": ''
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": '',
                        "image": '',
                        "ports": [],
                        "env": []
                    }
                ]
            }
        }
    }
}
