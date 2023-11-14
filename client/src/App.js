import React, { useState } from "react";

import Wrapper from "./Wrapper";
import SelectionChat from "./SelectionChat";

function App() {
  const [chat, setChat] = useState(0);
  
  if (chat === 0) {
    return <SelectionChat setChat={setChat}></SelectionChat>;
  } else {
    return <Wrapper chat={chat} setChat={setChat}></Wrapper>;
  }
}

export default App;
