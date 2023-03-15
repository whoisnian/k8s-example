const SERVER_CONFIG_MAP = {
  prod: {
    namespace: 'k8s-example',
    componentList: ['@nfs:prod']
  }
}

module.exports = {
  SERVER_CONFIG_MAP
}
