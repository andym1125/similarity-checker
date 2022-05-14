import logo from './logo.svg';
import './App.css';
import React, {useState} from 'react';

function App() 
{
  const [taContent, setTaContent] = useState();
  const [fmtContent, setFmtContent] = useState();

  function checkPlagarism()
  {

    alert("test")

  }

  function onTaChange(e)
  {
    setTaContent(e.target.value)
    processInput(e.target.value)
  }

  let processInput = debounce(content =>
  {
    let fetchParams = 
    {
      method: "POST",
      //mode: "cors",
      headers: 
      {
        "Content-type": "application/json; charset=UTF-8",
        //"Access-Control-Allow-Origin": "*"
      },
      body: JSON.stringify({
        text: content
      })
    }

    alert("fetch")

    fetch('http://localhost:8080/process', fetchParams)
      .then(response => {
        console.log(response)
        return response.json()
      })
      .then(data => 
      {
        /* Substrings should be wrapped in HTML */
        let els = []
        let lastSlice = 0
        //PRECONDITION: None of the substrings overlap
        for(let i = 0; i < content.length; i++)
        {
          //Find one with start at current index
          for(let j = 0; j < data.length; j++)
          {
            if(data[j].start === i)
            {
              els.push(content.slice(lastSlice, data[j].start))
              els.push(<span class="highlight">{content.slice(data[j].start, data[j].end)}</span>)
              lastSlice = data[j].end;
              i = data[j].end - 1;
              break;
            }
          }
        }
        els.push(content.slice(lastSlice))

        console.log(els)
        setFmtContent(els)
      })
      .catch(err => alert(err));
  }, 1000);

  return (
    <div className="App">
      <h1>Enter text to check for plagarism:</h1>
      <textarea id="content" name="content" onChange={onTaChange} value={taContent}></textarea>
      <br></br>
      <button onClick={checkPlagarism}>Check for Plagarism</button>
      <p>{fmtContent}</p>
    </div>
  );
}

/* Source: Modified from Web Dev Simplified
 * https://www.youtube.com/watch?v=cjIswDCKgu0
 */
let timeout;
function debounce(callback, delay = 500)
{
  return (...args) =>
  {
    clearTimeout(timeout);
    timeout = setTimeout(() =>
    {
      callback(...args)
    }, delay)
  }
}

export default App;
