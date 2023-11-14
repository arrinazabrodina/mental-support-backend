import { useState, useEffect } from "react";

import "./selectionChat.css";
const fetchChats = async () => {
  console.log("Fetching this url: ", `http://${process.env.REACT_APP_HOST}:${process.env.REACT_APP_Port}/chats`)
  try {
    const response = await fetch(`http://${process.env.REACT_APP_HOST}:${process.env.REACT_APP_Port}/chats`);
    const chats = await response.json();

    return chats;
  } catch (e) {
    console.log(e);
    return [];
  }
};

const SelectionChat = ({ setChat }) => {
  const [chats, setChats] = useState([]);

  useEffect(() => {
    fetchChats().then((fetchedChats) => {
      setChats(fetchedChats);
    });
  }, []);
  


  return chats.length !== 0 ? chats.map((chat, index) => (
    <div className="chat-item" onClick={() => setChat(chat.id)} key={index}>
      {chat.name}
    </div> 
  )): <div style={{padding: "10px"}}>There is no chats</div>;
};

export default SelectionChat;
