function UpVoteIcon({ fillColor }) {
    let color = "#377044";

    if (fillColor) {
       color = fillColor; 
    } 

    return (
        <svg width="26" height="26" viewBox="0 0 26 26" fill="none" xmlns="http://www.w3.org/2000/svg" shapeRendering="crispEdges">
            <rect x="11.8958" y="8.08112" width="1.60839" height="1.60839" fill={color}/>
            <rect x="13.5043" y="9.68951" width="1.60839" height="1.60839" fill={color}/>
            <rect x="10.2874" y="9.68951" width="1.60839" height="1.60839" fill={color}/>
            <rect x="15.1127" y="11.2979" width="1.60839" height="1.60839" fill={color}/>
            <rect x="8.67908" y="11.2979" width="1.60839" height="1.60839" fill={color}/>
            <rect x="7.07068" y="12.9063" width="1.60839" height="1.60839" fill={color}/>
            <rect x="5.46228" y="14.5146" width="1.60839" height="1.60839" fill={color}/>
            <rect x="3.85388" y="16.123" width="1.60839" height="1.60839" fill={color}/>
            <rect x="5.46228" y="16.123" width="14.4755" height="1.60839" fill={color}/>
            <rect x="16.7211" y="12.9063" width="1.60839" height="1.60839" fill={color}/>
            <rect x="18.3293" y="14.5146" width="1.60839" height="1.60839" fill={color}/>
            <rect x="7.07068" y="14.5146" width="11.2587" height="1.60839" fill={color}/>
            <rect x="8.67908" y="12.9063" width="8.04196" height="1.60839" fill={color}/>
            <rect x="10.2874" y="11.2979" width="4.82517" height="1.60839" fill={color}/>
            <rect x="11.8958" y="9.68951" width="1.60839" height="1.60839" fill={color}/>
            <rect x="19.9377" y="16.123" width="1.60839" height="1.60839" fill={color}/>
        </svg>
    )
}

export default UpVoteIcon;
