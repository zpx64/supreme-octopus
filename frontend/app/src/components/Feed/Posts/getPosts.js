import notificationStore from "utils/Notifications/notificationsStore";
import { getTokens } from "utils/TokensManagment/TokensManagment";

async function getPosts() {
    const jsonData = {
        "access_token": getTokens().access,
        "offset": 0,
        "limit": 35
    }

    const jsonDataString = JSON.stringify(jsonData);

    try {
        const response = await fetch(`${process.env.REACT_APP_BACKEND_DOMAIN}/api/list_posts`, {
            method: 'POST',
            body: jsonDataString
        })

        const data = await response.json();

        if (data.error === "null") {
            return data.posts;
        } else {
            return false;
        }
    } catch(error) {
        return false;
    }
}

export { getPosts };