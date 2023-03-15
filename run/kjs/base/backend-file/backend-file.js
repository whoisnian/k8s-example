const { generateSVC } = require('../../kind/service')

const baseBackendFile = ({ namespace }) => {
  const backendFileSVC = generateSVC({
    namespace,
    name: 'backend-file',
    ports: [{ port: 8081, targetPort: 8081, protocol: 'TCP' }]
  })

  return [backendFileSVC]
}

module.exports = {
  baseBackendFile
}
