import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getTokens } from '../TokensManagment/TokensManagment';

function Feed() {
  const navigate = useNavigate();
  
  useEffect(() => {
    if (!getTokens().access || !getTokens().refresh) {
      navigate("/login");
    }
  }, [navigate]);
  
  return (
    <>
      <p style={{width: "100%", position: "absolute", top: "50%", textAlign: "center"}}>It's Your Feed</p>
    </>
  )
}

export default Feed;
