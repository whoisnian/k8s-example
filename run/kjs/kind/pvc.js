const generatePVC = ({
  namespace,
  basename,
  accessModes,
  resources,
  storageClassName,
  args // maybe undefined
}) => ({
  apiVersion: 'v1',
  kind: 'PersistentVolumeClaim',
  metadata: { name: `${basename}-pvc`, namespace },
  spec: {
    volumeName: `${basename}-pv`,
    accessModes,
    resources,
    storageClassName,
    ...args
  }
})

const generatePVCNFS = ({ namespace, basename, storage }) => {
  const accessModes = ['ReadWriteMany']
  const resources = { requests: { storage } }
  const storageClassName = ''
  return generatePVC({ namespace, basename, accessModes, resources, storageClassName })
}

const generatePVCHostPath = ({ namespace, basename, storage }) => {
  const accessModes = ['ReadWriteOnce']
  const resources = { requests: { storage } }
  const storageClassName = ''
  return generatePVC({ namespace, basename, accessModes, resources, storageClassName })
}

const generatePVCLocal = ({ namespace, basename, storage }) => {
  const accessModes = ['ReadWriteOnce']
  const resources = { requests: { storage } }
  const storageClassName = 'local-storage'
  return generatePVC({ namespace, basename, accessModes, resources, storageClassName })
}

module.exports = {
  generatePVC,
  generatePVCNFS,
  generatePVCHostPath,
  generatePVCLocal
}
