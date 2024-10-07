// WebSocket连接
let socket = new WebSocket("ws://your-server-address");

// 打开WebSocket连接
socket.onopen = function () {
    console.log("已连接到服务器");
    // 请求客户端列表
    socket.send(JSON.stringify({ action: "get_clients" }));
};

// 接收WebSocket消息
socket.onmessage = function (event) {
    const data = JSON.parse(event.data);

    // 根据收到的消息类型执行不同的操作
    if (data.type === "clients") {
        loadClients(data.clients);
    } else if (data.type === "client_info") {
        showClientInfo(data.client);
    } else if (data.type === "log") {
        addLog(data.log);
    } else if (data.type === "system_status") {
        updateSystemStatus(data.status);
    }
};

// 初始化客户端列表
function loadClients(clients) {
    const clientList = document.getElementById('client-list');
    clientList.innerHTML = ''; // 清空当前列表
    clients.forEach(client => {
        const div = document.createElement('div');
        div.className = 'client-item';
        div.innerText = client.name;
        div.onclick = () => requestClientInfo(client.id);
        clientList.appendChild(div);
    });
}

// 请求选中客户端的详细信息
function requestClientInfo(clientId) {
    socket.send(JSON.stringify({ action: "get_client_info", id: clientId }));
}

// 显示选中客户端的信息
function showClientInfo(client) {
    const clientDetails = document.getElementById('client-details');
    clientDetails.innerHTML = `
        <p>名称: ${client.name}</p>
        <p>IP: ${client.ip}</p>
        <p>状态: ${client.status}</p>
        <button onclick="performAction('${client.id}')">操作</button>
    `;
}

// 模拟操作客户端的功能
function performAction(clientId) {
    alert(`正在对客户端ID ${clientId} 执行操作`);
    socket.send(JSON.stringify({ action: "perform_action", id: clientId }));
}

// 加载服务器日志
function addLog(logMessage) {
    const logContent = document.getElementById('log-content');
    const p = document.createElement('p');
    p.innerText = logMessage;
    logContent.appendChild(p);
}

// 更新系统状态
function updateSystemStatus(status) {
    const statusElement = document.getElementById('system-status');
    statusElement.innerText = status;
}

// 处理WebSocket关闭
socket.onclose = function () {
    console.log("连接已关闭");
};

// 处理WebSocket错误
socket.onerror = function (error) {
    console.log("WebSocket错误: ", error);
};
