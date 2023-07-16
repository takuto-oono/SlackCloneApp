import React, { useState } from "react";
import { Workspace, getWorkspaces, postWorkspace } from '@fetchAPI/workspace';
import { useRouter } from "next/router";
import { userChannelsState, workspaceChannelsState, workspaceIdState, workspacesState } from "@src/utils/atom";
import { useSetRecoilState } from "recoil";
import { Channel, getChannelsInW, getUserChannelsInW } from "@src/fetchAPI/channel";

function CreateWorkspace() {
  const [name, setName] = useState("");
  const setWorkspaceId = useSetRecoilState(workspaceIdState);
  const setWorkspaces = useSetRecoilState(workspacesState);
  const setUserChannels = useSetRecoilState(userChannelsState);
  const setWorkspaceChannels = useSetRecoilState(workspaceChannelsState)

  const router = useRouter()
  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement> ) => {
    e.preventDefault();
    let workspaceName = name;
    postWorkspace(workspaceName);
    // ワークスペースのリストを更新する(Todo)
  };
  return (
    <div>
      <form onSubmit={handleSubmit}  className="px-8 py-8">
        <p className="text-2xl p-1">Create Workspace</p>
        <div className="mb-4">
        <label  className="block mb-2 font-bold">ワークスペースの名前</label>
          <input className="border border-black w-full py-2 px-3" type="text" value={ name } name="name" onChange={(e) => nameChange(e)} maxLength={50} required/>
        </div>
        <div>
          <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">作成</button>
        </div>
      </form>
    </div>
  );
}

export default CreateWorkspace;
