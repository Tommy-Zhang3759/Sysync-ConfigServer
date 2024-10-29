async function loadClients() {
    try {
        const response = await fetch('/api/cliInfo');
        const data = await response.json(); // 使用 await 解析 JSON

        console.log(data); // 打印返回的数据
        const clientList = document.getElementById('client-list');
        clientList.innerHTML = ''; // 清空当前列表

        if (data) {
            data.forEach(hostname => {
                const div = document.createElement('div');
                div.className = 'client-item';
                div.innerText = hostname; // 直接使用 hostname
                clientList.appendChild(div);
                div.onclick = () => clientDetailedMenu(hostname);
            });
        } else {
            console.error('返回的数据中没有 clients 字段');
        }
    } catch (error) {
        console.error('获取客户端列表时出错:', error);
    }
}

async function clientDetailedMenu(hostname) {
    const c = await requestClientInfo(hostname);
    const clientDetails = document.getElementById('client-details');

    const paddedStatusCode = String(c.status_code).padStart(3, '0');

    clientDetails.innerHTML = `
        <p>Host name: ${c.host_name}</p>
        <p>IP: ${c.ip_addr}</p>
        <p>状态: ${paddedStatusCode}</p>
        <button onclick="performAction('${c.id}')">操作</button>
    `;
}

async function requestClientInfo(clientId) {
    const resp = await fetch(`/api/cliInfo?id=${clientId}`);
    const client = (await resp.json())[0];

    console.log(client);
    return client;
}

// 模拟操作客户端的功能
function performAction(clientId) {
    alert(`正在对客户端ID ${clientId} 执行操作`);
    fetch('/api/perform_action', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id: clientId })
    })
        .then(response => response.json())
        .then(data => {
            console.log('操作成功:', data);
        })
        .catch(error => console.error('执行操作时出错:', error));
}

function addLog(logMessage) {
    const logContent = document.getElementById('log-content');
    const p = document.createElement('p');
    p.innerText = logMessage;
    logContent.appendChild(p);
}

function updateSystemStatus() {
    fetch('/api/system_status')
        .then(response => response.json())
        .then(status => {
            const statusElement = document.getElementById('system-status');
            statusElement.innerText = status;
        })
        .catch(error => console.error('更新系统状态时出错:', error));
}

loadClients()