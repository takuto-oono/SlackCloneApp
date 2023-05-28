import { useState } from "react";
import { addUserInWorkspace } from "@src/fetchAPI/workspace";
import {
	DialogTitle,
	DialogContent,
	DialogActions,
	Dialog,
	Button,
} from "@mui/material";
import { User, getAllUsers } from "@src/fetchAPI/user";

interface Props {
	workspaceID: number;
}

export const AddUserInWorkspaceForm: React.FC<Props> = (props: Props) => {
	const [open, setOpen] = useState<boolean>(false);
	const [userName, setUserName] = useState<string>("");

	const handleOpen = () => {
		setOpen(true);
	};
	const handleClose = () => {
		setOpen(false);
	};

	const changeUserName = (e: React.ChangeEvent<HTMLInputElement>) => {
		setUserName(e.target.value);
	};

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		const users: User[] | null = await getAllUsers();
		if (users == null) {
			return;
		}
		for (const user of users) {
			if (user.name == userName) {
				addUserInWorkspace(props.workspaceID, Number(user.id));
				handleClose();
				return;
			}
		}
		console.log("not found user");
	};

	return (
		<>
			<div>
				<Button onClick={handleOpen}>
					<p className="text-black">新しいユーザーを追加</p>
				</Button>
			</div>
			<Dialog open={open} onClose={handleClose}>
				<form onSubmit={handleSubmit}>
					<DialogTitle>Create a channel</DialogTitle>
					<DialogContent>
						<div className="mb-4">
							<label className="block mb-2 font-bold">ユーザー名</label>
							<input
								className="border border-black w-full py-2 px-3"
								type="text"
								value={userName}
								name="name"
								onChange={changeUserName}
								maxLength={80}
								required
							/>
						</div>
					</DialogContent>
					<DialogActions>
						<Button variant="outlined" onClick={handleClose}>
							閉じる
						</Button>
						<Button type="submit" variant="contained" color="success">
							作成
						</Button>
					</DialogActions>
				</form>
			</Dialog>
		</>
	);
};
