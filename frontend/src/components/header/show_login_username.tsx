import { loginUserState } from "@src/utils/atom";
import React from "react";
import { useRecoilValue } from 'recoil';

export const ShowLoginUserName: React.FC = () => {
  const userName: string = useRecoilValue(loginUserState);

  return (
    <>
      <p className="font-bold text-pink-800 text-lg">{userName}</p>
    </>
  );
}
