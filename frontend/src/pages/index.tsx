import React from 'react';
import classes from '@styles/Home.module.css'

import { createRoot } from 'react-dom/client';
import { RouterConfig } from './route';
import Header from './header';

export default function Home() {
  if (typeof window === 'object') {
    const rootElement = document.getElementById('__next')!;
    const root = createRoot(rootElement);
    root.render(
      <React.StrictMode>
        <Header />
        <div className={classes.container}>
          <div className={classes.item}>
            <RouterConfig />
          </div>
        </div>
      </React.StrictMode>
    );
  }
}
