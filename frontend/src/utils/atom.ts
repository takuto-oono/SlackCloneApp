import { Channel, UserInChannel } from "@src/fetchAPI/channel";
import { UserInWorkspace, Workspace } from "@src/fetchAPI/workspace";
import { atom } from "recoil";

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

export const usersInWState = atom<UserInWorkspace[]>({
  key: "usersInW",
  default: []
})

export const usersInCState = atom<UserInChannel[]>({
  key: "usersInC",
  default: []
})
