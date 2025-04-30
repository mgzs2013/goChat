// register.js

// Function to handle user registration
async function RegisterUserRequest() {
    

    // Get the form data
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value.trim();
    console.log("Attempting to register with username:", username);
    console.log("Attempting to register in with username:", password);

    try {
        const data = await RegisterUser(username, password); 
        console.log("Registration successful with username and password:", data);
        window.location.href = '/index.html'; // Update this path as needed
       
    } catch (error) {
        console.error("Registration Error:", error.message);
        // Handle error feedback to the user if needed
        alert("Registration failed: " + error.message);
    }
}

async function RegisterUser(username, password) {
    const API_URL = `http://localhost:8080/register`; // Ensure this is correct
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

    if (response.status === 204) {
        // No content
        throw new Error("No content returned from server.");
    }

    const data = await response.json(); // Parse the response JSON
    console.log("Parsed Response:", data); // Log the parsed response
    return data; // Return the parsed data
}

// Attach the registerUser function to the form submission

// Set up event listeners for page load
document.addEventListener("DOMContentLoaded", () => {
    // Event listener for login form submission
    const registrationForm = document.getElementById('registration-form');
    if (registrationForm) {
        
        registrationForm.addEventListener("submit", async (event) => {
            event.preventDefault(); // Prevent default form submission
            await RegisterUserRequest(); // Trigger the Registration function

        });
    }

});