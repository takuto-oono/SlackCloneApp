import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Login from "./main/login_form";
import SignUp from "./main/signUp_form";
import SideNav2 from "./sideNav2";
import SideNav1 from "./sideNav1";
import CreateW from "./main/create_workspace";
import Main from "./main";

export const RouterConfig: React.FC = () => {
	return (
    <BrowserRouter>
      <Routes>
        <Route path="/">
          <Route index element={<Login />} />
          <Route path="signUp_form" element={<SignUp />} />
          <Route path="workspace" element={<SideNav1 />}>
            <Route path="create" element={<CreateW />} />
            <Route path=":workspaceId" element={<SideNav2 />}>
              {/* ここでメインページのルーティングを設定する */}
              <Route path="channel/:channelID" element={<Main />} />
            </Route>
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
	);
};

