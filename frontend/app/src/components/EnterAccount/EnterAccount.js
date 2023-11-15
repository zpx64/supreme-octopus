import React, { useState } from 'react';
import './EnterAccount.css';


function sendFormDataToServer(formData, fullNameEnabled) {
  const { nickname, email, password, name, surname } = formData;

  // Conditionally include name and surname only if fullNameEnabled is true
  const dataToSend = fullNameEnabled
    ? { nickname, email, password, name, surname }
    : { nickname, email, password };

  const jsonData = JSON.stringify(dataToSend);

  console.log(jsonData);
}

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

  const handleChangeValue = (e) => {
    setFormData({...formData, [e.target.id]: e.target.value})
    console.log(formData);
    console.log(e.target.value);
  }

  const handleFullNameEnable = (e) => {
    setFullNameEnabled(!fullNameEnabled);
  }
  
  const handleSignUp = (e) => {
    e.preventDefault();
    sendFormDataToServer(formData, fullNameEnabled);
  }

  return (
    <>
      <div className="windowS">
        <div className="windowHeader">
          <p>account / signup</p>
          <a href="https://google.com/"> </a>
        </div>
        <form className="windowArea" onSubmit={handleSignUp}>
          <div className="section">
            <p>Login</p>
            <input type="text" id="nickname" name="login" placeholder="login" minLength="1" maxLength="256" onChange={handleChangeValue}/>
          </div>
          <div className="section">
            <p>Email</p>
            <input type="email" id="email" name="email" placeholder="email" minLength="1" maxLength="256" onChange={handleChangeValue} />
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="password" name="password" placeholder="password" minLength="1" maxLength="256" onChange={handleChangeValue} />
          </div>
          <div className="sectionOption">
            <input type="checkbox" id="full-name-checkbox" name="full-name" checked={isFullNameChecked} onChange={() => {setIsFullNameChecked(!isFullNameChecked); handleFullNameEnable()}} />
            <label>Full Name</label>
          </div>
          <div className={`section ${!isFullNameChecked ? 'sectionInactive' : ''}`}>
            <p>First Name</p>
            <input type="text" id="name" name="first-name" placeholder="first name" minLength="1" maxLength="256" disabled={!isFullNameChecked} onChange={handleChangeValue} />
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
            <input type="text" id="login" name="login" placeholder="login" minLength="1" maxLength="256"/>
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="password" name="password" placeholder="password" minLength="1" maxLength="256"/>
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
