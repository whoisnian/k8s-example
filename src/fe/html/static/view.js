import { createElement, downloadFile, dialogCloser, calcFromBytes, calcRelativeTime, reloadPage, openUrl } from './function.js'
import { fetchGetJSON, fetchDeleteHead, fetchGetJSONWithStatus, fetchPostJSONWithStatus } from './request.js'

/** @param {{ id: number, name: string, size: number, created_at: string }} */
const createFileItem = ({ id, name, size, created_at }) => {
  const tr = createElement('tr')

  const nameTd = createElement('td', { style: 'text-align:left;' })
  nameTd.textContent = name

  const sizeTd = createElement('td', { style: 'text-align:right;' })
  sizeTd.textContent = calcFromBytes(size)

  const ctime = new Date(created_at)
  const timeTd = createElement('td', { style: 'text-align:right;', title: ctime.toLocaleString() })
  timeTd.textContent = calcRelativeTime(ctime)

  const actionTd = createElement('td', { style: 'text-align:right;' })
  const downBtn = createElement('button')
  downBtn.textContent = 'SAVE'
  downBtn.onclick = () => downloadFile(`/file/objects/${id}`, name)
  const delBtn = createElement('button')
  delBtn.textContent = 'DEL'
  delBtn.onclick = async () => {
    await fetchDeleteHead(`/file/objects/${id}`)
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

const fileListInput = document.getElementById('fileListInput')
const fileListButton = document.getElementById('fileListButton')
const fileListText = document.getElementById('fileListText')
fileListButton.onclick = () => fileListInput.click()
fileListText.value = 'No file selected'
fileListInput.onchange = () => {
  if (fileListInput.files.length == 1) fileListText.value = fileListInput.files.item(0).name
  else fileListText.value = `${fileListInput.files.length} files selected`
}

const uploadFormButton = document.getElementById('uploadFormButton')
uploadFormButton.onclick = (e) => {
  e.preventDefault()

  // https://html.spec.whatwg.org/multipage/form-control-infrastructure.html#multipart-form-data
	// The order of parts must be the same as the order of fields in entry list.
  const formData = new FormData()
  const sizes = []
  for (const f of fileListInput.files) sizes.push(f.size)
  formData.append('fileSize', new Blob([JSON.stringify(sizes)], { type: 'application/json' }), 'sizes.json')
  for (const f of fileListInput.files) formData.append('fileList', f, f.name)

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

const usernameSpan = document.getElementById('usernameSpan')
const logoutButton = document.getElementById('logoutButton')
logoutButton.onclick = () => openUrl('/user/logout')

const signinButton = document.getElementById('signinButton')
const signinDialog = document.getElementById('signinDialog')
signinButton.onclick = () => signinDialog.showModal()
signinDialog.onclick = dialogCloser(signinDialog)

const signinFormButton = document.getElementById('signinFormButton')
signinFormButton.onclick = async (e) => {
  e.preventDefault()
  const formData = new FormData(document.getElementById('signinForm'))
  const { ok } = await fetchPostJSONWithStatus('/user/signin', Object.fromEntries(formData.entries()))
  if (ok) reloadPage()
}

const signupButton = document.getElementById('signupButton')
const signupDialog = document.getElementById('signupDialog')
signupButton.onclick = () => signupDialog.showModal()
signupDialog.onclick = dialogCloser(signupDialog)

const signupFormButton = document.getElementById('signupFormButton')
signupFormButton.onclick = async (e) => {
  e.preventDefault()
  const formData = new FormData(document.getElementById('signupForm'))
  const { ok } = await fetchPostJSONWithStatus('/user/signup', Object.fromEntries(formData.entries()))
  if (ok) reloadPage()
}

const { ok, content: userInfo } = await fetchGetJSONWithStatus('/user/info')
if (ok) {
  usernameSpan.textContent = userInfo.name
  usernameSpan.style = logoutButton.style = 'display:block;'
  signinButton.style = signupButton.style = 'display:none;'
}

const infoTable = document.getElementById('infoTable')
const fileinfos = await fetchGetJSON('/file/objects')
fileinfos.forEach(info => infoTable.appendChild(createFileItem(info)))
