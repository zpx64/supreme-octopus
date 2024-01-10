import Cookies from 'js-cookie';

const access_token_name = 'access_token';
const refresh_token_name = 'refresh_token';

const setTokens = (accessToken, refreshToken) => {
  Cookies.set(access_token_name, accessToken, { expires: Infinity, sameSite: 'strict', secure: true });
  Cookies.set(refresh_token_name, refreshToken, { expires: Infinity, sameSite: 'strict', secure: true });
}

const removeTokens = (token) => {
  if (token === access_token_name) {
    Cookies.remove(access_token_name);
  } else if (token === refresh_token_name) {
    Cookies.remove(refresh_token_name);
  } else {
    Cookies.remove(access_token_name);
    Cookies.remove(refresh_token_name);
  }
}

const getTokens = () => {
  const cookies = {
    access: Cookies.get(access_token_name),
    refresh: Cookies.get(refresh_token_name),
  }

  return cookies;
}

const refreshTokens = async () => {
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
    console.log(data);

    if (data.error === "null") {
      setTokens(data[0], data[1]);
      return "success"
    } else {
      return null
    }

  } catch(err) {
    console.log(err);
  }
}

export { setTokens, removeTokens, getTokens, refreshTokens }
