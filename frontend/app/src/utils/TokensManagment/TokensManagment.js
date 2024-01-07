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

const fetchTokens = async () => {
  const cookies = {
    access: "asd",
    refresh: "123",
  }

  const jsonData = JSON.stringify(cookies);
  
  try {
    const response = await fetch('http://localhost:80/api/refresh', {
      method: 'POST',
      body: jsonData,
    });

    const data = await response.json();

    if (data[1] === "null") {
      setTokens(data[0], data[2]);
    } else {
      console.log("Error on tokenFetch");
    }
    console.log(data);

  } catch(err) {
    console.log(err);
  }
}

export { setTokens, removeTokens, getTokens, fetchTokens }
