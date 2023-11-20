import React, { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import notificationStore from './NotificationsStore';
import './Notifications.css'

const NotificationBuilder = ({ message, status, removeNotification, id }) => {
  const getStatusClass = () => {
    switch (status) {
      case 'success':
        return 'NotificationStatusGreen';
      case 'warn':
        return 'NotificationStatusYellow';
      default:
        return 'NotificationStatusRed';
    };
  };

  const notificationStatus = `NotificationStatus ${getStatusClass()}`;
  
  return (
    <div id={id} className="Notification">
      <div className="NotificationHeader">
        <div className={notificationStatus}></div>
        <button className="NotificationCloseButton" onClick={removeNotification}></button>
      </div>
      <p>{message}</p>
    </div>
  )
}

const Notifications = observer(() => {
  useEffect(() => {
      notificationStore.notifications = []; // Avoid directly mutating the store
  }, []);

  let sliced = notificationStore.notifications.slice(0, 5);

  return (
    <div className="NotificationContainer">
      {sliced.map(notification => (
        <NotificationBuilder key={notification.id} status={notification.status} message={notification.message} id={notification.id} removeNotification={() => {notificationStore.removeNotification(notification.id)}} />
      ))}

    </div>
  );
})
  // <button onClick={() => notificationStore.addNotification("asd", "err", timeAdd)}>Notification</button>
  // const [notifications, setNotifications] = useState([]);

  // useEffect(() => {
  //   console.log(notifications[1]);
  //   if (notifications.length > 0) {
  //     setTimeout(() => {
  //       setNotifications((prevNotifications) => prevNotifications.slice(1));
  //     }, 5000);
  //   }
  // }, [notifications]);

  // const addNotification = (message) => {
  //   // setNotifications((prevNotifications) => [...prevNotifications, message]);
  //   setNotifications((prevNotifications) => [...prevNotifications, message]);
  // }


export default Notifications
