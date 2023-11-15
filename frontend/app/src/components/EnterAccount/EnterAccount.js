import React, { useState } from 'react';
import './EnterAccount.css';


function SignUpScreen({onLoginClick}) {
  const [isFullNameChecked, setIsFullNameChecked] = useState(false);

  return (
    <>
      <div className="windowS">
        <div className="windowHeader">
          <p>account / signup</p>
          <a href="https://google.com/"> </a>
        </div>
        <form className="windowArea">
          <div className="section">
            <p>Login</p>
            <input type="text" id="login-input" name="login" placeholder="login" minLength="1" maxLength="256"/>
          </div>
          <div className="section">
            <p>Email</p>
            <input type="email" id="login-email" name="email" placeholder="email" minLength="1" maxLength="256"/>
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="login-password" name="password" placeholder="password" minLength="1" maxLength="256"/>
          </div>
          <div className="sectionOption">
            <input type="checkbox" id="full-name-checkbox" name="full-name" checked={isFullNameChecked} onChange={() => setIsFullNameChecked(!isFullNameChecked)} />
            <label>Full Name</label>
          </div>
          <div className={`section ${!isFullNameChecked ? 'sectionInactive' : ''}`}>
            <p>First Name</p>
            <input type="text" id="first-name-input" name="first-name" placeholder="first name" minLength="1" maxLength="256" disabled={!isFullNameChecked}/>
          </div>
          <div className={`section ${!isFullNameChecked ? 'sectionInactive' : ''}`}>
            <p>Last Name</p>
            <input type="text" id="last-name-input" name="last-name" placeholder="last name" minLength="1" maxLength="256" disabled={!isFullNameChecked} />
          </div>
        </form>
        <div className="buttonsContainer">
          <button className="login"><p>Sign Up</p></button>
          <button className="signup" onClick={onLoginClick}><p>Login</p></button>
        </div>
      </div>
    </>
  );
}

function LoginScreen({onSignUpClick}) {
  return (
    <>
      <div className="windowL">
        <div className="windowHeader">
          <p>account / login</p>
          <a href="https://google.com/"> </a>
        </div>
        <div className="windowArea">
          <div className="section">
            <p>Login</p>
            <input type="text" id="login-input" name="login" placeholder="login" minLength="1" maxLength="256"/>
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="login-input" name="password" placeholder="password" minLength="1" maxLength="256"/>
          </div>
        </div>
        <div className="buttonsContainer">
          <button className="login"><p>Login</p></button>
          <button onClick={onSignUpClick}><p>Sign Up</p></button>
        </div>
      </div>
    </>
  );
}

function EnterAccount() {
  const [renderLogin, setRenderLogin] = useState(true);

  function handleChange() {
    if (renderLogin) {
      setRenderLogin(false);
      console.log("Render is changed to: 'SignUp'")
    } else {
      setRenderLogin(true);
      console.log("Render is changed to: 'Login'")
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
