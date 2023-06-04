import { loginUserState } from "../main/user"
import React from "react";
import { useRecoilValue } from 'recoil';

export const ShowLoginUserName: React.FC = () => {
  const userName: string = useRecoilValue(loginUserState);
  console.log(loginUserState)
  console.log(userName);

  return (
    <>
      <p>{userName}</p>
    </>
  );
}
