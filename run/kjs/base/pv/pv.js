const { generatePVHostPath, generatePVNFS } = require('../../kind/pv')
const { generatePVCHostPath, generatePVCNFS } = require('../../kind/pvc')

const basePV = ({ namespace }) => {
  const localPublicPV = generatePVHostPath({
    basename: 'local-public',
    storage: '10Gi',
    path: '/mnt/public',
    type: 'Directory'
  })

  const localPublicPVC = generatePVCHostPath({
    namespace,
    basename: 'local-public',
    storage: '10Gi'
  })

  const uploadsPV = generatePVNFS({
    basename: 'uploads',
    storage: '10Gi',
    server: '192.168.49.1',
    path: '/'
  })

  const uploadsPVC = generatePVCNFS({
    namespace,
    basename: 'uploads',
    storage: '10Gi'
  })

  return [localPublicPV, localPublicPVC, uploadsPV, uploadsPVC]
}

module.exports = {
  basePV
}
