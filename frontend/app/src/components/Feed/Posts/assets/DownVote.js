function DownVoteIcon({ fillColor }) {
    let color = "#793636";

    if (fillColor) {
       color = fillColor; 
    } 

    return (
        <svg width="26" height="26" viewBox="0 0 26 26" fill="none" xmlns="http://www.w3.org/2000/svg" shapeRendering="crispEdges">
            <rect x="13.5042" y="18.0229" width="1.60839" height="1.60839" transform="rotate(-180 13.5042 18.0229)" fill={color}/>
            <rect x="11.8958" y="16.4145" width="1.60839" height="1.60839" transform="rotate(-180 11.8958 16.4145)" fill={color}/>
            <rect x="15.1125" y="16.4145" width="1.60839" height="1.60839" transform="rotate(-180 15.1125 16.4145)" fill={color}/>
            <rect x="10.2874" y="14.8061" width="1.60839" height="1.60839" transform="rotate(-180 10.2874 14.8061)" fill={color}/>
            <rect x="16.7209" y="14.8061" width="1.60839" height="1.60839" transform="rotate(-180 16.7209 14.8061)" fill={color}/>
            <rect x="18.3293" y="13.1977" width="1.60839" height="1.60839" transform="rotate(-180 18.3293 13.1977)" fill={color}/>
            <rect x="19.9377" y="11.5893" width="1.60839" height="1.60839" transform="rotate(-180 19.9377 11.5893)" fill={color}/>
            <rect x="21.5461" y="9.98093" width="1.60839" height="1.60839" transform="rotate(-180 21.5461 9.98093)" fill={color}/>
            <rect x="19.9377" y="9.98093" width="14.4755" height="1.60839" transform="rotate(-180 19.9377 9.98093)" fill={color}/>
            <rect x="8.67896" y="13.1977" width="1.60839" height="1.60839" transform="rotate(-180 8.67896 13.1977)" fill={color}/>
            <rect x="7.07056" y="11.5893" width="1.60839" height="1.60839" transform="rotate(-180 7.07056 11.5893)" fill={color}/>
            <rect x="18.3293" y="11.5893" width="11.2587" height="1.60839" transform="rotate(-180 18.3293 11.5893)" fill={color}/>
            <rect x="16.7209" y="13.1977" width="8.04195" height="1.60839" transform="rotate(-180 16.7209 13.1977)" fill={color}/>
            <rect x="15.1125" y="14.8061" width="4.82517" height="1.60839" transform="rotate(-180 15.1125 14.8061)" fill={color}/>
            <rect x="13.5042" y="16.4145" width="1.60839" height="1.60839" transform="rotate(-180 13.5042 16.4145)" fill={color}/>
            <rect x="5.46216" y="9.98093" width="1.60839" height="1.60839" transform="rotate(-180 5.46216 9.98093)" fill={color}/>
        </svg>
    )
}

export default DownVoteIcon;
