let selectedBucket = null;
let bucketsCache = null;

function el(id) { return document.getElementById(id); }

document.addEventListener('DOMContentLoaded', loadBuckets);

function show(elid, show = true) { el(elid).style.display = show ? '' : 'none'; }

function loadBuckets() {
  selectedBucket = null;
  show('bucketPanel', false);
  show('bucketsPanel', true);
  show('bucketError', false);
  show('bucketsLoading', true);
  fetch('/')
    .then(r => r.text())
    .then(xml => {
      show('bucketsLoading', false);
      renderBuckets(xml);
    })
    .catch(e => {
      show('bucketsLoading', false);
      el('bucketError').innerText = 'Ошибка загрузки списка бакетов';
      show('bucketError', true);
    });
}

function renderBuckets(xml) {
  const buckets = Array.from(xml.matchAll(/<Name>([^<]*)<\/Name>/g)).map(x => x[1]);
  bucketsCache = buckets;
  const wrap = el('buckets');
  wrap.innerHTML = '';
  if (!buckets.length) {
    show('emptyBuckets', true);
    return;
  }
  show('emptyBuckets', false);
  buckets.forEach(b => {
    const item = document.createElement('div');
    item.className = 'item';
    const left = document.createElement('span');
    left.textContent = b;
    item.append(left);
    const open = document.createElement('button');
    open.textContent = 'Открыть';
    open.onclick = () => showObjects(b);
    const del = document.createElement('button');
    del.textContent = 'Удалить';
    del.onclick = () => {
      fetch('/' + b, { method: 'DELETE' }).then(() => loadBuckets());
    };
    item.append(open); item.append(del);
    wrap.append(item);
  });
}

function createBucket() {
  const bname = el('bucketName').value;
  if (!bname) return;
  fetch('/' + encodeURIComponent(bname), { method: 'PUT' })
    .then(resp => {
      if (resp.status === 409)
        throw new Error('Бакет уже существует.');
      if (!resp.ok)
        throw new Error('Ошибка создания бакета.');
      el('bucketName').value = '';
      loadBuckets();
    }).catch(e => {
      el('bucketError').innerText = e.message;
      show('bucketError', true);
    });
}

function showObjects(bname) {
  selectedBucket = bname;
  show('bucketPanel', true);
  show('bucketsPanel', false);
  show('objError', false);
  el('bname').textContent = bname;
  updateObjects();
  el('uploadForm').onsubmit = uploadObject;
}

function updateObjects() {
  show('objectsLoading', true);
  fetchListObjects(selectedBucket);
}

function fetchListObjects(bucket) {
  fetch('/' + encodeURIComponent(bucket) + '/objects.csv')
    .then(r => r.ok ? r.text() : '')
    .then(text => {
      show('objectsLoading', false);
      renderObjects(text, bucket);
    })
    .catch(() => {
      show('objectsLoading', false);
      el('objError').innerText = 'Ошибка загрузки объектов';
      show('objError', true);
    });
}

function renderObjects(csv, bucket) {
  const wrap = el('objects');
  let lines = csv.split(/\r?\n/g).filter(Boolean);
  wrap.innerHTML = '';
  if (!lines.length) {
    show('emptyObjects', true);
    return;
  }
  show('emptyObjects', false);
  const table = document.createElement('table');
  table.className = 'table';
  const head = document.createElement('tr');
  ['Имя', 'Размер', 'Тип', 'Изменено', 'Действия'].forEach(t => {
    const th = document.createElement('th'); th.textContent = t; head.append(th);
  });
  table.append(head);
  lines.forEach(row => {
    const cols = row.split(',');
    if (!cols[0]) return;
    const tr = document.createElement('tr');
    [0, 1, 2, 3].forEach(i => {
      const td = document.createElement('td');
      if (i === 1) td.textContent = cols[1] + ' байт';
      else if (i === 3) td.textContent = formatTime(cols[3]);
      else td.textContent = cols[i] || '';
      tr.append(td);
    });
    const action = document.createElement('td');
    const download = document.createElement('button');
    download.textContent = '⤓';
    download.title = 'Скачать';
    download.onclick = () => window.open(`/${bucket}/${encodeURIComponent(cols[0])}`);
    const del = document.createElement('button');
    del.textContent = '✕';
    del.title = 'Удалить';
    del.onclick = () => {
      fetch(`/${bucket}/${encodeURIComponent(cols[0])}`, { method: 'DELETE' })
        .then(() => updateObjects())
    };
    action.append(download); action.append(del);
    tr.append(action);
    table.append(tr);
  });
  wrap.append(table);
}

function formatTime(t) {
  if (!t) return '';
  const d = new Date(t);
  return d.toLocaleString();
}

function uploadObject(e) {
  e.preventDefault();
  const fileInput = el('fileInput');
  if (!fileInput.files.length) return;
  const file = fileInput.files[0];
  show('objectsLoading', true);
  fetch(`/${selectedBucket}/${encodeURIComponent(file.name)}`, {
    method: 'PUT',
    headers: { 'Content-Type': file.type || 'application/octet-stream' },
    body: file
  }).then(resp => {
    show('objectsLoading', false);
    if (!resp.ok) {
      el('objError').innerText = 'Ошибка загрузки объекта.';
      show('objError', true);
      return;
    }
    fileInput.value = '';
    updateObjects();
  });
}

function backToBuckets() {
  show('bucketPanel', false);
  show('bucketsPanel', true);
  loadBuckets();
}
