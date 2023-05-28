import React, { useEffect, useState, useRef } from "react";
import {
	getWorkspaces,
	Workspace,
	UserInWorkspace,
	getUsersInWorkspace,
} from "@fetchAPI/workspace";
import { Link } from "react-router-dom";
import { MenuItem } from "react-pro-sidebar";
import { atom, useRecoilState } from "recoil";
import Button from "@mui/material/Button";
import Popover from "@mui/material/Popover";
import { AddUserInWorkspaceForm } from "../popUp/add_user_in_workspace_form";

export const usersInWState = atom<UserInWorkspace[]>({
	key: "usersInW",
	default: [],
});

function ShowWorkspaces() {
	const [open, setOpen] = useState(false);
	const [usersInW, setUsersInW] = useRecoilState(usersInWState);

	const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);
	const divRef = useRef(null);

	const handleClickOpen = () => {
		setOpen(true);
	};

	const handleClickClose = () => {
		setOpen(false);
	};

	const getWorkspaceInfo = (workspaceId: number) => {
		getUsersInWorkspace(workspaceId).then((usersInW: UserInWorkspace[]) => {
			setUsersInW(usersInW);
		});
	};

	const list = workspaceList.map((workspace, index) => (
		<div key={index}>
			<MenuItem>
				<Link
					to={`${workspace.id}`}
					onClick={() => getWorkspaceInfo(workspace.id)}
				>
					<span>{workspace.name}</span>
				</Link>
				<div ref={divRef} className="bg-purple-200 text-pink-700">
					<Button onClick={handleClickOpen}>
						<p className="bg-purple-200 text-pink-700">
							+ ユーザーを追加
						</p>
					</Button>
					<Popover
						open={open}
						anchorEl={divRef.current}
						onClose={handleClickClose}
						anchorOrigin={{
							vertical: "bottom",
							horizontal: "left",
						}}
					>
						<AddUserInWorkspaceForm workspaceID={workspace.id} />
					</Popover>
				</div>
			</MenuItem>
		</div>
	));

	useEffect(() => {
		getWorkspaces().then((workspaces: Workspace[]) => {
			setWorkspaceList(workspaces);
		});
	}, []);

	return (
		<div className="App">
			<div>{list}</div>
		</div>
	);
}

export default ShowWorkspaces;
