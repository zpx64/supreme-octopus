function NewEntryIcon({ fillColor }) {
    return (
        <svg width="20" height="25" viewBox="0 0 20 25" fill="none" xmlns="http://www.w3.org/2000/svg" shapeRendering="crispEdges">
            <rect x="1.25" y="1.25" width="1.25" height="23.75" fill={fillColor}/>
            <rect y="1.25" width="1.25" height="23.75" fill={fillColor}/>
            <rect x="2.5" y="25" width="1.25" height="15" transform="rotate(-90 2.5 25)" fill={fillColor}/>
            <rect x="2.5" y="23.75" width="1.25" height="15" transform="rotate(-90 2.5 23.75)" fill={fillColor}/>
            <rect x="12.5" y="7.5" width="1.25" height="6.25" transform="rotate(-90 12.5 7.5)" fill={fillColor}/>
            <rect x="2.5" y="2.5" width="1.25" height="12.5" transform="rotate(-90 2.5 2.5)" fill={fillColor}/>
            <rect x="1.25" y="1.25" width="1.25" height="13.75" transform="rotate(-90 1.25 1.25)" fill={fillColor}/>
            <rect y="1.25" width="1.25" height="1.25" transform="rotate(-90 0 1.25)" fill={fillColor}/>
            <rect x="15" y="3.75" width="1.25" height="1.25" transform="rotate(-90 15 3.75)" fill={fillColor}/>
            <rect x="16.25" y="3.75" width="1.25" height="1.25" transform="rotate(-90 16.25 3.75)" fill={fillColor}/>
            <rect x="17.5" y="6.25" width="1.25" height="1.25" transform="rotate(-90 17.5 6.25)" fill={fillColor}/>
            <rect x="18.75" y="6.25" width="1.25" height="1.25" transform="rotate(-90 18.75 6.25)" fill={fillColor}/>
            <rect x="12.5" y="6.25" width="2.5" height="5" transform="rotate(-90 12.5 6.25)" fill={fillColor}/>
            <rect x="12.5" y="3.75" width="1.25" height="2.5" transform="rotate(-90 12.5 3.75)" fill={fillColor}/>
            <rect x="17.5" y="6.25" width="1.25" height="18.75" fill={fillColor}/>
            <rect x="18.75" y="6.25" width="1.25" height="18.75" fill={fillColor}/>
        </svg>
    )
}

export default NewEntryIcon;