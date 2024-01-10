import React, { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import notificationStore from './NotificationsStore';
import './Notifications.css'

const NotificationBuilder = ({ message, status, removeNotification, id }) => {
  const getStatusClass = () => {
    switch (status) {
      case 'success':
        return 'notification-status-green';
      case 'warn':
        return 'notification-status-yellow';
      default:
        return 'notification-status-red';
    };
  };

  const notificationStatus = `NotificationStatus ${getStatusClass()}`;
  
  return (
    <div id={id} className="notification">
      <div className="notification-header">
        <div className={notificationStatus}></div>
        <button className="notification-close-button" onClick={removeNotification}></button>
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
    <div className="notification-container">
      {sliced.map(notification => (
        <NotificationBuilder 
          key={notification.id}
          status={notification.status}
          message={notification.message}
          id={notification.id}
          removeNotification={() => {
           notificationStore.removeNotification(notification.id)
          }}
        />
      ))}

    </div>
  );
})
  // <button onClick={() => notificationStore.addNotification("asd", "err", timeAdd)}>Notification</button>

export default Notifications
