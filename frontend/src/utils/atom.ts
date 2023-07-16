import { Channel } from "@src/fetchAPI/channel";
import { UserInWorkspace, Workspace } from "@src/fetchAPI/workspace";
import { atom } from "recoil";

export const loginUserState = atom<string>({
  key: "userName",
  default: ""
})

export const workspaceIdState = atom<number>({
  key: "WorkspaceId",
  default: 0,
})

export const workspacesState = atom<Workspace[]>({
  key: "Workspaces",
  default: [],
})

export const channelsState = atom<Channel[]>({
  key: "channels",
  default: []
})

export const joinedChannelsState = atom<Channel[]>({
  key: "joinedChannels",
  default: []
})

export const usersInWState = atom<UserInWorkspace[]>({
  key: "usersInW",
  default: []
})

export const usersInCState = atom<UserInWorkspace[]>({
  key: "usersInC",
  default: []
})
