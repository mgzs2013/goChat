
function handleLogout() {
    // Clear username and token from local storage
    localStorage.removeItem("username");
    localStorage.removeItem("jwtToken");

    // Redirect to login page
    window.location.href = '/index.html'; // Update this path as needed
}

// Set up event listener for logout link
document.addEventListener("DOMContentLoaded", () => {
    const logoutLink = document.getElementById("logout-link");
    if (logoutLink) {
        logoutLink.addEventListener("click", (event) => {
            event.preventDefault(); // Prevent default link behavior
            handleLogout(); // Call the logout function
        });
    }
});
