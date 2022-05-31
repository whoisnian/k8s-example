const createElement = (tag, options = {}) => {
  const element = document.createElement(tag)
  Object.entries(options).forEach(([k, v]) => {
    element.setAttribute(k, v)
  })
  return element
}

const reloadPage = () => window.location.reload()

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

(async () => {
  const urlParams = new URLSearchParams(window.location.search)
  document.getElementById('errorSpan').textContent = urlParams.get('err')

  const infoTable = document.getElementById('infoTable')
  while (infoTable.firstChild) { infoTable.removeChild(infoTable.firstChild) }

  /** @type { { Cid: string, Name: string, Size: number, Time: number }[] } */
  const fileinfos = await (await fetch('/api/files')).json()
  fileinfos.forEach(({ Cid: cid, Name: name, Size: size, Time: time }) => {
    const tr = createElement('tr')
    // 名称
    const nameTd = createElement('td')
    const rawLink = createElement('a', { href: `/file/data?${new URLSearchParams({ cid, name })}` })
    rawLink.textContent = name
    nameTd.appendChild(rawLink)
    tr.appendChild(nameTd)
    // 大小
    const sizeTd = createElement('td')
    sizeTd.textContent = calcFromBytes(size)
    tr.appendChild(sizeTd)
    // 时间
    const timeTd = createElement('td')
    timeTd.textContent = calcRelativeTime(time * 1000)
    tr.appendChild(timeTd)
    // 下载
    const downTd = createElement('td')
    const downloadLink = createElement('a', { href: `/file/data?${new URLSearchParams({ cid, name, download: true })}`, download: name })
    downloadLink.textContent = "下载"
    downTd.appendChild(downloadLink)
    tr.appendChild(downTd)
    // 删除
    const delTd = createElement('td')
    const deleteLink = createElement('a')
    deleteLink.textContent = "删除"
    deleteLink.onclick = async () => {
      await fetch(`/api/file?${new URLSearchParams({ cid })}`, { method: 'DELETE' })
      reloadPage()
    }
    delTd.appendChild(deleteLink)
    tr.appendChild(delTd)
    infoTable.appendChild(tr)
  })
})()