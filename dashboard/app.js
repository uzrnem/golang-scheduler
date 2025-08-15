const API_BASE = "http://localhost:8080";

function login() {
    const serviceId = document.getElementById("serviceId").value;
    const token = document.getElementById("token").value;

    if (!serviceId || !token) {
        alert("Please enter Service ID and Token");
        return;
    }

    // Save credentials locally
    localStorage.setItem("service_id", serviceId);
    localStorage.setItem("token", token);

    // Load tasks
    loadTasks();

    document.getElementById("login-page").style.display = "none";
    document.getElementById("dashboard").style.display = "block";
}

function logout() {
    localStorage.clear();
    document.getElementById("login-page").style.display = "block";
    document.getElementById("dashboard").style.display = "none";
}

async function loadTasks() {
    const res = await fetch(`${API_BASE}/tasks`, {
        headers: getAuthHeaders()
    });
    const tasks = await res.json();
    const container = document.getElementById("task-list");
    container.innerHTML = "";
    tasks.forEach(task => {
        const div = document.createElement("div");
        div.innerHTML = `<b>${task.name}</b> (${task.method}) - ${task.url}`;
        div.onclick = () => loadTaskDetail(task.id);
        container.appendChild(div);
    });
}

async function loadTaskDetail(taskId) {
    const res = await fetch(`${API_BASE}/tasks/${taskId}`, {
        headers: getAuthHeaders()
    });
    const task = await res.json();
    document.getElementById("task-detail").innerHTML = `
        <h4>${task.name}</h4>
        <p>URL: ${task.url}</p>
        <p>Method: ${task.method}</p>
        <p>Next Run: ${task.scheduled_at}</p>
        <button onclick="loadExecutions('${task.id}')">View Executions</button>
    `;
}

async function loadExecutions(taskId, pageNumber = 0, count = 10) {
    const res = await fetch(`${API_BASE}/tasks/${taskId}/executions?pageNumber=${pageNumber}&count=${count}`, {
        headers: getAuthHeaders()
    });
    const data = await res.json();
    const container = document.getElementById("execution-list");
    container.innerHTML = `<h4>Executions (Total: ${data.total_count})</h4>`;
    data.records.forEach(exec => {
        const div = document.createElement("div");
        div.innerHTML = `[${exec.status}] ${exec.status_code} - ${exec.response}`;
        container.appendChild(div);
    });
}

function showCreateTask() {
    document.getElementById("task-detail").innerHTML = `
        <h4>Create Task</h4>
        <input type="text" id="name" placeholder="Task Name">
        <input type="text" id="url" placeholder="URL">
        <input type="text" id="method" placeholder="Method (GET, POST...)">
        <textarea id="header" placeholder='Headers (JSON)'></textarea>
        <textarea id="payload" placeholder="Payload"></textarea>
        <input type="number" id="frequency" placeholder="Frequency">
        <input type="text" id="unit" placeholder="hour/day">
        <button onclick="createTask()">Save Task</button>
    `;
}

async function createTask() {
    const body = {
        name: document.getElementById("name").value,
        url: document.getElementById("url").value,
        method: document.getElementById("method").value,
        header: JSON.parse(document.getElementById("header").value || "{}"),
        payload: document.getElementById("payload").value,
        frequency: parseInt(document.getElementById("frequency").value),
        unit: document.getElementById("unit").value
    };

    const res = await fetch(`${API_BASE}/tasks`, {
        method: "POST",
        headers: { ...getAuthHeaders(), "Content-Type": "application/json" },
        body: JSON.stringify(body)
    });

    if (res.ok) {
        alert("Task created successfully");
        loadTasks();
    } else {
        alert("Error creating task");
    }
}

function getAuthHeaders() {
    return {
        "service_id": localStorage.getItem("service_id"),
        "token": localStorage.getItem("token")
    };
}

// Auto login if creds in storage
if (localStorage.getItem("service_id") && localStorage.getItem("token")) {
    document.getElementById("login-page").style.display = "none";
    document.getElementById("dashboard").style.display = "block";
    loadTasks();
}
