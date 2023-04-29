import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Login from './main/login_form';
import SignUp from './main/signUp_form';
import SideNav2 from './sideNav2';
import SideNav1 from './sideNav1';


export const RouterConfig: React.FC = () => {
  return (
    <>
     <BrowserRouter>
        <Routes >
          <Route path="/">
            <Route index element={<Login />} />
            <Route path="signUp_form" element={<SignUp />} />
          </Route>
          <Route path="main" element={<SideNav1 />} >
            <Route path=":workspaceid" element={<SideNav2 />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </>
  );
}
