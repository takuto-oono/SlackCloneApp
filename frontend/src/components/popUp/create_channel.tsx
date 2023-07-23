import { postChannel } from '@src/fetchAPI/channel'
import React, { useState } from 'react'
import { DialogTitle, DialogContent, DialogActions, Dialog, Button } from '@mui/material'
import { workspaceIdState } from '@src/utils/atom'
import { useRecoilValue } from 'recoil'

const CreateChannelForm = () => {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [isPrivate, setIsPrivate] = useState(false)
  const workspaceId = useRecoilValue(workspaceIdState)

  const handleOpen = () => {
    setOpen(true)
  }
  const handleClose = () => {
    setOpen(false)
  }

  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value)
  }

  const isPrivateChangeTrue = () => {
    setIsPrivate(true)
  }
  const isPrivateChangeFalse = () => {
    setIsPrivate(false)
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    await postChannel(name, description, isPrivate, workspaceId)
    setOpen(false)
    // チャンネルのリストを更新する(Todo)
  }

  return (
    <div>
      <div>
        <button
          className='px-4 py-1 border border-slate-400 hover:bg-slate-50'
          onClick={handleOpen}
        >
          <p className='text-black'>新しいチャンネルを作成</p>
        </button>
      </div>
      <Dialog open={open} onClose={handleClose}>
        <form onSubmit={handleSubmit}>
          <DialogTitle>Create a channel</DialogTitle>
          <DialogContent>
            <div className='mb-4'>
              <label className='block mb-2 font-bold'>名前</label>
              <input
                className='border border-black w-full py-2 px-3'
                type='text'
                value={name}
                name='name'
                onChange={nameChange}
                maxLength={80}
                required
              />
            </div>
            <fieldset>
              <legend className='block mb-2 font-bold'>可視性</legend>
              <label className='block'>
                <input
                  className='mr-2'
                  type='radio'
                  name='isPrivate'
                  onChange={isPrivateChangeTrue}
                />
                <span>プライベート : 特定のメンバーのみ</span>
              </label>
              <label className='block'>
                <input
                  className='mr-2'
                  type='radio'
                  name='isPrivate'
                  onChange={isPrivateChangeFalse}
                  checked
                />
                <>パブリック : Slack 内の全員</>
              </label>
            </fieldset>
          </DialogContent>
          <DialogActions>
            <Button variant='outlined' onClick={handleClose}>
              閉じる
            </Button>
            <Button type='submit' variant='contained' color='success'>
              作成
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    </div>
  )
}

export default CreateChannelForm
