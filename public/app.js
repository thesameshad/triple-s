let selectedBucket = null;

document.addEventListener('DOMContentLoaded', loadBuckets);

function loadBuckets() {
  selectedBucket = null;
  document.getElementById('bucketPanel').style.display = 'none';
  fetch('/')
    .then(r => r.text())
    .then(renderBuckets);
}

function renderBuckets(xml) {
  const buckets = Array.from(xml.matchAll(/<Name>([^<]*)<\/Name>/g)).map(x=>x[1]);
  const wrap = document.getElementById('buckets');
  wrap.innerHTML = '';
  buckets.forEach(b => {
    const el = document.createElement('div');
    el.className = 'item';
    el.textContent = b;
    const open = document.createElement('button');
    open.textContent = 'Открыть';
    open.onclick = () => showObjects(b);
    const del = document.createElement('button');
    del.textContent = 'Удалить';
    del.onclick = () => {
      fetch('/'+b, {method:'DELETE'})
        .then(()=>loadBuckets());
    };
    el.append(open); el.append(del);
    wrap.append(el);
  })
}

function createBucket() {
  const bname = document.getElementById('bucketName').value;
  fetch('/' + encodeURIComponent(bname), {method:'PUT'})
    .then(() => {
      document.getElementById('bucketName').value = '';
      loadBuckets();
    });
}

function showObjects(bname) {
  selectedBucket = bname;
  document.getElementById('bucketPanel').style.display = '';
  document.getElementById('bname').textContent = bname;
  updateObjects();
  document.getElementById('uploadForm').onsubmit = uploadObject;
}

function updateObjects() {
  fetch('/')
    .then(r => r.text())
    .then(xml => {
      if (!selectedBucket) return;
      fetchListObjects(selectedBucket);
    });
}

function fetchListObjects(bucket) {
  fetch('/'+bucket+'/objects.csv')
    .then(r => r.ok ? r.text() : '')
    .then(text => {
      renderObjects(text, bucket);
    });
}

function renderObjects(csv, bucket) {
  const wrap = document.getElementById('objects');
  const lines = csv.split(/\r?\n/g).filter(Boolean);
  wrap.innerHTML = '';
  lines.forEach(row => {
    const cols = row.split(',');
    if (!cols[0]) return;
    const el = document.createElement('div');
    el.className = 'item';
    el.textContent = cols[0] + ` (${cols[1]} байт)`;
    const download = document.createElement('button');
    download.textContent = 'Скачать';
    download.onclick = () => {
      window.open(`/${bucket}/${encodeURIComponent(cols[0])}`);
    };
    const del = document.createElement('button');
    del.textContent = 'Удалить';
    del.onclick = () => {
      fetch(`/${bucket}/${encodeURIComponent(cols[0])}`, {method:'DELETE'})
        .then(()=>updateObjects());
    };
    el.append(download); el.append(del);
    wrap.append(el);
  });
}

function uploadObject(e) {
  e.preventDefault();
  const fileInput = document.getElementById('fileInput');
  if (!fileInput.files.length) return;
  const file = fileInput.files[0];
  fetch(`/${selectedBucket}/${encodeURIComponent(file.name)}`, {
    method:'PUT',
    headers: { 'Content-Type': file.type || 'application/octet-stream' },
    body: file
  }).then(()=>{
    fileInput.value = '';
    updateObjects();
  });
}

function backToBuckets() {
  document.getElementById('bucketPanel').style.display = 'none';
  loadBuckets();
}
