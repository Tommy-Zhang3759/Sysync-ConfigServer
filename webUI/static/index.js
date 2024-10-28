function loadClients() {
    fetch('/api/cliInfo')
        .then(response => response.json())
        .then(data => {
            console.log(data); // 打印返回的数据
            const clientList = document.getElementById('client-list');
            clientList.innerHTML = ''; // 清空当前列表

            if (data.clients) { // 检查 clients 是否存在
                data.clients.forEach(hostname => {
                    const div = document.createElement('div');
                    div.className = 'client-item';
                    div.innerText = hostname; // 直接使用 hostname
                    clientList.appendChild(div);
                });
            } else {
                console.error('返回的数据中没有 clients 字段');
            }
        })
        .catch(error => console.error('获取客户端列表时出错:', error));
}

// 请求选中客户端的详细信息
function requestClientInfo(clientId) {
    fetch(`/api/cliInfo?id=${clientId}`)
        .then(response => response.json())
        .then(client => {
            const clientDetails = document.getElementById('client-details');
            clientDetails.innerHTML = `
                <p>名称: ${client.name}</p>
                <p>IP: ${client.ip}</p>
                <p>状态: ${client.status}</p>
                <button onclick="performAction('${client.id}')">操作</button>
            `;
        })
        .catch(error => console.error('获取客户端信息时出错:', error));
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

// 加载服务器日志
function addLog(logMessage) {
    const logContent = document.getElementById('log-content');
    const p = document.createElement('p');
    p.innerText = logMessage;
    logContent.appendChild(p);
}

// 更新系统状态
function updateSystemStatus() {
    fetch('/api/system_status')
        .then(response => response.json())
        .then(status => {
            const statusElement = document.getElementById('system-status');
            statusElement.innerText = status;
        })
        .catch(error => console.error('更新系统状态时出错:', error));
}

