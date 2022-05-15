import logo from './logo.svg';
import './App.css';
import React, {useState} from 'react';

function App() 
{
  const [taContent, setTaContent] = useState();
  const [fmtContent, setFmtContent] = useState();
  const [percent, setPercent] = useState();

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
      headers: 
      {
        "Content-type": "application/json; charset=UTF-8",
      },
      body: JSON.stringify({
        text: content
      })
    }

    fetch('http://localhost:8080/process', fetchParams)
      .then(response => response.json())
      .then(json => {
        console.log(json.Substrings)
        setPercent(json.Percentage*100)
        return json.Substrings
      })
      .then(data => 
      {
        let splitContent = content.split(/(\w+\s+)/).filter(e => e.length > 1)

        for(let i = 0; i < data.length; i++)
        {
          splitContent.splice(
            data[i].start, 
            data[i].len, 
            (<span class="highlight">{splitContent.slice(data[i].start, data[i].end)}</span>))
        }

        console.log("Format content: " + splitContent)
        setFmtContent(splitContent)
      })
      .catch(err => alert(err));
  }, 1000);

  return (
    <div className="App">
      <h1>Enter text to check for plagarism:</h1>
      <textarea id="content" name="content" onChange={onTaChange} value={taContent}></textarea>
      <br></br>
      <p>We think it may be</p>
      <h1>{percent}% Similar</h1>
      <p>To other documents</p>
      <br></br>
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
