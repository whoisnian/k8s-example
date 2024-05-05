/** @returns { HTMLElement } */
const createElement = (tag, options = {}) => {
  const element = document.createElement(tag)
  Object.entries(options).forEach(([k, v]) => {
    element.setAttribute(k, v)
  })
  return element
}

const downloadFile = (url, filename) => {
  const link = createElement('a', {
    href: url,
    download: filename
  })
  link.click()
  link.remove()
}

const calcFromBytes = (raw) => {
  if (typeof raw === 'string') {
    raw = parseInt(raw)
  }
  if (raw >= 1125899906842624) {
    return (raw / 1125899906842624).toFixed(1) + ' P'
  } else if (raw >= 1099511627776) {
    return (raw / 1099511627776).toFixed(1) + ' T'
  } else if (raw > 1073741824) {
    return (raw / 1073741824).toFixed(1) + ' G'
  } else if (raw > 1048576) {
    return (raw / 1048576).toFixed(1) + ' M'
  } else if (raw > 1024) {
    return (raw / 1024).toFixed(1) + ' K'
  } else {
    return raw.toFixed(1) + ' B'
  }
}

const calcRelativeTime = (raw) => {
  if (typeof raw === 'object') {
    raw = raw.getTime()
  }
  const now = Date.now()
  const rtf = new Intl.RelativeTimeFormat('en-US', { style: 'long' })
  if (now - raw < 60000) {
    return rtf.format(Math.floor((raw - now) / 1000), 'second')
  } else if (now - raw < 3600000) {
    return rtf.format(Math.floor((raw - now) / 60000), 'minute')
  } else if (now - raw < 86400000) {
    return rtf.format(Math.floor((raw - now) / 3600000), 'hour')
  } else if (now - raw < 604800000) {
    return rtf.format(Math.floor((raw - now) / 86400000), 'day')
  } else if (now - raw < 2592000000) {
    return rtf.format(Math.floor((raw - now) / 604800000), 'week')
  } else if (now - raw < 31104000000) {
    return rtf.format(Math.floor((raw - now) / 2592000000), 'month')
  } else {
    return rtf.format(Math.floor((raw - now) / 31536000000), 'year')
  }
}

const reloadPage = () => window.location.reload()
const openUrl = (url) => { window.location.href = url }
const openUrlInNewTab = (url) => window.open(url, '_blank')

export {
  createElement,
  downloadFile,
  calcFromBytes,
  calcRelativeTime,
  reloadPage,
  openUrl,
  openUrlInNewTab
}
