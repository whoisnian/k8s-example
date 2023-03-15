const generateSVC = ({
  namespace,
  name,
  ports,
  args
}) => ({
  apiVersion: 'v1',
  kind: 'Service',
  metadata: { namespace, name, labels: { [`${namespace}`]: name } },
  spec: {
    ports,
    selector: { [`${namespace}`]: name },
    ...args
  }
})

module.exports = {
  generateSVC
}
