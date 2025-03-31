import websocketService from "./services/websocketService";

// Define action types as constants
const ActionTypes = {
    LOGIN_CONFIRMATION: 'LOGIN_CONFIRMATION',
    WELCOME: 'WELCOME',
    ERROR: 'ERROR',
};

// Main function to handle login
async function HandleLoginRequest() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    try {
        const data = await loginUser(username, password);
        localStorage.setItem("jwtToken", data.accessToken);
        
        // Attempt to set up WebSocket connection after successful login
        const socket = await setupWebSocket(data.accessToken, data.id);
        console.log("WebSocket connection established:", socket);
        
        return data; // Return data for further processing if needed
    } catch (error) {
        console.error("Login Error:", error.message);
        throw error; // Re-throw for handling by caller
    }
}


// Function to perform the login request
async function loginUser(username, password) {
    const API_URL = import.meta.env.VITE_API_BASE_URL;
    const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    });

    console.log("Raw Response:", response);
    console.log("Response Status:", response.status);

    if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data = await response.json();
    console.log("Parsed Response:", data);
    return data; // Return the parsed data
}

// Function to set up the WebSocket connection and message handling
async function setupWebSocket(accessToken, userId) {
    
    const socket = websocketService.connect(accessToken).then(() => {
        console.log("Websocket connection established. You can now send messages.");
    }).catch(error => {
        console.error("Failed to connect to WebSocket", error);
    });

    websocketService.sendMessage({
        type: ActionTypes.LOGIN_CONFIRMATION,
        sender_id: userId || 'anonymous',
    });

    websocketService.onMessage((message) => {
        console.log("Login handler received message:", message);
        handleWebSocketMessage(message);
    });

    return socket; // Return the socket if needed
}

// Function to handle sending messages
function sendMessage() {
    const messageInput = document.getElementById("message-input").value; // Get message input
    console.log("Message:", messageInput)
    websocketService.sendMessage(messageInput); // Call the sendMessage method from WebSocketService
    document.getElementById("message-input").value = ""; // Clear input
}


// Function to handle incoming WebSocket messages
function handleWebSocketMessage(message) {
    if (message.type === ActionTypes.WELCOME) {
        displayWelcomeMessage(message.content);
    } else if (message.type === ActionTypes.ERROR) {
        handleError(message.error);
    }
}

// Function to display a welcome message
function displayWelcomeMessage(content) {
    console.log("Welcome message:", content);
    const messageElement = document.getElementById('welcome-message');
    if (messageElement) {
        messageElement.textContent = content;
        messageElement.classList.remove('hidden');
    }
}

// Function to handle errors
function handleError(error) {
    console.error("Server error:", error);
    alert(`Server error: ${error}`);
}

// Event listener for login form submission
document.getElementById("login-form").addEventListener("submit", (event) => {
    event.preventDefault(); // Prevent default form submission
    HandleLoginRequest(); // Trigger the login function
});

// Add event listener for the send button
document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("sendButton").addEventListener("click", sendMessage); // Attach event listener
});


export default HandleLoginRequest;

