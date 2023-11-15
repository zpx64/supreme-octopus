import React from 'react';
import Login from './components/EnterAccount/EnterAccount';
import Background from './components/Background/Background';

import { useState } from 'react';
import './App.css';
import './styles.css'

function Page() {
  return (
    <>
      <Login />
      <Background />
    </>
  );
}

export default Page;
