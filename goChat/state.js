// state.js

class StateManager {
    constructor() {
        this.state = {
            user: null, // To store user information
            messages: [], // To store messages
            token: null, // To store JWT token
        };
        this.listeners = []; // To hold listeners for state changes
    }

    // Method to set user information
    setUser(user) {
        this.state.user = user;
        this.notifyListeners();
    }

    // Method to set messages
    setMessages(messages) {
        this.state.messages = messages;
        this.notifyListeners();
    }

    // Method to set the token
    setToken(token) {
        this.state.token = token;
        this.notifyListeners();
    }

    // Method to get the current state
    getState() {
        return this.state;
    }

    // Method to add listeners for state changes
    addListener(listener) {
        this.listeners.push(listener);
    }

    // Method to notify all listeners of state changes
    notifyListeners() {
        this.listeners.forEach(listener => listener(this.state));
    }
}

// Create a single instance of the StateManager
const stateManager = new StateManager();
export default stateManager;
