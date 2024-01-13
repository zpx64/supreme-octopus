import UpVoteIcon from './assets/UpVote.js';
import DownVoteIcon from './assets/DownVote.js';

function UpVoteButton({ handleUpVote, postId, buttonStatus, reqAction }) {
    const ACTIVE_BG = "#377044";
    
    return (
        <button 
            onClick={() => handleUpVote(postId)}
            style={buttonStatus == reqAction ? { backgroundColor: ACTIVE_BG } : {}}
        >
            <UpVoteIcon fillColor={buttonStatus == reqAction ? "white" : ""} className="post-action-vote" />
        </button> 
    )
}

function DownVoteButton({ handleDownVote, postId, buttonStatus, reqAction }) {
    const ACTIVE_BG = "#793636";
    return (
        <button 
            onClick={() => handleDownVote(postId)}
            style={buttonStatus == reqAction ? { backgroundColor: ACTIVE_BG } : {}}
        >
            <DownVoteIcon fillColor={buttonStatus == reqAction ? "white" : ""} className="post-action-vote" />
        </button> 
    )
}

export { UpVoteButton, DownVoteButton };