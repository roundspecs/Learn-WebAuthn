const registerBtn = document.querySelector("#register-btn")
const loginBtn = document.querySelector("#login-btn")
const statusText = document.querySelector("#status-text")
const usernameField = document.querySelector("#username-field")

registerBtn.addEventListener('click', register)
loginBtn.addEventListener('click', login)

function setStatus(status) {
  statusText.innerHTML = status
}

function register() {
  setStatus("Registering...")
}

function login() {
  setStatus("Logging in...")
}
