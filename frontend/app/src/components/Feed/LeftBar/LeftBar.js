import './LeftBar.css';
import NewEntry_Icon from './assets/NewEntry.svg';

function LeftBar() {
  return (
    <div className="leftbar-wrapper">
      <div className="leftbar">
        <button>
          <img src={NewEntry_Icon}/>
        </button>
        <div className="leftbar-sep"></div>
        <button>
          <img />
        </button>
        <button>
          <img />
        </button>
        <button>
          <img />
        </button>
      </div>
    </div>
  )
}

export default LeftBar;
