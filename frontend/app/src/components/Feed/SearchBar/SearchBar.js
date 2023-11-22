import React from 'react';
import './SearchBar.css';

function Search() {
  return (
    <>
      <div className="feedSearchBarContainer">
        <div className="feedSearchBar">
          <input type="text" name="search" id="globalSearch" placeholder='search' />
        </div>
      </div>
    </>
  )
}

export default Search;
