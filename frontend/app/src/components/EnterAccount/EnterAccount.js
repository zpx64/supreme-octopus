import React, { useState } from 'react';
import notificationStore from '../Notifications/NotificationsStore';
import './EnterAccount.css';


function validateData(type, data) {
  switch (type) {
    case ("nickname"):
      if (data.length >=3 && data.length <= 256) { return true } else { return false };
    case ("email"):
        if (data.length >= 5 && data.length <= 256) { return data.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
        )} else { return false };
    case ("password"):
      if (data.length >=6 && data.length <= 256) { return true } else { return false };
    case ("name"):
      if (data.length >=2 && data.length <= 256) { return true } else { return false };
    case ("surname"):
      if (data.length >=2 && data.length <= 256) { return true } else { return false };
    default:
      return false;
  }
}

async function sendFormDataToServer(formData, fullNameEnabled) {
  const { nickname, email, password, name, surname } = formData;

  const sendData = async (jsonData) => {
    notificationStore.addNotification("Registering account...", "warn");
    try {
      const response = await fetch('http://localhost:80/api/reg', {
        method: 'POST',
        body: jsonData,
      });

      if (!response.ok) {
        notificationStore.addNotification(`HTTP Error: ${response.status}`, "err");
      }

      const data = await response.json();
      console.log(data);
    } catch(err) {
        // notificationStore.addNotification(`HTTP Error: ${err.message}`, "err");
    }
  }

  if (!fullNameEnabled) {
    if (validateData("nickname", formData.nickname) === true &&
      validateData("email", formData.email) &&
      validateData("password", formData.password) === true) {

      const dataToSend = { nickname, email, password };
      const jsonData = JSON.stringify(dataToSend);
      sendData(jsonData);
    } else {
      console.log(validateData("email", formData.email));
      notificationStore.addNotification(`Data is incorrect`, "err");
    }
  } else {
    if (validateData("nickname", formData.nickname) === true &&
      validateData("email", formData.email) &&
      validateData("password", formData.password) === true &&
      validateData("name", formData.name) === true &&
      validateData("surname", formData.surname) === true) {

      const dataToSend = { nickname, email, password, name, surname };
      const jsonData = JSON.stringify(dataToSend);
      sendData(jsonData);
    } else {
      notificationStore.addNotification(`Data is incorrect`, "err");
    }
  }
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
  
  return (
    <>
      <div className="windowL">
        <div className="windowHeader">
          <p>account / login</p>
          <a href="https://google.com/"> </a>
        </div>
        <div className="windowArea">
          <div className="section">
            <p>Email</p>
            <input type="email" id="email" name="email" placeholder="email" minLength="3" maxLength="256"/>
          </div>
          <div className="section">
            <p>Password</p>
            <input type="password" id="password" name="password" placeholder="password" minLength="6" maxLength="256"/>
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
