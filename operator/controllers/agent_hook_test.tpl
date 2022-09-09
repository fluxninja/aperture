{
  "apiVersion": "fluxninja.com/v1alpha1",
  "kind": "Agent",
  "metadata": {
    "name": "agent-sample"
  },
  "spec": {
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
