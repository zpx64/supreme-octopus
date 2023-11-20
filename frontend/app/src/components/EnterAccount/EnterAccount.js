import React, { useState } from 'react';
import Joi from 'joi';
import notificationStore from '../Notifications/NotificationsStore';
import './EnterAccount.css';


function returnValidationScheme(action) {
  if (action == "SignUpMin") {
    return Joi.object({
      nickname: Joi.string()
                .alphanum()
                .min(3)
                .max(256)
                .label('Login')
                .required(),

      email: Joi.string()
             .email({ minDomainSegments: 2, tlds: { allow: ['com', 'net'] } })
             .min(5)
             .max(256)
                .label('Email')
             .required(),

      password: Joi.string()
                .pattern(new RegExp('^[a-zA-Z0-9]{3,30}$'))
                .min(6)
                .max(256)
                .label('Password')
                .required(),
    })
  } else if (action == "SignUpFull") {
    return Joi.object({
      nickname: Joi.string()
                .alphanum()
                .min(3)
                .max(256)
                .label('Login')
                .required(),

      email: Joi.string()
             .email({ minDomainSegments: 2, tlds: { allow: ['com', 'net'] } })
             .min(5)
             .max(256)
                .label('Email')
             .required(),

      password: Joi.string()
                .pattern(new RegExp('^[a-zA-Z0-9]{3,30}$'))
                .min(6)
                .max(256)
                .label('Password')
                .required(),

      name: Joi.string()
                .alphanum()
                .min(2)
                .max(256)
                .label('First name')
                .required(),

      surname: Joi.string()
                .alphanum()
                .min(2)
                .max(256)
                .label('Last Name')
                .required(),
    })
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
        console.log(response.status);
      }

      const data = await response.json();
      console.log(data);

      if (data.error != "null") {
        if (data.error.includes("Already in db")) {
          notificationStore.addNotification("Account already registered", "err");
        } else {
          notificationStore.addNotification(data.error, "err");
        }
      } else {
        notificationStore.addNotification("Registration successful", "success");
      }
    } catch(err) {
        // TODO: Rewrite with normal error handling
        notificationStore.addNotification(err.message, "err");
    }
  }

  if (!fullNameEnabled) {
    try {
      const schema = returnValidationScheme("SignUpMin");
      const value = await schema.validateAsync({ nickname: formData.nickname,
                                                 email: formData.email,
                                                 password: formData.password });

      const dataToSend = value;
      const jsonData = JSON.stringify(dataToSend);
      sendData(jsonData);
    }
    catch (err) {
      notificationStore.addNotification(err.details[0].message, "err");
    }
    
  } else {
    try {
      const schema = returnValidationScheme("SignUpFull");
      const value = await schema.validateAsync({ nickname: formData.nickname,
                                                 email: formData.email,
                                                 password: formData.password,
                                                 name: formData.name,
                                                 surname: formData.surname });

      const dataToSend = value;
      const jsonData = JSON.stringify(dataToSend);
      sendData(jsonData);
    }
    catch (err) {
      notificationStore.addNotification(err.details[0].message, "err");
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
