import { Channel } from "@src/fetchAPI/channel";
import { UserInWorkspace, Workspace } from "@src/fetchAPI/workspace";
import { RecoilState, atom } from "recoil";


export const workspacesState = atom<Workspace[]>({
  key: "Workspaces",
  default: [],
})

export const channelsState = atom<Channel[]>({
  key: "channels",
  default: []
})

export const usersInWState = atom<UserInWorkspace[]>({
  key: "usersInW",
  default: []
})
