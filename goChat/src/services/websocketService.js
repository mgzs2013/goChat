class WebSocketService {
    constructor() {
        this.socket = null;
        this.messageCallback = null;
    }

    
    connect(wsUrl) {
        return new Promise((resolve, reject) => {
            
            console.log("Attempting to connect WebSocket with access token:", wsUrl);
            if (!wsUrl) {
                console.error("[ERROR] No websocket URL found.");
                return reject("No access token found");
            }

            // Initialize the WebSocket instance with your URL
            this.socket = new WebSocket(wsUrl);

            this.socket.onopen = () => {
                console.log("[DEBUG] WebSocket connected");
                resolve(this.socket);
            };

            this.socket.onerror = (error) => {
                console.error("[ERROR] WebSocket error:", error);
                reject(error);
            };

            this.socket.onclose = (event) => {
                console.warn("[DEBUG] WebSocket connection closed. Code:", event.code, "Reason:", event.reason);
            };

            this.socket.onmessage = (event) => {
                const message = JSON.parse(event.data);
                console.log("[DEBUG] Message from server:", event.data);
                this.displayMessage(message); // Call displayMessage as a method of the class
            };
        });
    }

    displayMessage(message) {
        const messagesContainer = document.getElementById("messages");
        const messageElement = document.createElement("div");
        messageElement.textContent = `${message.sender_id}: ${message.content} (${message.timestamp})`;
        messagesContainer.appendChild(messageElement);
    }

    sendMessage(message) {
        if (!this.socket) {
            console.error("WebSocket is not initialized.");
            return; // Exit the function if the socket is not initialized
        }
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            const payload = {
                sender_id: 2, 
                content: message,
                timestamp: new Date().toISOString(),
            };
            this.socket.send(JSON.stringify(payload));
        } else {
            console.error("WebSocket is not open. Ready state:", this.socket.readyState);
        }
    }

    onmessage(callback) {
        this.messageCallback = callback; // Store the callback
        this.socket.onmessage = (event) => {
            // Call the provided callback with the parsed message data
            const message = JSON.parse(event.data); // Parse the message data
            this.messageCallback(message); // Call the callback with the message
        };
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }
}

// Create a singleton instance
const websocketService = new WebSocketService();

export default websocketService;
