import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import Feed from './components/Feed/Feed';
import EnterAccount from './components/EnterAccount/EnterAccount';
import Background from './components/Background/Background';
import Notifications from './components/Notifications/Notifications';

import './App.css';
import './styles.css'
import './assets/styles/root.css'

function Page() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <Feed />
    },
    {
      path: "/login",
      element:
      <>
        <Notifications />
        <EnterAccount action="login" />
        <Background />
      </>
    },
    {
      path: "/signup",
      element:
      <>
        <Notifications />
        <EnterAccount action="signup" />
        <Background />
      </>
    },
  ]);
    
  return (
    <>
      <RouterProvider router={router} />
    </>
  );
}

export default Page;
