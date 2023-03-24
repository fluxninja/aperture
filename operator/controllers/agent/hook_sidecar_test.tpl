{
  "apiVersion": "fluxninja.com/v1alpha1",
  "kind": "Agent",
  "metadata": {
    "name": "agent-sample"
  },
  "spec": {
    "sidecar": {
      "enabled": true
    },
    "config": {
      "etcd": {
        "endpoints": [
          "test.com"
        ]
      },
      "prometheus": {
        "address": "test.com"
      }
    }
  }
}
