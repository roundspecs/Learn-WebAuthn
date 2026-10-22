const registerBtn = document.querySelector("#register-btn");
const loginBtn = document.querySelector("#login-btn");
const statusText = document.querySelector("#status-text");
const usernameField = document.querySelector("#username-field");

registerBtn.addEventListener("click", register);
loginBtn.addEventListener("click", login);

function setStatus(status) {
  statusText.innerHTML = status;
}

async function register() {
  const username = usernameField.value;
  if (username == "") {
    setStatus("Username is required");
    return;
  }

  setStatus("Registering as " + username + "...");
  const response = await fetch("/register/begin", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: username }),
  });
  if (!response.ok) {
    setStatus("Error");
    return;
  }
  const options = await response.json();
  setStatus("Got challenge: " + options.publicKey.challenge);
  options.publicKey.challenge = Uint8Array.fromBase64(
    options.publicKey.challenge,
    { alphabet: "base64url" },
  );
  options.publicKey.user.id = Uint8Array.fromBase64(options.publicKey.user.id, {
    alphabet: "base64url",
  });
  const data = await navigator.credentials.create(options);

  const payload = {
    type: data.type,
    id: data.id,
    rawId: new Uint8Array(data.rawId).toBase64({ alphabet: "base64url" }),
    response: {
      attestationObject: new Uint8Array(
        data.response.attestationObject,
      ).toBase64({ alphabet: "base64url" }),
      clientDataJSON: new Uint8Array(data.response.clientDataJSON).toBase64({
        alphabet: "base64url",
      }),
    },
  };
  response = await fetch("/register/finish", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    setStatus("Error");
    return;
  }
}

function login() {
  setStatus("Logging in...");
}
