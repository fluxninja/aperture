{
  "apiVersion": "fluxninja.com/v1alpha1",
  "kind": "Agent",
  "metadata": {
    "name": "agent-sample"
  },
  "spec": {
    "image":{
      "digest": "sha:432234"
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
