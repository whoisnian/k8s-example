const { basePV } = require('../base/pv/pv')
const { sortByKind } = require('../kind/sort')
const { SERVER_CONFIG_MAP } = require('./config')

const runMain = async (CLUSTER_NAME) => {
  const SERVER_CONFIG = SERVER_CONFIG_MAP[CLUSTER_NAME]
  if (!SERVER_CONFIG) throw new Error(`invalid cluster name: ${CLUSTER_NAME}`)

  const result = [
    ...basePV(SERVER_CONFIG)
  ]

  result.sort(sortByKind)

  console.log(JSON.stringify(result, null, '  '))
}

const [
  , // node
  , // script.js
  CLUSTER_NAME = 'prod'
] = process.argv

runMain(CLUSTER_NAME)
