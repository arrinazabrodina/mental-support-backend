import "./message.css";

const Message = ({ data }) => {
  const { author, date, message } = data;
  return (
    <div className="message">
      <div className="message-text">{message}</div>
      <div className="message-author">{author.name}</div>
      <div className="message-datetime">
        {new Date(date).toLocaleTimeString("en-US", {
          hour: "numeric",
          minute: "numeric",
          hour12: false,
        })}
      </div>
    </div>
  );
};

export default Message;
