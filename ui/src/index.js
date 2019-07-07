import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import TickerNameSearch from "./TickerNameSearch";
import TickerName from "./TickerName";

componentDidMount()
{
    let initialTickers = [];
    fetch('localhost:7666/api/ticker-names/')
        .then(response => {
            return response.json();
        }).then(data => {
        initialTickers = data.results.map((ticker) => {
            return ticker
        });
        console.log(initialTickers);
        this.setState({
            names: initialTickers,
        });
    });
}

render()
{
    return (
        <TickerName state={this.state}/>
    );
}

ReactDOM.render(<App />, document.getElementById('root'));
ReactDOM.render(<TickerNameSearch/>, document.getElementById('react-search'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
