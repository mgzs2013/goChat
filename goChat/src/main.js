async function HandleLoginRequest() {
  try {
    const response = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(
        {
          username: "adminuser",
          password: "adminpassword",
        }),
    });
    console.log("Raw Response:", response);
    console.log("Response Status:", response.status);

    // Check if the response is okay (status 2xx)
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    // Try to parse the JSON response
    // const text = await response.text();
    // let data;
    // try {
    //   data = JSON.parse(text);
    // } catch (e) {
    //   console.error("Invalid JSON received from server:", e);
    //   throw new Error("Unexpected response format");
    // }
    const data = await response.json();
    console.log("Parsed Response:", data);

    console.log("Response Data:", data);

    // Handle success (e.g., store token, connect WebSocket)
    localStorage.setItem("jwtToken", data.accessToken);
    console.log("Stored Token:", localStorage.getItem("jwtToken"));
    connectWebSocket();
    console.log("accessToken:", data.accessToken);
  } catch (error) {
    console.error("Login Error:", error.message);
  }
}

// Function to establish WebSocket connection with JWT token
function connectWebSocket() {
  var secret = import.meta.env.VITE_ACCESS_TOKEN;
  console.log("secret var:", secret)
  const token = localStorage.getItem("jwtToken"); // Retrieve the stored token
  if (!token) {
    console.error("[ERROR] No access token found in local storage.");
    alert("Session expired. Please log in again.");
    return;
  }
  const socket = new WebSocket(`ws://localhost:8080/ws?accessToken=${token}`);
  console.log("Access Token:", token);
  console.log("Loaded Environment Variables:", import.meta.env);


  socket.onopen = () => {
    console.log("[DEBUG] WebSocket connected");
  };

  socket.onerror = (error) => {
    console.error("[ERROR] WebSocket error:", error);
  };

  socket.onclose = (event) => {
    console.warn(
      "[DEBUG] WebSocket connection closed. Code:",
      event.code,
      "Reason:",
      event.reason
    );
    alert("Session expired. Please log in again.");
  };

  socket.onmessage = (event) => {
    console.log("[DEBUG] Message from server:", event.data);
  };
}

// Function to send messages through WebSocket
export function sendMessage(socket) {
  const sender_id = 1; // Replace with the actual SenderID (e.g., fetched from server-side claims)
  const message = document.getElementById("message").value;
  const payload = {
    sender_id,
    content: message,
    timestamp: new Date().toISOString(), // Use ISO format for better compatibility
  };

  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(payload));
  } else {
    console.error(
      "WebSocket is not open. Ready state:",
      socket.readyState
    );
  }

  // Clear the input field
  document.getElementById("message").value = "";
}

// Trigger login when the page loads
document.addEventListener("DOMContentLoaded", () => {
  console.log("DOMContentLoaded event fired");
  HandleLoginRequest();
}
);


