import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getTokens } from 'utils/TokensManagment/TokensManagment';

import Scrollbar from '../Scrollbar/Scrollbar';
import AccountLink from './AccountLink/AccountLink';
import SearchBar from './SearchBar/SearchBar';
import Posts from './Posts/Posts';
import LeftBar from './LeftBar/LeftBar.js'
import CreateEntry from './LeftBar/CreateEntry/CreateEntry.js';


function Feed() {
  const [PostCreateWindow, setPostCreateWindow] = useState(false);
  const navigate = useNavigate();
  
  useEffect(() => {
    if (!getTokens().access || !getTokens().refresh) {
      navigate("/login");
    }
  }, []);
  
  return (
    <>
      <LeftBar setPostWindowSwitch={setPostCreateWindow} PostWindowSwitch={PostCreateWindow} />
      <AccountLink />
      <SearchBar />
      <Scrollbar />
      <Posts />
      { PostCreateWindow ? <CreateEntry /> : "" }
    </>
  )
}

export default Feed;
