// https://github.com/kubernetes-sigs/kustomize/blob/ce3e394a414387ce09679523902e981414b09a1a/plugin/builtin/sortordertransformer/SortOrderTransformer.go#L216
const orderMap = {
  Namespace: 0,
  ResourceQuota: 1,
  StorageClass: 2,
  CustomResourceDefinition: 3,
  ServiceAccount: 4,
  PodSecurityPolicy: 5,
  Role: 6,
  ClusterRole: 7,
  RoleBinding: 8,
  ClusterRoleBinding: 9,
  ConfigMap: 10,
  Secret: 11,
  Endpoints: 12,
  Service: 13,
  LimitRange: 14,
  PriorityClass: 15,
  PersistentVolume: 16,
  PersistentVolumeClaim: 17,
  Deployment: 18,
  StatefulSet: 19,
  CronJob: 20,
  PodDisruptionBudget: 21,

  Unknown: 100
}

const sortByKind = (a, b) => {
  return (orderMap[a.kind] || orderMap.Unknown) - (orderMap[b.kind] || orderMap.Unknown)
}

module.exports = {
  sortByKind
}
