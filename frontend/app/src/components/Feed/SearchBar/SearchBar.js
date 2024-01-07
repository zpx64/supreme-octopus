import './SearchBar.css';

function Search() {
  return (
    <>
      <div className="feed-search-bar-container">
        <div className="feed-search-bar">
          <input type="text" name="search" id="globalSearch" placeholder='search' />
        </div>
      </div>
    </>
  )
}

export default Search;
