import { useState } from "react";

function User({ setShowModal }) {
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");

  function handler(e) {
    e.preventDefault();

    const userData = {
      name,
      email,
    };

    localStorage.setItem(
      "user",
      JSON.stringify(userData)
    );

    console.log("Saved:", userData);

    // Modal close
    setShowModal(false);
  }

  return (
    <form onSubmit={handler}>
      <h2>User Registration</h2>

      <input
        type="text"
        placeholder="Enter Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />

      <br />
      <br />

      <input
        type="email"
        placeholder="Enter Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />

      <br />
      <br />

      <button type="submit">
        Save
      </button>
    </form>
  );
}

export default User;