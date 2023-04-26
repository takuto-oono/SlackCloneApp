import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Login from './main/login_form';
import SignUp from './main/signUp_form';
import Top from './main/Top';


export const RouterConfig: React.FC = () => {
  return (
    <>
     <BrowserRouter>
        <Routes >
          <Route path="/">
            <Route index element={<Top />} />
            <Route path="login_form" element={<Login />} />
            <Route path="signUp_form" element={<SignUp />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </>
  );
}
