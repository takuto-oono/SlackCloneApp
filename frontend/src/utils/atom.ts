import { Channel } from "@src/fetchAPI/channel";
import { UserInWorkspace, Workspace } from "@src/fetchAPI/workspace";
import { atom } from "recoil";

export const workspacesState = atom<Workspace[]>({
  key: "workspaces",
  default: []
})

export const channelsState = atom<Channel[]>({
  key: "channels",
  default: []
})

export const usersInWState = atom<UserInWorkspace[]>({
  key: "usersInW",
  default: []
})
