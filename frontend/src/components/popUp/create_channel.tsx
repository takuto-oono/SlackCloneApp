import { postChannel } from "@src/fetchAPI/channel";
import React, { useState } from "react";
import { useNavigate, useParams } from 'react-router-dom';
import { DialogTitle, DialogContent, DialogActions, Dialog, Button } from '@mui/material';

const CreateChannelForm = () => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const navigate = useNavigate();

  const handleOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  const isPrivateChangeTrue = () => {
    setIsPrivate(true);
  };
  const isPrivateChangeFalse = () => {
    setIsPrivate(false);
  };


  const handleSubmit = (e: React.FormEvent<HTMLFormElement> ) => {
    e.preventDefault();
    console.log("create channel");
    let channel = { name: name, description: description, is_private: isPrivate, workspace_id: Number(workspaceId) };
    postChannel(channel);
    setOpen(false);
    // チャンネルのリストを更新する(Todo)
  };

  return (
    <div>
      <div>
        <Button onClick={handleOpen}>
          <p style={{ color: 'black'}}>新しいチャンネルを作成</p>
        </Button>
      </div>
      <Dialog open={open} onClose={handleClose}>
        <form onSubmit={handleSubmit}>
          <DialogTitle>Create a channel</DialogTitle>
          <DialogContent>
            <label>
              名前
              <input type="text" value={name} name="name" onChange={nameChange} maxLength={80} required />
            </label>
            <fieldset>
              <legend>可視性</legend>
              <label>
                <input type="radio" name="isPrivate" onChange={isPrivateChangeTrue} />
                  プライベート : 特定のメンバーのみ
              </label><br />
              <label>
                <input type="radio" name="isPrivate" onChange={isPrivateChangeFalse} checked/>
                  パブリック : Slack 内の全員
              </label>
            </fieldset>
          </DialogContent>
          <DialogActions>
            <Button variant="outlined" onClick={handleClose}>閉じる</Button>
            <Button type="submit" variant="contained" color="success">作成</Button>
          </DialogActions>
        </form>
      </Dialog>
    </div>
  );
};

export default CreateChannelForm;
