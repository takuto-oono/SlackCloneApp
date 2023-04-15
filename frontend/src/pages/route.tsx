import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './_app';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { CookiesProvider } from 'react-cookie';
import Login from './main/login_form';
import SignUp from './main/signUp_form';
import IndexW from './workspace/workspace';
import CreateW from './workspace/create';
import ShowW from './workspace/show/[id]';
import Top from './main/Top';
import Header from './header';


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
