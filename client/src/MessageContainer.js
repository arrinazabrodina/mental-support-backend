/* eslint-disable react-hooks/exhaustive-deps */
import { useRef, useEffect } from "react";

import Message from "./Message";

import "./messageContainer.css";

const MessageContainer = ({ dataset, setPage, page, metadata }) => {
  const messagesEndRef = useRef(null);
  const currentPage = useRef(null);

  const { next } = metadata;

  useEffect(() => {
    if (currentPage.current === page) {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
    currentPage.current = page;
  }, [dataset]);

  return (
    <div className="message-field">
      {next ? (
        <div
          className="load-more"
          onClick={() =>
            setPage(() => {
              currentPage.current = page;
              return page + 1;
            })
          }
        >
          Load More
        </div>
      ) : null}
      {dataset.map((data, index) => (
        <Message key={index} data={data}></Message>
      ))}
      <div ref={messagesEndRef}></div>
    </div>
  );
};

export default MessageContainer;
