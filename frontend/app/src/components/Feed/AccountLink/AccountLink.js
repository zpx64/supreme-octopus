import User from './assets/img/User.svg';

import './AccountLink.css'

function AccountLink() {
  return (
    <>
      <div className="account-icon-wrapper">
        <div>
          <img src={User} alt="" />
        </div>
      </div>
    </>
  )
}

export default AccountLink;
