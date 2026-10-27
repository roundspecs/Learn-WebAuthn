const registerBtn = document.querySelector("#register-btn");
const loginBtn = document.querySelector("#login-btn");
const statusText = document.querySelector("#status-text");
const usernameField = document.querySelector("#username-field");

let isLoggedIn = false;

registerBtn.addEventListener("click", register);
loginBtn.addEventListener("click", () => {
  if (isLoggedIn) {
    logout();
  } else {
    login();
  }
});

function setStatus(status) {
  statusText.innerHTML = status;
}

function setLoggedIn(loggedIn) {
  isLoggedIn = loggedIn;
  if (loggedIn) {
    loginBtn.textContent = "Logout";
    setStatus("Logged in");
  } else {
    loginBtn.textContent = "Login";
    setStatus("Logged out");
    usernameField.value = "";
  }
}

checkSession();

async function checkSession() {
  try {
    const response = await fetch("/session");
    if (!response.ok) return;
    const data = await response.json();
    if (data.loggedIn) {
      setLoggedIn(true);
      if (data.username) {
        usernameField.value = data.username;
      }
    }
  } catch (err) {
    console.error(err);
  }
}

async function logout() {
  setStatus("Logging out...");
  try {
    const response = await fetch("/logout", { method: "POST" });
    if (response.ok) {
      setLoggedIn(false);
    } else {
      setStatus("Error logging out");
    }
  } catch (err) {
    setStatus("Error logging out");
  }
}

async function register() {
  const username = usernameField.value;
  if (username == "") {
    setStatus("Username is required");
    return;
  }

  setStatus("Registering as " + username + "...");
  response = await fetch("/register/begin", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username: username }),
  });
  if (!response.ok) {
    setStatus(await response.text());
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
    setStatus(await response.text());
  }
  else {
    setStatus("Registered")
  }
}

async function login() {
  const username = usernameField.value;
  if (username == "") {
    setStatus("Username is required");
    return;
  }

  setStatus("Logging in as " + username + "...");
  let response = await fetch("/login/begin", {
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
  if (options.publicKey.allowCredentials) {
    options.publicKey.allowCredentials = options.publicKey.allowCredentials.map(
      (cred) => ({
        ...cred,
        id: Uint8Array.fromBase64(cred.id, { alphabet: "base64url" }),
      }),
    );
  }
  const data = await navigator.credentials.get(options);

  const payload = {
    type: data.type,
    id: data.id,
    rawId: new Uint8Array(data.rawId).toBase64({ alphabet: "base64url" }),
    response: {
      clientDataJSON: new Uint8Array(data.response.clientDataJSON).toBase64({
        alphabet: "base64url",
      }),
      authenticatorData: new Uint8Array(
        data.response.authenticatorData,
      ).toBase64({ alphabet: "base64url" }),
      signature: new Uint8Array(data.response.signature).toBase64({
        alphabet: "base64url",
      }),
      userHandle: data.response.userHandle
        ? new Uint8Array(data.response.userHandle).toBase64({
            alphabet: "base64url",
          })
        : null,
    },
  };
  response = await fetch("/login/finish", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    setStatus("Error");
    return;
  }
  setLoggedIn(true);
}
