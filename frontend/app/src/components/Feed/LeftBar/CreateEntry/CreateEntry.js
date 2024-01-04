import { useState } from 'react';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import { registerImages } from './sendData';

import './CreateEntry.css';
import Newspaper from './assets/Newspaper.svg';
import Note from './assets/Note.svg';

function NoteMode() {
  return (
    <div className="create-note-window">
      <textarea name="note" placeholder="Enter Note"></textarea>
    </div>
  )
}

function ArticleMode() {
  const [renderedMarkdown, setRenderedMarkdown] = useState("");

  const renderMarkdown = (e) => {
    setRenderedMarkdown(marked.parse(e.target.value));
  }

  return (
    <div className="create-article-window">
      <div className="create-article-header">
        <p>Markdown</p>
        <p>Rendered</p>
      </div>
      <div className="create-article-wrapper">
        <div className="create-article-markdown">
          <textarea name="markdown" onChange={renderMarkdown} placeholder="Enter Markdown"></textarea>
        </div>
        <div className="create-article-separator"></div>
        <div className="create-article-rendered" dangerouslySetInnerHTML={{__html: DOMPurify.sanitize(renderedMarkdown)}}></div>
      </div>
    </div>
  )
}

function CreateEntry() {
  const [isArticleEnabled, setIsArticleEnabled] = useState(true);
  const [getUploadedImagesSRC, setUploadedImagesSRC] = useState([]);
  const [getUploadedImagesFiles, setUploadedImagesFiles] = useState([]);

  const enableArticles  = (e) => {
    setIsArticleEnabled(true);
  }

  const enableNotes = (e) => {
    setIsArticleEnabled(false);
  }

  const updateUploadedImages = (e) => {
    const files = e.target.files;
    let newImages = [];
    let newImagesFiles = [];

    for (const file of files) {
      if (file.type.startsWith('image/')) {
        newImages.push(URL.createObjectURL(file));
      }

      newImagesFiles.push(file);
    }

    setUploadedImagesFiles([...getUploadedImagesFiles, ...newImagesFiles]);
    setUploadedImagesSRC([...getUploadedImagesSRC, ...newImages]);
  }

  const uploadData = () => {
    registerImages(getUploadedImagesFiles);
  }

  return (
    <div className="create-entry-container">
      <div className="window-create-note">
        <div className="window-header">
          <p>{ isArticleEnabled ? "entry / createArticle" : "entry / createNote" }</p>
          <a href="https://google.com/"> </a>
        </div>
        <div className="window-area display-flex">
          <div className="create-entry-vertical-bar">
            <button className={`${isArticleEnabled ? "entry-vertical-bar-item-selected" : "entry-vertical-bar-item"}`} style={{borderRadius: "8px 8px 0 0"}} onClick={enableArticles}>
              <div className="entry-vertical-bar-item">
                <div>
                  <img src={Newspaper} alt="" />
                </div>
                <p>Article</p>
              </div>
            </button>
            <button className={`${isArticleEnabled ? "entry-vertical-bar-item" : "entry-vertical-bar-item-selected"}`} style={{borderRadius: "0 0 8px 8px"}} onClick={enableNotes}>
              <div className="entry-vertical-bar-item">
                <div>
                  <img src={Note} alt="" />
                </div>
                <p>Note</p>
              </div>
            </button>
            { isArticleEnabled ? <></> : <> 
              <label htmlFor="img-upload" className="entry-upload-image-button">Add Images</label>
              <input type="file" id="img-upload" name="img" accept="image/*" multiple="multiple" style={{display: "none"}} onChange={updateUploadedImages} />
              <div className="entry-uploaded-image-wrapper">
              {getUploadedImagesSRC.map((image, index) => (
                <img src={image} className="entry-uploaded-image" key={index} alt="" />
              ))}
              </div>
            </>
            }
            <button className="entry-upload-button" onClick={uploadData}>Upload</button>
          </div>
          { isArticleEnabled ? <ArticleMode /> : <NoteMode />  }
        </div>
      </div>
    </div>
  )
}

export default CreateEntry;
