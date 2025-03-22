            function HandleLogin() {
            fetch("http://localhost:8080/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                  username: "adminuser",
                  password: "adminpassword"
                }),
              })
                .then((response) => {
                  console.log("Response Status:", response.status);
                  if (!response.ok) {
                    throw new Error("Login failed!");
                  }
                  return response.json();
                })
                .then((data) => {
                  console.log("Response Data:", data)
                  
                  // Store the token for future use
                  localStorage.setItem("jwtToken", data.accessToken);
                  console.log("Stored Token:", localStorage.getItem("jwtToken"));

                  connectWebSocket(data.accessToken); // Connect WebSocket after login
                })
                .catch((error) => console.error("Login Error:", error));

              
            }

            // Function to establish WebSocket connection with JWT token
         function connectWebSocket() {
              const token = import.meta.env.VITE_ACCESS_TOKEN;
              const socket = new WebSocket(`ws://localhost:8080/ws?accessToken=${token}`);
              console.log("Access Token:", token);

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
              HandleLogin();
            });
    