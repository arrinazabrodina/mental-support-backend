import { useState } from "react";

import "./chatbox.css";

const ChatBox = ({ sendJsonMessage, storage, setStorage, chat, setChat }) => {
  const [message, setMessage] = useState("");

  const handleKeyDown = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (message.trim() !== "") {
      sendJsonMessage({ id: chat, text: message });
      setStorage(() => [
        ...storage,
        { message, date: new Date(), author: { name: "Me" } },
      ]);
      setMessage("");
    }
  };

  return (
    <form id="chatbox" onSubmit={handleSubmit}>
      <textarea
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        onKeyDown={handleKeyDown}
      ></textarea>
      <input type="submit" value="Send" />
      <button onClick={() => setChat(0)}>Chose chat</button>
    </form>
  );
};

export default ChatBox;
