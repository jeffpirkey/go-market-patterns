import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import {Route, withRouter} from 'react-router-dom';
import TickerName from './Merk';

class TickerNameSearch extends Component {
    constructor() {
        super();
        this.state = {
            names: [],
        };
    }
}

componentDidMount() {
    let initialNames = [];
    fetch('localhost:7666/api/ticker-names/')
        .then(response => {
            return response.json();
        }).then(data => {
        initialNames = data.results.map((name) => {
            return name
        });
        console.log(initialNames);
        this.setState({
            names: initialNames,
        });
    });
}

render() {
    return (
        <TickerName state={this.state}/>
);
}

// after component is finished

export default PlanetSearch;

ReactDOM.render(<PlanetSearch />, document.getElementById('react-search'));