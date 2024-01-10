import notificationStore from "utils/Notifications/NotificationsStore";
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
        console.log(data);

        if (data.error === "null") {
            return data.posts;
        } else {
            return null;
        }
    } catch(error) {
        console.error(error);
        return null;
    }
}

export { getPosts };