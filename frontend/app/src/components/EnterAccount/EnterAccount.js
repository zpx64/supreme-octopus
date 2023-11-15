import React from 'react';
import './EnterAccount.css';

function Login() {
  return (
    <>
      <div className="window">
        <div className="windowHeader">
          <p>account / login</p>
          <div></div>
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
          <button className="signup"><p>Sign Up</p></button>
        </div>
      </div>
    </>
  );
}

export default Login;
