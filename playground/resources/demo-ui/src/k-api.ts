import { CoreV1Api, KubeConfig } from '@kubernetes/client-node'

const kc = new KubeConfig()

kc.loadFromString(JSON.stringify(process.env.KUBECONFIG))

const k8sApi = kc.makeApiClient(CoreV1Api)

const NAMESPACE = 'default' // Replace with your desired namespace
const SERVICE_NAME = 'service1-demo-app' // Replace with your service name

export const kApi = async () => {
  try {
    const serviceResponse = await k8sApi.readNamespacedService(
      SERVICE_NAME,
      NAMESPACE
    )

    const serviceEndpoint = serviceResponse.body.spec.clusterIP

    console.log('serviceEndpoint', serviceEndpoint)
    return {
      ...serviceResponse,
      serviceEndpoint,
    }
  } catch (err) {
    console.log('err', err)
  }
}
