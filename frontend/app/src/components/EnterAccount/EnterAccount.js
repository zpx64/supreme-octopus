import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { sendSignUpDataToServer, sendLoginDataToServer } from './sendData';

import './EnterAccount.css';


function SignUpScreen({onLoginClick}) {
  const [isFullNameChecked, setIsFullNameChecked] = useState(false);
  const [fullNameEnabled, setFullNameEnabled] = useState(false);
  const [formData, setFormData] = useState({
    nickname: '',
    email: '',
    password: '',
    name: '',
    surname: '',
  });

  const navigator = useNavigate();

  const handleChangeValue = (e) => {
    setFormData({...formData, [e.target.id]: e.target.value})
  }

  const handleFullNameEnable = (e) => {
    setFullNameEnabled(!fullNameEnabled);
  }
  
  const handleSignUp = async(e) => {
    const result = await sendSignUpDataToServer(formData, fullNameEnabled);
    console.log(result);
    if (result) {
      navigator('/login');
    }
  }

  return (
    <>
      <div className="windowSignUp">
        <div className="windowHeader">
          <p>account / signup</p>
          <a href="https://google.com/"> </a>
        </div>
        <form className="windowArea" onSubmit={handleSignUp}>
          <div className="section">
            <p>Login</p>
            <input type="text" id="nickname" name="login" placeholder="login" minLength="3" maxLength="256" onChange={handleChangeValue}/>
          </div>
          <div className="section">
            <p>Email</p>
            <input type="email" id="email" name="email" placeholder="email" minLength="5" maxLength="256" onChange={handleChangeValue} />
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="password" name="password" placeholder="password" minLength="6" maxLength="256" onChange={handleChangeValue} />
          </div>
          <div className="sectionOption">
            <input type="checkbox" id="full-name-checkbox" name="full-name" checked={isFullNameChecked} onChange={() => {setIsFullNameChecked(!isFullNameChecked); handleFullNameEnable()}} />
            <label>Full Name</label>
          </div>
          <div className={`section ${!isFullNameChecked ? 'sectionInactive' : ''}`}>
            <p>First Name</p>
            <input type="text" id="name" name="first-name" placeholder="first name" minLength="2" maxLength="256" disabled={!isFullNameChecked} onChange={handleChangeValue} />
          </div>
          <div className={`section ${!isFullNameChecked ? 'sectionInactive' : ''}`}>
            <p>Last Name</p>
            <input type="text" id="surname" name="last-name" placeholder="last name" minLength="1" maxLength="256" disabled={!isFullNameChecked} onChange={handleChangeValue} />
          </div>
        </form>
        <div className="buttonsContainer">
          <button onClick={handleSignUp}><p>Sign Up</p></button>
          <button onClick={onLoginClick}><p>Login</p></button>
        </div>
      </div>
    </>
  );
}

function LoginScreen({onSignUpClick}) {
  const [formData, setFormData] = useState({});
  const navigator = useNavigate();

  const handleChangeValue = (e) => {
    setFormData({...formData, [e.target.id]: e.target.value});
  }

  const handleLogin = async() => {
    const result = await sendLoginDataToServer(formData);
    console.log(result);
    if (result) {
      navigator('/');
    }
  }
  
  return (
    <>
      <div className="windowLogin">
        <div className="windowHeader">
          <p>account / login</p>
          <a href="https://google.com/"> </a>
        </div>
        <div className="windowArea">
          <div className="section">
            <p>Email</p>
            <input type="email" id="email" name="email" placeholder="email" minLength="3" maxLength="256" onChange={handleChangeValue} />
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="password" name="password" placeholder="password" minLength="6" maxLength="256" onChange={handleChangeValue} />
          </div>
        </div>
        <div className="buttonsContainer">
          <button onClick={handleLogin} className="login"><p>Login</p></button>
          <button onClick={onSignUpClick}><p>Sign Up</p></button>
        </div>
      </div>
    </>
  );
}

function EnterAccount({ action }) {
  const [renderLogin, setRenderLogin] = useState(true);
  const navigator = useNavigate();

  useEffect(() => {
    if (action == "login") {
      setRenderLogin(true);
    } if (action == "signup") {
      setRenderLogin(false);
    }
  }, [action]);

  function handleChange() {
    if (renderLogin) {
      setRenderLogin(false);
      navigator("/signup");
    } else {
      setRenderLogin(true);
      navigator("/login");
    }
  }

  return (
    <>
      {renderLogin ? (
        <LoginScreen onSignUpClick={() => handleChange()} />
      ) : (
        <SignUpScreen onLoginClick={() => handleChange()} />
      )}
    </>
  )
}

export default EnterAccount;
