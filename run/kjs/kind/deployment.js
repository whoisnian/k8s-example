// {
//   "spec": {
//     "template": {
//       "metadata": {
//         "labels": {
//           "k8s-example": "backend-file"
//         }
//       },
//       "spec": {
//         "containers": [
//           {
//             "env": [
//               {
//                 "name": "CFG_APIPREFIX",
//                 "value": "http://backend-api:8080"
//               }
//             ],
//             "image": "ghcr.io/whoisnian/k8s-example-backend-file:v0.0.5",
//             "name": "backend-file",
//             "ports": [
//               {
//                 "containerPort": 8081
//               }
//             ],
//             "volumeMounts": [
//               {
//                 "mountPath": "/app/uploads",
//                 "name": "uploads"
//               }
//             ]
//           }
//         ],
//           "volumes": [
//             {
//               "name": "uploads",
//               "persistentVolumeClaim": {
//                 "claimName": "uploads-pvc"
//               }
//             }
//           ]
//       }
//     }
//   }
// }
const generateDeploy = ({
  namespace,
  name,
  replicas,
  spec,
  args
}) => ({
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: { namespace, name, labels: { [`${namespace}`]: name } },
  spec: {
    replicas,
    selector: { matchLabels: { [`${namespace}`]: name } },
    template: {
      metadata: { labels: { [`${namespace}`]: name } },
      spec
    },
    ...args
  }
})

module.exports = {
  generateDeploy
}
