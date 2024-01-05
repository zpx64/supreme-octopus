import { Buffer } from 'buffer';
import base85 from 'base85';
import { getTokens } from "utils/TokensManagment/TokensManagment"
import notificationStore from 'utils/Notifications/NotificationsStore';


const PadBuffer = (buffer) => {
  const size = buffer.length;
  const need = 4 - (size % 4);

  const buf = Buffer.alloc(need)
  for (let i = 0; i < need; ++i) {
    buf.writeUInt8(need, i);
  }

  return Buffer.concat([buffer, buf], size + need);
}

async function createPost(type, text, attachments) {
    const jsonData = {
        "access_token": getTokens().access,
        "post_type": type,
        "body": text,
        "attachments": attachments,
    }

    let postTypeName = null;
    if (type === 1) {
        postTypeName = "Article";
    } else {
        postTypeName = "Note";
    }
    
    const jsonDataString = JSON.stringify(jsonData);
    console.log(jsonDataString);
    notificationStore.addNotification('Running [2/2]: post creation phase...', 'warn');
    try {
        const response = await fetch(`${process.env.REACT_APP_BACKEND_DOMAIN}/api/post_new`, {
            method: 'POST',
            body: jsonDataString,
        })

        const data = await response.json();

        if (data.error === "null") {
            notificationStore.addNotification(`Post type "${postTypeName}" created`, 'success');
        } else {
            notificationStore.addNotification('Creation phase failed', 'err');
        }
    } catch(error) {
        console.log(error);
    }
}

async function registerImages(images) {
    let jsonData = {
        "access_token": getTokens().access,
        "images": []
    }

    const readFileAsync = (file) => new Promise((resolve) => {
        const reader = new FileReader();

        reader.onload = () => {
            const buffer = Buffer.from(reader.result);
            const result = base85.encode(PadBuffer(buffer), 'z85');
        
            resolve({
                "content_type": file.type,
                "encoded_image": result
            })
        }

        reader.readAsArrayBuffer(file);
    })

    try {
        const processedImages = await Promise.all(images.map(readFileAsync));
        jsonData.images = processedImages;
        const jsonDataString = JSON.stringify(jsonData);
        notificationStore.addNotification('Running [1/2]: images uploading phase...', 'warn');

        try {
            const response = await fetch(`${process.env.REACT_APP_BACKEND_DOMAIN}/api/post_image`, {
                method: 'POST',
                body: jsonDataString,
            })

            const data = await response.json();

            if (data.error === "null") {
                return data.writed_ids
            } else {
                console.log(data.error);
                notificationStore.addNotification('Error while processing images', 'err');
            }
        } catch(error) {
            console.log(error);
        }
    } catch(error) {
        console.log(error);
    }
}

export { registerImages, createPost }