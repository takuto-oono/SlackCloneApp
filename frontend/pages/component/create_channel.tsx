import React, { useState } from "react";
import { postChannel } from "pages/fetchAPI/create_channel";

function CreateChannel() {
  const [name, setName] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);

  const nameChange = (e: any) => {
    setName(e.target.value);
  };

  const privateChange = (e: any) => {
    setIsPrivate(true);
  };

  const handleCreate = () => {
    console.log("create");
    let channelName = name;
    postChannel(channelName);
  };

  return (
    <div className="CreateWorkspace">
      <h2>Create Channel</h2>
      <label>
        チャンネルの名前
        <input type="text" value={name} name="name" onChange={nameChange} />
      </label><br />
      <label>
        private
        <input type="radio" name="private" onChange={privateChange} />
      </label>
      <br />
      <button onClick={handleCreate}>作成</button>
    </div>
  );
}

export default CreateChannel;
