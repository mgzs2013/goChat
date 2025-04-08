import websocketService from "./services/websocketService";


// Define action types as constants
const ActionTypes = {
    LOGIN_CONFIRMATION: 'LOGIN_CONFIRMATION',
    WELCOME: 'WELCOME',
    ERROR: 'ERROR',
};

// Main function to handle login
async function HandleLoginRequest() {
    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value.trim();
    console.log("Attempting to log in with username:", username);
    console.log("Attempting to log in with username:", password);
    try {
        const data = await loginUser(username, password); // Destructure directly
        console.log("Login successful, received access token:", data);
        localStorage.setItem("jwtToken", data.accessToken);

        

        // Set up WebSocket connection after successful login
        await setupWebSocket(data.accessToken);
    } catch (error) {
        console.error("Login Error:", error.message);
        // Handle error feedback to the user if needed
        alert("Login failed: " + error.message);
    }
}



// Function to perform the login request
async function loginUser(username, password) {
    const API_URL = `http://localhost:8080/login`; // Ensure this is correct
    console.log("Using API_URL:", API_URL); // Log the API URL

    const response = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
    });

    console.log("Raw Response:", response);
    console.log("Response Status:", response.status);

    // Check if the response is OK
    if (!response.ok) {
        // Attempt to read the response body as text to log the error
        const errorText = await response.text(); // Read the response as text
        throw new Error(`HTTP error! Status: ${response.status}, Details: ${errorText}`);
    }

    const data = await response.json(); // Parse the response JSON
    console.log("Parsed Response:", data); // Log the parsed response
    return data; // Return the parsed data
}



// Function to set up the WebSocket connection and message handling
async function setupWebSocket(accessToken) {
    
    const wsUrl = `ws://localhost:8080/ws?token=${accessToken}`;
    console.log("Connecting to WebSocket at:", wsUrl);

    try {
        await websocketService.connect(wsUrl);
        console.log("WebSocket connection established. You can now send messages.");

        // Send the login confirmation message
        websocketService.sendMessage({
            type: ActionTypes.LOGIN_CONFIRMATION,
        });

        // Set up the message handler
        websocketService.onmessage((message) => {
            console.log("Login handler received message:", message);
            handleWebSocketMessage(message);
        });
    } catch (error) {
        console.error("Failed to connect to WebSocket:", error);
    }
}

// Function to handle incoming WebSocket messages
function handleWebSocketMessage(message) {
    if (message.type === ActionTypes.WELCOME) {
        displayWelcomeMessage(message.content);
    } else if (message.type === ActionTypes.ERROR) {
        handleError(message.error);
    }
}

// Function to handle errors
function handleError(error) {
    console.error("Server error:", error);
    alert(`Server error: ${error}`);
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






// Set up event listeners for page load
document.addEventListener("DOMContentLoaded", () => {
    // Event listener for login form submission
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (event) => {
            event.preventDefault(); // Prevent default form submission
            await HandleLoginRequest(); // Trigger the login function
        });
    }




});






export default HandleLoginRequest;

