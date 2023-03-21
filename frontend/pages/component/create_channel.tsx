import React, { useState } from "react";
import { postChannel, currentChannel } from "../fetchAPI/create_channel";
import { getChannelsByWorkspaceId, Channel } from "../fetchAPI/channel";
import Link from "next/link";

export default function CreateChannel(props: any) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);

  const nameChange = (e: any) => {
    setName(e.target.value);
  };

  const descriptionChange = (e: any) => {
    setDescription(e.target.value);
  };

  const privateChange = (e: any) => {
    setIsPrivate(isPrivate ? false : true);
  };

  const handleCreate = () => {
    console.log("create");
    let channel: currentChannel = {
      name: name,
      description: description,
      is_private: isPrivate,
      workspace_id: props.workspace_id,
    };
    postChannel(channel);
  };

  return (
    <div className="CreateChannel">
      <h2>Create Channel</h2>
      <label>
        チャンネルの名前
        <input type="text" value={name} name="name" onChange={nameChange} />
      </label>
      <br />
      <label>
        チャンネルの説明
        <input
          type="text"
          value={description}
          name="description"
          onChange={descriptionChange}
        />
      </label>
      <br />
      <label>
        private
        <input
          type="checkbox"
          id="isPrivate"
          value="isPrivate"
          checked={isPrivate}
          onChange={privateChange}
        />
      </label>
      <br />
      <button onClick={handleCreate}>作成</button>
      <br />

      {/* テスト用 */}
      <Link href="/">
        <button>ログイン画面へ</button>
      </Link>
    </div>
  );
}
