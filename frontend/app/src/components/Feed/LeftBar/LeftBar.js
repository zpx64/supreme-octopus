import './LeftBar.css';
import NewEntryIcon from './assets/NewEntryIcon';

function LeftBar({ setPostWindowSwitch, PostWindowSwitch }) {
  const SwitchPostCreationWindow = () => {
    setPostWindowSwitch(!PostWindowSwitch);
  }

  return (
    <div className="leftbar-wrapper">
      <div className="leftbar">
        <button className={PostWindowSwitch ? "leftbar-active-button" : "leftbar-inactive-button" } onClick={SwitchPostCreationWindow}>
          <NewEntryIcon fillColor={PostWindowSwitch ? "white" : "black"} />
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
