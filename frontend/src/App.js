import logo from './logo.svg';
import './App.css';

function App() {
  return (
    <div className="App">
      <h1>Enter text to check for plagarism:</h1>
      <form method="POST" action="/">
        <input type="text" name="txt"></input>
        <br></br>
        <input type="submit" value="Check for Plagarism"></input>
      </form>
    </div>
  );
}

export default App;
