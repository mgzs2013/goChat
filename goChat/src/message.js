import websocketService from './services/websocketService'; // Your WebSocket service


function initializeMessagingPage() {

    
    const accessToken = localStorage.getItem("jwtToken");
    if (accessToken) {
        console.log("accessToken from newState:", accessToken)
        const wsUrl = `ws://localhost:8080/ws?token=${accessToken}`; // WebSocket URL
        websocketService.connect(wsUrl)
            .then(() => {
                console.log("WebSocket connection established.");
            })
            .catch((error) => {
                console.error("Failed to connect to WebSocket", error);
            });
    } else {
        console.log("User is not logged in. Cannot connect to WebSocket.");
    }

    // Add event listener for message form submission
    const messageForm = document.getElementById("message-form");
    if (messageForm) {
        messageForm.addEventListener("submit", (event) => {
            event.preventDefault(); // Prevent default form submission
            sendMessage(); // Call the sendMessage function
        });
    }
}

// Function to send a message
// Function to handle sending messages
function sendMessage() {

    if (!websocketService.socket || websocketService.socket.readyState !== WebSocket.OPEN) {
        console.error("WebSocket is not connected.");
        return; // Exit if the WebSocket is not open
    }

    const messageInput = document.getElementById("message-input").value; // Get message input
    console.log("Message:", messageInput)
    if (messageInput === "") {
        console.warn("Cannot send an empty message.");
        return; // Exit if the message is empty
    }

    const payload = {
        sender_id: 2, // Example sender ID
        content: messageInput, // The actual message content
        timestamp: new Date().toISOString(), // Add a timestamp
    };

    console.log("Sending payload:", payload); // Log the payload being sent

    websocketService.sendMessage(messageInput); // Call the sendMessage method from WebSocketService
    document.getElementById("message-input").value = ""; // Clear input
}

// Event listener for Messaging page

// Set up event listeners for page load
document.addEventListener("DOMContentLoaded", () => {
    // Event listener for login form submission
    const messageForm = document.getElementById("message-form");
    if (messageForm) {
        initializeMessagingPage();
        messageForm.addEventListener("submit", async (event) => {
            event.preventDefault(); // Prevent default form submission

        });
    }

});





// Call this function when the messaging page loads

