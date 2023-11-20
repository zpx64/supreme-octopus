import { makeAutoObservable, runInAction } from 'mobx';

class NotificationStore {
  notifications = [];
  count = 0;
  normalTime = 5000;
  timeAdd = 5000;

  constructor() {
    makeAutoObservable(this);
  }

  addNotification(message, status, time) {
    const notification = { message: message, status, id: Date.now() };
    this.timeAdd = time;

    runInAction(() => {
      this.notifications.push(notification);
    });

    if (this.notifications.length > 5) {
      this.timeAdd += 1000;
      setTimeout(() => {
        this.removeNotification(notification.id);
      }, this.normalTime + this.timeAdd);
    } else {
      setTimeout(() => {
        this.removeNotification(notification.id);
      }, this.normalTime);
    }
  }

  removeNotification(id) {
    this.notifications = this.notifications.filter(notification => notification.id !== id);
  }
}

const notificationStore = new NotificationStore();
export default notificationStore;
