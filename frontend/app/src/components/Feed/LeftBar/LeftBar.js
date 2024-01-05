import './LeftBar.css';
import NewEntry_Icon from './assets/NewEntry.svg';

function LeftBar({ setPostWindowSwitch }) {
  const SwitchPostCreationWindow = () => {
    setPostWindowSwitch(true);
  }

  return (
    <div className="leftbar-wrapper">
      <div className="leftbar">
        <button onClick={SwitchPostCreationWindow}>
          <img src={NewEntry_Icon} alt="" />
        </button>
        <div className="leftbar-sep"></div>
        <button>
          <img alt="" />
        </button>
        <button>
          <img alt="" />
        </button>
        <button>
          <img alt="" />
        </button>
      </div>
    </div>
  )
}

export default LeftBar;
