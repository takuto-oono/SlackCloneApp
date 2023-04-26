import React from 'react';
import classes from '@styles/Home.module.css'

import { createRoot } from 'react-dom/client';
import { ProSidebarProvider } from 'react-pro-sidebar';
import { RouterConfig } from './route';
import Header from './header';
import SideNav1 from './sideNav1';
import SideNav2 from './sideNav2';

export default function Home() {
  if (typeof window === 'object') {
    const rootElement = document.getElementById('__next')!;
    const root = createRoot(rootElement);
    root.render(
      <React.StrictMode>
        <Header />
        <div className={classes.container}>
          <div className={classes.item}>
            <SideNav1 />
          </div>
          <div className={classes.item}>
            <ProSidebarProvider>
              <SideNav2 />
            </ProSidebarProvider>
          </div>
          <div className={classes.item}>
            <RouterConfig />
          </div>
        </div>
      </React.StrictMode>
    );
  }
}
