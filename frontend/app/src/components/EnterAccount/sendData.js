import Joi from 'joi';
import { setTokens } from '../TokensManagment/TokensManagment';
import notificationStore from '../Notifications/NotificationsStore';


function returnValidationScheme(action) {
  if (action === "SignUpMin") {
    return Joi.object({
      nickname: Joi.string()
                .alphanum()
                .min(3)
                .max(256)
                .label('Login')
                .required(),

      email: Joi.string()
                .email({ tlds: { allow: false } })
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
  } else if (action === "SignUpFull") {
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
  } else if (action === "Login") {
    return Joi.object({
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
  }
}



async function sendSignUpDataToServer(formData, fullNameEnabled) {
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

      if (data.error !== "null") {
        if (data.error.includes("Already in db")) {
          notificationStore.addNotification("Account already registered", "err");
        } else {
          notificationStore.addNotification(data.error, "err");
        }
        return false;

      } else {
        notificationStore.addNotification("Registration successful", "success");
        return true;
      }
    } catch(err) {
      // TODO: Rewrite with normal error handling
      notificationStore.addNotification(err.message, "err");
      return false;
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
      return sendData(jsonData);
    }
    catch (err) {
      notificationStore.addNotification(err.details[0].message, "err");
      return false;
    }
    
  } else {
    try {
      const schema = returnValidationScheme("SignUpFull");
      const value = await schema.validateAsync(
        {
          nickname: formData.nickname,
          email: formData.email,
          password: formData.password,
          name: formData.name,
          surname: formData.surname
        }
      );

      const dataToSend = value;
      const jsonData = JSON.stringify(dataToSend);
      return sendData(jsonData);
    }
    catch (err) {
      notificationStore.addNotification(err.details[0].message, "err");
      return false;
    }
  }
}

function generateDeviceId() {
  const userAgent = window.navigator.userAgent;

  const rndMin = 1000000000;
  const rndMax = 9999999999;
  const userId = Math.floor(Math.random() * (rndMax - rndMin)) + rndMin;

  const result = userAgent + "_" + userId;
  
  return result;
}

async function sendLoginDataToServer(formData) {
  formData.device_id = generateDeviceId();

  const sendData = async (jsonData) => {
    notificationStore.addNotification("Logging into your account...", "warn");
    try {
      const response = await fetch('http://localhost:80/api/login', {
        method: 'POST',
        body: jsonData,
      });

      const data = await response.json();

      if (data.error !== "null") {
        // TODO: Rewrite with normal error handling      
        notificationStore.addNotification(data.error, "err");

        return false;
      } else {
        setTokens(data.access_token, data.refresh_token);
        notificationStore.addNotification("Login successful", "success");

        return true;
      }

    } catch(err) {
      // TODO: Rewrite with normal error handling
      notificationStore.addNotification(err.message, "err");
      return false;
    }
  }

  try {
    const schema = returnValidationScheme("Login");
    await schema.validateAsync(
      {
        email: formData.email,
        password: formData.password
      }
    );

    const dataToSend = formData;
    const jsonData = JSON.stringify(dataToSend);
    return sendData(jsonData);
  }
  catch (err) {
    notificationStore.addNotification(err.details[0].message, "err");
    return false;
  }
}

export { sendSignUpDataToServer, sendLoginDataToServer }
