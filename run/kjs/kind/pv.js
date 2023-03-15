const generatePV = ({
  basename,
  accessModes,
  capacity,
  persistentVolumeReclaimPolicy,
  storageClassName,
  args
}) => ({
  apiVersion: 'v1',
  kind: 'PersistentVolume',
  metadata: { name: `${basename}-pv` },
  spec: {
    accessModes,
    capacity,
    persistentVolumeReclaimPolicy,
    storageClassName,
    ...args
  }
})

const generatePVNFS = ({ basename, storage, server, path }) => {
  const accessModes = ['ReadWriteMany']
  const capacity = { storage }
  const persistentVolumeReclaimPolicy = 'Retain'
  const storageClassName = ''
  const args = { nfs: { server, path } }
  return generatePV({ basename, accessModes, capacity, persistentVolumeReclaimPolicy, storageClassName, args })
}

const generatePVHostPath = ({ basename, storage, path, type }) => {
  const accessModes = ['ReadWriteOnce']
  const capacity = { storage }
  const persistentVolumeReclaimPolicy = 'Delete'
  const storageClassName = ''
  const args = { hostPath: { path, type } }
  return generatePV({ basename, accessModes, capacity, persistentVolumeReclaimPolicy, storageClassName, args })
}

const generatePVLocal = ({ basename, storage, path, hostname }) => {
  const accessModes = ['ReadWriteOnce']
  const capacity = { storage }
  const persistentVolumeReclaimPolicy = 'Retain'
  const storageClassName = 'local-storage'
  const nodeAffinity = {
    required: {
      nodeSelectorTerms: [{
        matchExpressions: [{
          key: 'kubernetes.io/hostname',
          operator: 'In',
          values: [hostname]
        }]
      }]
    }
  }
  const args = { local: { path }, nodeAffinity }
  return generatePV({ basename, accessModes, capacity, persistentVolumeReclaimPolicy, storageClassName, args })
}

module.exports = {
  generatePV,
  generatePVNFS,
  generatePVHostPath,
  generatePVLocal
}
