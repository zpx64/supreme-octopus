import React from 'react';
import Login from './components/EnterAccount/EnterAccount';
import Background from './components/Background/Background';
import Notifications from './components/Notifications/Notifications';

import './App.css';
import './styles.css'
import './assets/styles/root.css'

function Page() {
  return (
    <>
      <Notifications />
      <Login />
      <Background />
    </>
  );
}

export default Page;
