import { createElement, downloadFile, calcFromBytes, calcRelativeTime, reloadPage } from './function.js'
import { fetchGetJSON, fetchDeleteHead } from './request.js'

/** @param {{ id: number, name: string, size: number, created_at: string }} */
const createFileItem = ({ id, name, size, created_at }) => {
  const tr = createElement('tr')

  const nameTd = createElement('td', { style: "text-align:left;" })
  nameTd.textContent = name

  const sizeTd = createElement('td', { style: "text-align:right;" })
  sizeTd.textContent = calcFromBytes(size)

  const ctime = new Date(created_at)
  const timeTd = createElement('td', { style: "text-align:right;", title: ctime.toLocaleString() })
  timeTd.textContent = calcRelativeTime(ctime)

  const actionTd = createElement('td', { style: "text-align:right;" })
  const downBtn = createElement('button')
  downBtn.textContent = "SAVE"
  downBtn.onclick = () => downloadFile(`/file/object/${id}`, name)
  const delBtn = createElement('button')
  delBtn.textContent = "DEL"
  delBtn.onclick = async () => {
    await fetchDeleteHead(`/file/object/${id}`)
    tr.remove()
  }
  actionTd.appendChild(downBtn)
  actionTd.append(' ')
  actionTd.appendChild(delBtn)

  tr.appendChild(nameTd)
  tr.appendChild(sizeTd)
  tr.appendChild(timeTd)
  tr.appendChild(actionTd)
  return tr
}

const uploadForm = document.getElementById('uploadForm')
const fileSizeInput = document.getElementById('fileSizeInput')
const fileListInput = document.getElementById('fileListInput')

const fileListButton = document.getElementById('fileListButton')
fileListButton.onclick = () => fileListInput.click()

const fileListText = document.getElementById('fileListText')
fileListText.value = 'No file selected'
fileListInput.onchange = () => {
  if (fileListInput.files.length == 1) fileListText.value = fileListInput.files.item(0).name
  else fileListText.value = `${fileListInput.files.length} files selected`
}

const uploadButton = document.getElementById('uploadButton')
uploadButton.onclick = () => {
  const sizes = []
  for (const f of fileListInput.files) sizes.push(f.size)
  fileSizeInput.value = JSON.stringify(sizes)

  const formData = new FormData(uploadForm)
  new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.withCredentials = true
    xhr.onload = () => {
      if (200 <= xhr.status && xhr.status < 300) resolve(xhr.response)
      else reject(xhr.status)
    }
    xhr.onerror = () => reject(xhr.status)
    xhr.upload.onprogress = (event) => {
      fileListText.value = event.total ? `${Math.round(100 * event.loaded / event.total)} %` : `Uploading`
    }
    xhr.open('POST', '/file/objects')
    xhr.send(formData)
  }).then(reloadPage).catch(console.error)
}

const infoTable = document.getElementById('infoTable')
const fileinfos = await fetchGetJSON('/file/objects')
fileinfos.forEach(info => infoTable.appendChild(createFileItem(info)))
