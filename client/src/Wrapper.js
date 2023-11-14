/* eslint-disable react-hooks/exhaustive-deps */
import useWebSocket from "react-use-websocket";
import { useEffect, useState, useRef } from "react";

import MessageContainer from "./MessageContainer";
import ChatBox from "./ChatBox";

const WS_URL = `ws://${process.env.REACT_APP_HOST}:${process.env.REACT_APP_Port}/room`;

const getMessages = async (chat, page) => {
  try {
    const response = await fetch(
      `http://${process.env.REACT_APP_HOST}:${process.env.REACT_APP_Port}/messages?chatId=${chat}&page=${page}`
    );
    if (response.ok) {
      const { objects: messages, metadata } = await response.json();
      return [messages, metadata];
    }
    return [];
  } catch (error) {
    console.log(error);
  }
};

const Wrapper = ({ chat, setChat }) => {
  const [storage, setStorage] = useState([]);
  const [page, setPage] = useState(1);
  const metadata = useRef({});

  const { sendJsonMessage } = useWebSocket(WS_URL, {
    onOpen: () => {
      console.log("WebSocket connection established.");
    },
    onMessage: (e) => {
      const { data } = JSON.parse(e.data);
      if (data.chat.id === chat) setStorage(() => [...storage, data]);
    },
    onClose: () => {
      console.log("WebSocket connection closed.");
    },
  });

  useEffect(() => {
    getMessages(chat, page).then(([messages, meta]) => {
      metadata.current = meta;
      setStorage([...messages.reverse(), ...storage]);
    });
  }, [page]);

  return (
    <>
      <MessageContainer
        setPage={setPage}
        page={page}
        dataset={storage}
        metadata={metadata.current}
      ></MessageContainer>
      <ChatBox
        chat={chat}
        storage={storage}
        setStorage={setStorage}
        sendJsonMessage={sendJsonMessage}
        setChat={setChat}
      ></ChatBox>
    </>
  );
};

export default Wrapper;
