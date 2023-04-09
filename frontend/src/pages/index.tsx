import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './_app';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { CookiesProvider } from 'react-cookie';
import Login from './login_form';
import SignUp from './signUp_form';
import IndexW from './workspace/workspace';
import CreateW from './workspace/create';
import ShowW from './workspace/show/[id]';
import Top from './Top';


export default function Home() {
  if (typeof window === 'object') {
    const rootElement = document.getElementById('__next')!;
    const root = createRoot(rootElement);
    root.render(
      <React.StrictMode>
        <CookiesProvider>
          <BrowserRouter>
            <p>index</p>
            <Routes >
              <Route path="/">
                <Route index element={<Top />} />
                <Route path="login_form" element={<Login />} />
                <Route path="signUp_form" element={<SignUp />} />
                <Route path="workspace" >
                  <Route index element={<IndexW />} />
                  <Route path="create" element={<CreateW />} />
                  <Route path={"show/:id"} element={<ShowW />}/>
                </Route>
              </Route>
            </Routes>
            {/* <App /> */}
          </BrowserRouter>
        </CookiesProvider>
      </React.StrictMode>
    );
  }
}
