import { makeAutoObservable, runInAction } from 'mobx';

class NotificationStore {
  notifications = [];
  count = 0;
  normalTime = 5000;
  timeAdd = 5000;
  defaultTimeAdd = 5000;

  constructor() {
    makeAutoObservable(this);
  }

  addNotification(message, status, time) {
    const notification = { message: message, status, id: Date.now() + this.count };


    if (time) {
      this.timeAdd = time;
    }
    this.count++;

    // console.log(this.timeAdd);
    // console.log(this.count);

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
        this.timeAdd = this.defaultTimeAdd;
      }, this.normalTime);
    }
  }

  removeNotification(id) {
    this.notifications = this.notifications.filter(notification => notification.id !== id);
  }
}

const notificationStore = new NotificationStore();
export default notificationStore;
