import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getTokens } from '../TokensManagment/TokensManagment';

import Scrollbar from '../Scrollbar/Scrollbar';
import AccountLink from './AccountLink/AccountLink';
import SearchBar from './SearchBar/SearchBar';
import Posts from './Posts/Posts';


function Feed() {
  const navigate = useNavigate();
  
  useEffect(() => {
    if (!getTokens().access || !getTokens().refresh) {
      navigate("/login");
    }
  }, [navigate]);
  
  return (
    <>
      <AccountLink />
      <SearchBar />
      <Posts />
      <p style={{width: "100%", position: "absolute", top: "49%", textAlign: "center"}}>It's Your Feed.</p>
      <p style={{width: "100%", position: "absolute", top: "51%", textAlign: "center"}}>Enjoy this label.</p>
    </>
  )
}

export default Feed;
