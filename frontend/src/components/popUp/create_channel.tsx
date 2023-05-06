import { postChannel } from "@src/fetchAPI/channel";
import React, { useState } from "react";
import { useParams } from 'react-router-dom';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';

const CreatChannelForm = () => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const { workspaceId } = useParams<{ workspaceId: string }>();

  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const nameChange = (e: any) => {
    setName(e.target.value);
  };

  const isPrivateChangeTrue = () => {
    setIsPrivate(true);
  };
  const isPrivateChangeFalse = () => {
    setIsPrivate(false);
  };


  const handleSubmit = () => {
    console.log("create channel");
    let channel = { name: name, description: description, is_private: isPrivate, workspace_id: parseInt(workspaceId) };
    postChannel(channel);
    setOpen(false);
  };

  return (
    <div>
      <div>
        <Button onClick={handleClickOpen}>
          <p style={{ color: 'black'}}>新しいチャンネルを作成</p>
        </Button>
      </div>
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>Create a channel</DialogTitle>
        <DialogContent>
          <label htmlFor="name">
            名前<br />
            <input type="text" value={name} name="name" onChange={nameChange} />
          </label>

          <fieldset>
            <legend>可視性</legend>
            <label htmlFor="is_private">
              <input type="radio" name="isPrivate" onChange={isPrivateChangeTrue} />
              プライベート : 特定のメンバーのみ
            </label><br />
            <label htmlFor="is_private">
              <input type="radio" name="isPrivate" onChange={isPrivateChangeFalse} checked/>
              パブリック : Slack 内の全員
            </label>
          </fieldset> 
        </DialogContent>
        <DialogActions>
          <Button variant="outlined" onClick={handleClose}>閉じる</Button>
          <Button variant="contained" color="success" onClick={handleSubmit}>作成</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default CreatChannelForm;
