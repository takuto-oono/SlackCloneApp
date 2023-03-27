import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './_app';
import { BrowserRouter } from 'react-router-dom';
import { CookiesProvider } from 'react-cookie';


export default function Index() {
  if (typeof window === 'object') {
    const rootElement = document.getElementById('__next')!;
    const root = createRoot(rootElement);
    root.render(
      <React.StrictMode>
        <CookiesProvider>
          <BrowserRouter>
            <App />
          </BrowserRouter>
        </CookiesProvider>
      </React.StrictMode>
    );
  }
}
