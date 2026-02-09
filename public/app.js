function el(id) { return document.getElementById(id); }

document.addEventListener('DOMContentLoaded', loadBuckets);

function show(elid, b = true) { el(elid).style.display = b ? '' : 'none'; }

function loadBuckets() {
  show('bucketPanel', false);
  show('bucketsPanel', true);
  show('bucketError', false);
  fetch('/')
    .then(r => r.text())
    .then(xml => renderBuckets(xml))
    .catch(() => {
      el('bucketError').innerText = 'Ошибка загрузки бакетов.';
      show('bucketError', true);
    });
}

function renderBuckets(xml) {
  // Fixed regex: removed extra space before </Name>
  const buckets = Array.from(xml.matchAll(/<Name>([^<]+)<\/Name>/g)).map(x => x[1]);
  const wrap = el('buckets');
  wrap.innerHTML = '';
  if (!buckets.length) {
    show('emptyBuckets', true);
    return;
  }
  show('emptyBuckets', false);
  for (const b of buckets) {
    const d = document.createElement('div');
    d.className = 'item'; d.textContent = b + ' ';
    const open = document.createElement('button');
    open.textContent = 'Открыть'; open.onclick = () => showObjects(b);
    const del = document.createElement('button');
    del.textContent = 'Удалить'; del.onclick = () => fetch('/' + b, { method: 'DELETE' }).then(loadBuckets);
    d.append(open); d.append(del); wrap.append(d);
  }
}

function createBucket() {
  const b = el('bucketName').value.trim();
  if (!b) return;
  fetch('/' + encodeURIComponent(b), { method: 'PUT' })
    .then(r => { if (!r.ok) throw new Error(); loadBuckets(); })
    .catch(() => { el('bucketError').innerText = 'Ошибка создания.'; show('bucketError', true); });
  el('bucketName').value = '';
}

function showObjects(bucket) {
  el('bname').textContent = bucket;
  show('bucketPanel', true); show('bucketsPanel', false); show('objError', false);
  updateObjects(bucket);
  el('uploadForm').onsubmit = e => uploadObject(e, bucket);
}

function updateObjects(bucket) {
  fetch('/' + encodeURIComponent(bucket))
    .then(r => r.ok ? r.text() : '')
    .then(xml => renderObjectsXml(xml, bucket))
    .catch(() => {
      el('objError').innerText = 'Ошибка загрузки объектов.';
      show('objError', true);
    });
}

function renderObjectsXml(xml, bucket) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(xml, 'text/xml');
  const contents = Array.from(doc.getElementsByTagName('Contents'));
  const wrap = el('objects');
  wrap.innerHTML = '';
  if (!contents.length) { show('emptyObjects', true); return; }
  show('emptyObjects', false);
  const table = document.createElement('table');
  table.className = 'table';
  table.innerHTML = `<tr><th>Имя</th><th>Размер</th><th>Тип</th><th>Изменено</th><th></th></tr>`;
  contents.forEach(c => {
    const key = c.getElementsByTagName('Key')[0]?.textContent || '';
    const size = c.getElementsByTagName('Size')[0]?.textContent || '';
    const ctype = c.getElementsByTagName('ContentType')[0]?.textContent || '';
    const date = c.getElementsByTagName('LastModified')[0]?.textContent || '';
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${key}</td><td>${size} байт</td><td>${ctype}</td><td>${date.replace('T', ' ').replace(/\+.*/, '')}</td>`;
    const td = document.createElement('td');
    const download = document.createElement('button');
    download.textContent = '⤓'; download.title = 'Скачать';
    download.onclick = () => window.open(`/${bucket}/${encodeURIComponent(key)}`);
    const del = document.createElement('button');
    del.textContent = '✕'; del.title = 'Удалить';
    del.onclick = () => fetch(`/${bucket}/${encodeURIComponent(key)}`, { method: 'DELETE' }).then(() => updateObjects(bucket));
    td.append(download); td.append(del); tr.append(td);
    table.append(tr);
  });
  wrap.append(table);
}

function uploadObject(e, bucket) {
  e.preventDefault();
  const fin = el('fileInput');
  if (!fin.files.length) return;
  const file = fin.files[0];
  fetch(`/${bucket}/${encodeURIComponent(file.name)}`, {
    method: 'PUT',
    headers: { 'Content-Type': file.type || 'application/octet-stream' },
    body: file
  }).then(r => { if (r.ok) updateObjects(bucket); else alert('Ошибка загрузки файла'); });
  fin.value = '';
}

function backToBuckets() {
  show('bucketPanel', false);
  show('bucketsPanel', true);
  loadBuckets();
}