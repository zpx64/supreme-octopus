import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getTokens } from '../TokensManagment/TokensManagment';

import Scrollbar from '../Scrollbar/Scrollbar';
import AccountLink from './AccountLink/AccountLink';
import SearchBar from './SearchBar/SearchBar';
import Posts from './Posts/Posts';
import LeftBar from './LeftBar/LeftBar.js'
import CreateEntry from './LeftBar/CreateEntry/CreateEntry.js';


function Feed() {
  const navigate = useNavigate();
  
  useEffect(() => {
    if (!getTokens().access || !getTokens().refresh) {
      navigate("/login");
    }
  }, [navigate]);
  
  return (
    <>
      <LeftBar />
      <AccountLink />
      <SearchBar />
      <Scrollbar />
      <Posts />
      <CreateEntry />
    </>
  )
}

export default Feed;
