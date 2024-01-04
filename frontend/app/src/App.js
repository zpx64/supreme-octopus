import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import Feed from './components/Feed/Feed';
import EnterAccount from './components/EnterAccount/EnterAccount';
import Background from './components/Background/Background';
import Notifications from 'utils/Notifications/Notifications';

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
        <EnterAccount action="login" />
        <Background />
      </>
    },
    {
      path: "/signup",
      element:
      <>
        <EnterAccount action="signup" />
        <Background />
      </>
    },
  ]);
    
  return (
    <>
      <Notifications />
      <RouterProvider router={router} />
    </>
  );
}

export default Page;
