import Cookies from 'js-cookie';

const access_token_name = 'access_token';
const refresh_token_name = 'refresh_token';

function setTokens(accessToken, refreshToken) {
  Cookies.set(access_token_name, accessToken, { expires: Infinity, sameSite: 'strict', secure: true });
  Cookies.set(refresh_token_name, refreshToken, { expires: Infinity, sameSite: 'strict', secure: true });
}

function removeTokens(token) {
  if (token === access_token_name) {
    Cookies.remove(access_token_name);
  } else if (token === refresh_token_name) {
    Cookies.remove(refresh_token_name);
  } else {
    Cookies.remove(access_token_name);
    Cookies.remove(refresh_token_name);
  }
}

function getTokens() {
  const cookies = {
    access: Cookies.get(access_token_name),
    refresh: Cookies.get(refresh_token_name),
  }

  return cookies;
}

async function refreshTokens() {
  const jsonData = {
    access_token: getTokens().access,
    refresh_token: getTokens().refresh,
  }

  const jsonDataString = JSON.stringify(jsonData);
  
  try {
    const response = await fetch(`${process.env.REACT_APP_BACKEND_DOMAIN}/api/refresh`, {
      method: 'POST',
      body: jsonDataString,
    });

    const data = await response.json();
    // console.log(data);
    // ARSENIY WHAT THE FUCK? WHY R U LOG FUCKING TOKKENS
    // ITS A PRIVATE U KNOW PRIVATE TOKKENS

    if (data.error === "null") {
      setTokens(data.access_token, data.refresh_token);
      return "success"
    } else {
      return null
    }

  } catch(err) {
    console.log(err);
  }

  // ARSENIY INSERT FUCKING RETURN ON FUNCTION END
  // I GOT FUCKING UNDEFINED IN VAR AND CANT REALLY
  // UNDERSTAND WHATS GOING INSIDE PROGRAM
  return null;
}

export { setTokens, removeTokens, getTokens, refreshTokens }
