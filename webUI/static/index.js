function checkType(variable, name, type = 'string') {
    if (typeof variable !== type) {
        throw new TypeError(`${name} must be a ${type}`);
    }
}




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
                div.innerText = hostname['host_name']; // 直接使用 hostname
                clientList.appendChild(div);
                div.onclick = () => clientDetailedMenu(hostname['sysync_id']);
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
        <button onclick="syncSettings(null, '${c.ip_addr}', null, null, null)">Syncing System Configrations</button>
    `;
}

async function requestClientInfo(clientId) {
    const resp = await fetch(`/api/cliInfo?id=${clientId}`);
    const client = (await resp.json())[0];

    console.log(client);
    return client;
}

// 模拟操作客户端的功能
function syncSettings(destSysyncId, destIpAddr, destPort, server_ip_addr, server_port) {
    if (!Array.isArray(destSysyncId) || !destSysyncId.length) {
        destSysyncId = [destSysyncId];
    }
    if (!Array.isArray(destIpAddr) || !destIpAddr.length) {
        destIpAddr = [destIpAddr];
    }
    fetch('/api/func', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            dest_sysync_id: destSysyncId,
            f_name: "update_host_name",
            dest_ip: destIpAddr,
            dest_port: destPort,
            host_ip: server_ip_addr,
            host_port: server_port
        })
    })
        .then(response => response.json())
        .then(data => {
            console.log('Request sent:', data);
        })
        .catch(error => console.error('Error while sending:', error));
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
requestClientInfo()