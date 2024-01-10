const { getTokens } = require("utils/TokensManagment/TokensManagment");

async function sendVoteData(postid, action) {
    const jsonData = {
        "access_token": getTokens().access,
        "post_id": postid,
        "action": action
    };

    const jsonDataString = JSON.stringify(jsonData);

    try {
        const response = await fetch(`${process.env.REACT_APP_BACKEND_DOMAIN}/api/post_vote`, {
            method: 'POST',
            body: jsonDataString,
        })

        const data = await response.json();

        if (data.error === "null") {
            return true
        } else {
            return false
        }
    } catch(error) {
        console.error(error);
        return false;
    }
}

export { sendVoteData };